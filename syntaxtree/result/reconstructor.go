package result

import (
	"github.com/SSripilaipong/go-common/optional"

	"github.com/SSripilaipong/muto/syntaxtree/pattern"
)

type Reconstructor struct {
	extractor optional.Of[pattern.ParamPart]
	builder   Object
}

func NewReconstructor(extractor pattern.ParamPart, builder Object) Reconstructor {
	return Reconstructor{
		extractor: optional.Value(extractor),
		builder:   builder,
	}
}

func NewConstructor(builder Object) Reconstructor {
	return Reconstructor{
		extractor: optional.Empty[pattern.ParamPart](),
		builder:   builder,
	}
}

func (Reconstructor) RuleResultNodeType() NodeType { return NodeTypeReconstructor }

func (Reconstructor) ObjectParamType() ParamType { return ParamTypeSingle }

func (r Reconstructor) Extractor() optional.Of[pattern.ParamPart] {
	return r.extractor
}

func (r Reconstructor) HasExtractor() bool {
	return r.extractor.IsNotEmpty()
}

func (r Reconstructor) Builder() Object {
	return r.builder
}

func UnsafeNodeToReconstructor(x Node) Reconstructor { return x.(Reconstructor) }
