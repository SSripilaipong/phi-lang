package result

import (
	"testing"

	"github.com/stretchr/testify/assert"

	ps "github.com/SSripilaipong/muto/common/parsing"
	psBase "github.com/SSripilaipong/muto/parser/base"
	st "github.com/SSripilaipong/muto/syntaxtree"
	stResult "github.com/SSripilaipong/muto/syntaxtree/result"
)

func TestConstructor(t *testing.T) {
	t.Run("should parse constructor with local class builder", func(t *testing.T) {
		result := constructor()(psBase.StringToCharTokens(`\[$]abc`))
		expectedResult := stResult.NewConstructor(
			stResult.NewObject(st.NewLocalClass("$"), stResult.FixedParamPart{}),
		)
		expectedRemainder := psBase.IgnoreLineAndColumn(psBase.StringToCharTokens("abc"))
		assert.Equal(t, expectedResult, ps.ResultValue(result))
		assert.Equal(t, expectedRemainder, psBase.IgnoreLineAndColumn(result.X2()))
	})

	t.Run("should parse constructor with object builder", func(t *testing.T) {
		result := constructor()(psBase.StringToCharTokens(`\[print! "Hello"]`))
		expectedResult := stResult.NewConstructor(
			stResult.NewObject(st.NewLocalClass("print!"), stResult.ParamsToFixedParamPart([]stResult.Param{st.NewString(`"Hello"`)})),
		)
		assert.Equal(t, expectedResult, ps.ResultValue(result))
	})

	t.Run("should parse constructor with tag builder", func(t *testing.T) {
		result := constructor()(psBase.StringToCharTokens(`\[.ok]`))
		expectedResult := stResult.NewConstructor(
			stResult.NewObject(st.NewTag(".ok"), stResult.FixedParamPart{}),
		)
		assert.Equal(t, expectedResult, ps.ResultValue(result))
	})

	t.Run("should parse constructor with whitespace after backslash", func(t *testing.T) {
		result := constructor()(psBase.StringToCharTokens(`\ [$]`))
		expectedResult := stResult.NewConstructor(
			stResult.NewObject(st.NewLocalClass("$"), stResult.FixedParamPart{}),
		)
		assert.Equal(t, expectedResult, ps.ResultValue(result))
	})
}
