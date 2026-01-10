package global

import (
	"reflect"

	"github.com/SSripilaipong/go-common/optional"

	"github.com/SSripilaipong/muto/builtin/portal"
	"github.com/SSripilaipong/muto/core/base"
	"github.com/SSripilaipong/muto/core/mutation/rule/mutator"
)

const selectMutatorName = "select"

type selectMutator struct{}

func newSelectMutator() *selectMutator {
	return &selectMutator{}
}

func (s *selectMutator) Name() string { return selectMutatorName }

func (s *selectMutator) Mutate(obj base.Object) optional.Of[base.Node] {
	params := obj.ParamChain().DirectParams()
	if len(params) == 0 {
		return optional.Empty[base.Node]()
	}

	selectCases, callbacks, err := s.parseCases(params)
	if err != nil {
		return optional.Empty[base.Node]()
	}

	chosenIdx, recvValue, recvOK, panicked := s.doSelect(selectCases)

	if panicked {
		return optional.Value[base.Node](base.NewTag("failure"))
	}

	callback := callbacks[chosenIdx]
	caseType := selectCases[chosenIdx].Dir

	switch caseType {
	case reflect.SelectRecv:
		if !recvOK {
			return optional.Value[base.Node](
				base.NewOneLayerObject(callback, base.NewTag("closed")),
			)
		}
		receivedNode := recvValue.Interface().(base.Node)
		return optional.Value[base.Node](
			base.NewOneLayerObject(callback, base.NewOneLayerObject(base.NewTag("value"), receivedNode)),
		)
	case reflect.SelectSend:
		return optional.Value[base.Node](
			base.NewOneLayerObject(callback, base.NewTag("ok")),
		)
	case reflect.SelectDefault:
		return optional.Value[base.Node](
			base.NewOneLayerObject(callback),
		)
	}

	return optional.Empty[base.Node]()
}

func (s *selectMutator) doSelect(selectCases []reflect.SelectCase) (chosenIdx int, recvValue reflect.Value, recvOK bool, panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()
	chosenIdx, recvValue, recvOK = reflect.Select(selectCases)
	return
}

func (s *selectMutator) VisitClass(mutator.ClassVisitor) {}

func (s *selectMutator) parseCases(params []base.Node) ([]reflect.SelectCase, []base.Node, error) {
	var selectCases []reflect.SelectCase
	var callbacks []base.Node

	for _, param := range params {
		if !base.IsObjectNode(param) {
			return nil, nil, errInvalidCase
		}

		caseObj := base.UnsafeNodeToObject(param)
		head := caseObj.Head()

		if !base.IsTagNode(head) {
			return nil, nil, errInvalidCase
		}

		tag := base.UnsafeNodeToTag(head)
		caseParams := caseObj.ParamChain().DirectParams()

		switch tag.Name() {
		case "recv":
			if len(caseParams) != 2 {
				return nil, nil, errInvalidCase
			}
			holder, err := extractChannelHolder(caseParams[0])
			if err != nil {
				return nil, nil, err
			}
			if holder.Direction() != portal.ChannelDirectionRecv {
				return nil, nil, errInvalidCase
			}
			selectCases = append(selectCases, reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(holder.Channel()),
			})
			callbacks = append(callbacks, caseParams[1])

		case "send":
			if len(caseParams) != 3 {
				return nil, nil, errInvalidCase
			}
			holder, err := extractChannelHolder(caseParams[0])
			if err != nil {
				return nil, nil, err
			}
			if holder.Direction() != portal.ChannelDirectionSend {
				return nil, nil, errInvalidCase
			}
			selectCases = append(selectCases, reflect.SelectCase{
				Dir:  reflect.SelectSend,
				Chan: reflect.ValueOf(holder.Channel()),
				Send: reflect.ValueOf(caseParams[1]),
			})
			callbacks = append(callbacks, caseParams[2])

		case "default":
			if len(caseParams) != 1 {
				return nil, nil, errInvalidCase
			}
			selectCases = append(selectCases, reflect.SelectCase{
				Dir: reflect.SelectDefault,
			})
			callbacks = append(callbacks, caseParams[0])

		default:
			return nil, nil, errInvalidCase
		}
	}

	return selectCases, callbacks, nil
}

func extractChannelHolder(node base.Node) (portal.ChannelHolder, error) {
	if !base.IsClassNode(node) {
		return nil, errInvalidCase
	}
	class := base.UnsafeNodeToClass(node)
	if class.ClassType() != base.ClassTypeRuleBased {
		return nil, errInvalidCase
	}
	rbc := base.UnsafeClassToRuleBasedClass(class)
	rule := rbc.Rule()
	if holder, ok := rule.(portal.ChannelHolder); ok {
		return holder, nil
	}
	return nil, errInvalidCase
}

type selectError struct {
	msg string
}

func (e selectError) Error() string {
	return e.msg
}

var errInvalidCase = selectError{msg: "invalid select case"}
