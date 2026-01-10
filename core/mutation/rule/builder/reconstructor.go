package builder

import (
	"fmt"

	"github.com/SSripilaipong/go-common/optional"

	"github.com/SSripilaipong/muto/core/base"
	ruleExtractor "github.com/SSripilaipong/muto/core/mutation/rule/extractor"
	ruleMutator "github.com/SSripilaipong/muto/core/mutation/rule/mutator"
	"github.com/SSripilaipong/muto/core/pattern/extractor"
	"github.com/SSripilaipong/muto/core/pattern/parameter"
	stPattern "github.com/SSripilaipong/muto/syntaxtree/pattern"
	stResult "github.com/SSripilaipong/muto/syntaxtree/result"
)

type reconstructorBuilderFactory struct {
	node nodeBuilderFactory
}

func newReconstructorBuilderFactory(nodeFactory nodeBuilderFactory) reconstructorBuilderFactory {
	return reconstructorBuilderFactory{node: nodeFactory}
}

func (f reconstructorBuilderFactory) NewBuilder(recon stResult.Reconstructor) optional.Of[ruleMutator.Builder] { // TODO unit test
	// Handle constructor case (no extractor)
	if !recon.HasExtractor() {
		return optional.Value[ruleMutator.Builder](constructorBuilder{
			builder: NewObjectBuilderFactory().NewBuilder(recon.Builder()),
		})
	}

	// Existing reconstructor logic
	ext := recon.Extractor().Value()
	extractorSample, isValidExtractor := newExtractorWithVariableFactory(ext, extractor.NewVariableFactory()).Return()
	if !isValidExtractor {
		return optional.Empty[ruleMutator.Builder]()
	}

	return optional.Value[ruleMutator.Builder](reconstructorBuilder{
		extractor:       ext,
		builder:         NewObjectBuilderFactory().NewBuilder(recon.Builder()),
		extractorSample: extractorSample,
	})
}

type reconstructorBuilder struct {
	extractor       stPattern.ParamPart
	builder         ruleMutator.Builder
	extractorSample extractor.NodeListExtractor
}

func (r reconstructorBuilder) Build(parameter *parameter.Parameter) optional.Of[base.Node] {
	variableFactory := extractor.NewEmbeddedVariableFactory(parameter.VariableMap(), parameter.VariadicVarMap())
	ext, isValidExtractor := newExtractorWithVariableFactory(r.extractor, variableFactory).Return()
	if !isValidExtractor {
		return optional.Empty[base.Node]()
	}

	embeddedBuilder := withVariablesEmbedded(parameter.VariableMappings(), parameter.VariadicVarMappings(), r.builder)
	return optional.Value[base.Node](NewReconstructor(ext, embeddedBuilder))
}

func (r reconstructorBuilder) VisitClass(f func(base.Class)) {
	ruleMutator.VisitClass(f, r.builder)
}

func (r reconstructorBuilder) DisplayString() string {
	return fmt.Sprintf("\\%s [%s]", extractor.DisplayString(r.extractorSample), NakedDisplayString(r.builder))
}

type constructorBuilder struct {
	builder ruleMutator.Builder
}

func (c constructorBuilder) Build(param *parameter.Parameter) optional.Of[base.Node] {
	embeddedBuilder := withVariablesEmbedded(param.VariableMappings(), param.VariadicVarMappings(), c.builder)
	return optional.Value[base.Node](NewReconstructor(extractor.NewImplicitRightVariadic(nil), embeddedBuilder))
}

func (c constructorBuilder) VisitClass(f func(base.Class)) {
	ruleMutator.VisitClass(f, c.builder)
}

func (c constructorBuilder) DisplayString() string {
	return fmt.Sprintf("\\[%s]", NakedDisplayString(c.builder))
}

func newExtractorWithVariableFactory(pattern stPattern.ParamPart, variableFactory ruleExtractor.VariableFactory) optional.Of[extractor.NodeListExtractor] {
	extractorFactory := ruleExtractor.NewTopLevelFactory(variableFactory)
	return extractorFactory.TopLevel(pattern)
}

type Reconstructor struct {
	extractor extractor.NodeListExtractor
	builder   ruleMutator.Builder
}

func NewReconstructor(extractor extractor.NodeListExtractor, builder ruleMutator.Builder) Reconstructor {
	return Reconstructor{extractor: extractor, builder: builder}
}

func (Reconstructor) NodeType() base.NodeType { return base.NodeTypeReconstructor }

func (s Reconstructor) MutateAsHead(params base.ParamChain) optional.Of[base.Node] {
	newChildren := base.MutateParamChain(params)
	if newChildren.IsNotEmpty() {
		return optional.Value[base.Node](base.NewCompoundObject(s, newChildren.Value()))
	}

	build := optional.JoinFmap(s.builder.Build)
	appendRemainingParams := optional.JoinFmap(appendRemainingParamToNode(params.SliceFromOrEmpty(1)))

	return appendRemainingParams(build(s.extractor.Extract(params.DirectParams())))
}

func (s Reconstructor) TopLevelString() string {
	return s.String()
}

func (s Reconstructor) String() string {
	ext := extractor.DisplayString(s.extractor)
	if len(ext) > 0 {
		ext += " "
	}
	return fmt.Sprintf("\\%s[%s]", ext, NakedDisplayString(s.builder))
}

var _ base.Node = Reconstructor{}
