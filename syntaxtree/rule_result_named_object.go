package syntaxtree

type RuleResultNamedObject struct {
	objectName string
	paramPart  ObjectParamPart
}

func NewRuleResultNamedObject(objectName string, paramPart ObjectParamPart) RuleResultNamedObject {
	return RuleResultNamedObject{objectName: objectName, paramPart: paramPart}
}

func (RuleResultNamedObject) RuleResultType() RuleResultType { return RuleResultTypeNamedObject }

func (obj RuleResultNamedObject) ObjectName() string {
	return obj.objectName
}

func (obj RuleResultNamedObject) ParamPart() ObjectParamPart {
	return obj.paramPart
}
