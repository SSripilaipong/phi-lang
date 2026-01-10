package global

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/SSripilaipong/muto/builtin/portal"
	"github.com/SSripilaipong/muto/core/base"
)

func createChannel() (sender base.Class, receiver base.Class) {
	port := portal.NewLocalChannel()
	result, ok := port.Call([]base.Node{base.Null()}).Return()
	if !ok {
		panic("failed to create channel")
	}
	obj := base.UnsafeNodeToObject(result)
	params := obj.ParamChain().DirectParams()
	return base.UnsafeNodeToClass(params[0]), base.UnsafeNodeToClass(params[1])
}

func TestSelectMutator(t *testing.T) {
	t.Run("RecvCase", func(t *testing.T) {
		sender, receiver := createChannel()
		callback := base.NewUnlinkedRuleBasedClass("callback")

		go func() {
			time.Sleep(10 * time.Millisecond)
			base.MutateUntilTerminated(base.NewOneLayerObject(sender, base.NewString("hello")))
		}()

		selectObj := base.NewOneLayerObject(
			base.NewUnlinkedRuleBasedClass("select"),
			base.NewOneLayerObject(base.NewTag("recv"), receiver, callback),
		)

		mutator := newSelectMutator()
		result, ok := mutator.Mutate(base.UnsafeNodeToObject(selectObj)).Return()

		expected := base.NewOneLayerObject(callback,
			base.NewOneLayerObject(base.NewTag("value"), base.NewString("hello")),
		)
		assert.True(t, ok)
		assert.True(t, base.NodeEqual(result, expected))
	})

	t.Run("SendCase", func(t *testing.T) {
		sender, receiver := createChannel()
		callback := base.NewUnlinkedRuleBasedClass("callback")

		go func() {
			time.Sleep(10 * time.Millisecond)
			base.MutateUntilTerminated(base.NewOneLayerObject(receiver))
		}()

		selectObj := base.NewOneLayerObject(
			base.NewUnlinkedRuleBasedClass("select"),
			base.NewOneLayerObject(base.NewTag("send"), sender, base.NewString("hello"), callback),
		)

		mutator := newSelectMutator()
		result, ok := mutator.Mutate(base.UnsafeNodeToObject(selectObj)).Return()

		expected := base.NewOneLayerObject(callback, base.NewTag("ok"))
		assert.True(t, ok)
		assert.True(t, base.NodeEqual(result, expected))
	})

	t.Run("DefaultCaseNoChannelReady", func(t *testing.T) {
		_, receiver := createChannel()
		recvCallback := base.NewUnlinkedRuleBasedClass("recvCallback")
		defaultCallback := base.NewUnlinkedRuleBasedClass("defaultCallback")

		selectObj := base.NewOneLayerObject(
			base.NewUnlinkedRuleBasedClass("select"),
			base.NewOneLayerObject(base.NewTag("recv"), receiver, recvCallback),
			base.NewOneLayerObject(base.NewTag("default"), defaultCallback),
		)

		mutator := newSelectMutator()
		result, ok := mutator.Mutate(base.UnsafeNodeToObject(selectObj)).Return()

		expected := base.NewOneLayerObject(defaultCallback)
		assert.True(t, ok)
		assert.True(t, base.NodeEqual(result, expected))
	})

	t.Run("RecvCaseChannelClosed", func(t *testing.T) {
		sender, receiver := createChannel()
		callback := base.NewUnlinkedRuleBasedClass("callback")

		rbc := base.UnsafeClassToRuleBasedClass(sender)
		holder := rbc.Rule().(portal.ChannelHolder)
		close(holder.Channel())

		selectObj := base.NewOneLayerObject(
			base.NewUnlinkedRuleBasedClass("select"),
			base.NewOneLayerObject(base.NewTag("recv"), receiver, callback),
		)

		mutator := newSelectMutator()
		result, ok := mutator.Mutate(base.UnsafeNodeToObject(selectObj)).Return()

		expected := base.NewOneLayerObject(callback, base.NewTag("closed"))
		assert.True(t, ok)
		assert.True(t, base.NodeEqual(result, expected))
	})

	t.Run("MultipleCasesFirstReady", func(t *testing.T) {
		sender1, receiver1 := createChannel()
		_, receiver2 := createChannel()
		callback1 := base.NewUnlinkedRuleBasedClass("callback1")
		callback2 := base.NewUnlinkedRuleBasedClass("callback2")

		go func() {
			time.Sleep(10 * time.Millisecond)
			base.MutateUntilTerminated(base.NewOneLayerObject(sender1, base.NewString("from-ch1")))
		}()

		selectObj := base.NewOneLayerObject(
			base.NewUnlinkedRuleBasedClass("select"),
			base.NewOneLayerObject(base.NewTag("recv"), receiver1, callback1),
			base.NewOneLayerObject(base.NewTag("recv"), receiver2, callback2),
		)

		mutator := newSelectMutator()
		result, ok := mutator.Mutate(base.UnsafeNodeToObject(selectObj)).Return()

		expected := base.NewOneLayerObject(callback1,
			base.NewOneLayerObject(base.NewTag("value"), base.NewString("from-ch1")),
		)
		assert.True(t, ok)
		assert.True(t, base.NodeEqual(result, expected))
	})

	t.Run("InvalidCaseWrongDirection", func(t *testing.T) {
		sender, _ := createChannel()
		callback := base.NewUnlinkedRuleBasedClass("callback")

		selectObj := base.NewOneLayerObject(
			base.NewUnlinkedRuleBasedClass("select"),
			base.NewOneLayerObject(base.NewTag("recv"), sender, callback),
		)

		mutator := newSelectMutator()
		_, ok := mutator.Mutate(base.UnsafeNodeToObject(selectObj)).Return()

		assert.False(t, ok)
	})

	t.Run("EmptyParams", func(t *testing.T) {
		selectObj := base.NewOneLayerObject(
			base.NewUnlinkedRuleBasedClass("select"),
		)

		mutator := newSelectMutator()
		_, ok := mutator.Mutate(base.UnsafeNodeToObject(selectObj)).Return()

		assert.False(t, ok)
	})

	t.Run("SendToClosedChannelReturnsFailure", func(t *testing.T) {
		sender, _ := createChannel()
		callback := base.NewUnlinkedRuleBasedClass("callback")

		rbc := base.UnsafeClassToRuleBasedClass(sender)
		holder := rbc.Rule().(portal.ChannelHolder)
		close(holder.Channel())

		selectObj := base.NewOneLayerObject(
			base.NewUnlinkedRuleBasedClass("select"),
			base.NewOneLayerObject(base.NewTag("send"), sender, base.NewString("hello"), callback),
		)

		mutator := newSelectMutator()
		result, ok := mutator.Mutate(base.UnsafeNodeToObject(selectObj)).Return()

		assert.True(t, ok)
		assert.True(t, base.NodeEqual(result, base.NewTag("failure")))
	})
}
