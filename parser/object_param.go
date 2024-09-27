package parser

import (
	ps "github.com/SSripilaipong/muto/common/parsing"
	"github.com/SSripilaipong/muto/common/tuple"
	"github.com/SSripilaipong/muto/parser/tokenizer"
	st "github.com/SSripilaipong/muto/syntaxtree"
	stResult "github.com/SSripilaipong/muto/syntaxtree/result"
)

func objectParamPart(xs []tokenizer.Token) []tuple.Of2[stResult.ParamPart, []tokenizer.Token] {
	return ps.Map(stResult.ParamsToParamPart, ps.OptionalGreedyRepeat(objectParam))(xs)
}

func objectParam(xs []tokenizer.Token) []tuple.Of2[stResult.Param, []tokenizer.Token] {
	return ps.Or(
		objectParamValue,
		ps.Map(variadicVariableToObjectParam, variadicVar),
		ps.Map(variableToObjectParam, variable),
		ps.Map(classToObjectParam, classIncludingSymbols),
		ps.Map(objectToObjectParam, inParentheses(object)),
	)(xs)
}

var objectParamValue = ps.Or(
	ps.Map(booleanToObjectParam, boolean),
	ps.Map(stringToObjectParam, string_),
	ps.Map(numberToObjectParam, number),
)

func objectToObjectParam(obj objectNode) stResult.Param {
	return stResult.NewObject(obj.Head(), obj.ParamPart())
}

func classToObjectParam(objName tokenizer.Token) stResult.Param {
	return st.NewClass(objName.Value())
}

func booleanToObjectParam(x tokenizer.Token) stResult.Param {
	return st.NewBoolean(x.Value())
}

func numberToObjectParam(x tokenizer.Token) stResult.Param {
	return st.NewNumber(x.Value())
}

func stringToObjectParam(x tokenizer.Token) stResult.Param {
	s := x.Value()
	return st.NewString(s[1 : len(s)-1])
}

func variableToObjectParam(x tokenizer.Token) stResult.Param {
	return st.NewVariable(x.Value())
}

func variadicVariableToObjectParam(x variadicVarNode) stResult.Param {
	return stResult.NewVariadicVariable(x.Name())
}
