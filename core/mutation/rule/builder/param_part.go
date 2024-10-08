package builder

import (
	"github.com/SSripilaipong/muto/common/fn"
	"github.com/SSripilaipong/muto/common/optional"
	"github.com/SSripilaipong/muto/common/slc"
	"github.com/SSripilaipong/muto/core/base"
	"github.com/SSripilaipong/muto/core/mutation/rule/data"
	stResult "github.com/SSripilaipong/muto/syntaxtree/result"
)

func buildChildren(paramPart stResult.ParamPart) func(mapping *data.Mutation) optional.Of[[]base.Node] {
	switch {
	case paramPart == nil:
		return func(*data.Mutation) optional.Of[[]base.Node] { return optional.Value[[]base.Node](nil) }
	case stResult.IsParamPartTypeFixed(paramPart):
		return buildFixedParamPart(stResult.UnsafeParamPartToFixedParamPart(paramPart))
	}
	panic("not implemented")
}

func buildFixedParamPart(part stResult.FixedParamPart) func(mapping *data.Mutation) optional.Of[[]base.Node] {
	buildParams := slc.Map(buildObjectParam)(part)

	return func(mapping *data.Mutation) optional.Of[[]base.Node] {
		return fn.Compose(
			slc.Fold(aggregateNodes)(optional.Value[[]base.Node](nil)),
			slc.Map(fn.CallWith[optional.Of[[]base.Node]](mapping)),
		)(buildParams)
	}
}

var aggregateNodes = optional.MergeFn(func(nodes []base.Node, xs []base.Node) optional.Of[[]base.Node] {
	return optional.Value(append(nodes, xs...))
})
