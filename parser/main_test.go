package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"phi-lang/common/tuple"
	st "phi-lang/parser/syntaxtree"
	tk "phi-lang/tokenizer"
)

func TestParser(t *testing.T) {
	parser := NewParser()
	tokens := []tk.Token{
		tk.NewToken("hello", tk.Identifier),
		tk.NewToken("X", tk.Identifier),
		tk.NewToken("Y", tk.Identifier),
		tk.NewToken("Z", tk.Identifier),
		tk.NewToken("=", tk.Symbol),
		tk.NewToken("Y", tk.Identifier),
		tk.NewToken("\\n", tk.LineBreak),
		tk.NewToken("main", tk.Identifier),
		tk.NewToken("X", tk.Identifier),
		tk.NewToken("=", tk.Symbol),
		tk.NewToken("hello", tk.Identifier),
		tk.NewToken(`"world"`, tk.String),
		tk.NewToken(`123`, tk.Number),
		tk.NewToken("X", tk.Identifier),
	}
	expectedParsedTree := st.NewFile([]st.Statement{
		st.NewRule(st.NewRulePattern("hello", []st.RuleParamPattern{st.NewVariable("X"), st.NewVariable("Y"), st.NewVariable("Z")}), st.NewVariable("Y")),
		st.NewRule(st.NewRulePattern("main", []st.RuleParamPattern{st.NewVariable("X")}), st.NewRuleResultObject("hello", []st.ObjectParam{st.NewString("world"), st.NewNumber("123"), st.NewVariable("X")})),
	})
	assert.Equal(t,
		[]tuple.Of2[st.File, []tk.Token]{tuple.New2(expectedParsedTree, []tk.Token{})},
		parser(tokens),
	)
}
