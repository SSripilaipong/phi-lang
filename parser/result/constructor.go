package result

import (
	"github.com/SSripilaipong/go-common/rslt"
	"github.com/SSripilaipong/go-common/tuple"

	ps "github.com/SSripilaipong/muto/common/parsing"
	psBase "github.com/SSripilaipong/muto/parser/base"
	stResult "github.com/SSripilaipong/muto/syntaxtree/result"
)

func constructor() func(xs []psBase.Character) tuple.Of2[rslt.Of[stResult.Reconstructor], []psBase.Character] {
	backSlash := ps.Sequence2(
		ps.ToParser(psBase.BackSlash),
		ps.OptionalGreedyRepeat(ps.ToParser(psBase.WhiteSpace)),
	)
	builderPart := ps.Map(
		castObjectResult,
		ps.ToParser(psBase.InSquareBracketsWhiteSpacesAllowed(nakedObjectMultilines)),
	)
	return ps.Map(
		stResult.NewConstructor,
		ps.Prefix(backSlash, builderPart),
	).Legacy
}
