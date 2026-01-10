package base

import "github.com/SSripilaipong/go-common/optional"

type Rule interface {
	Active(obj Object) optional.Of[Node]
	Normal(obj Object) optional.Of[Node]
}

type NoActiveRule struct{}

func (NoActiveRule) Active(Object) optional.Of[Node] {
	return optional.Empty[Node]()
}

type RuleBasedClass struct {
	name    string
	mutator Rule
}

func (c *RuleBasedClass) Children() []Node {
	return nil
}

func (c *RuleBasedClass) NodeType() NodeType {
	return NodeTypeClass
}

func (c *RuleBasedClass) ClassType() ClassType { return ClassTypeRuleBased }

func (c *RuleBasedClass) MutateAsHead(params ParamChain) optional.Of[Node] {
	if result, ok := c.ActivelyMutateWithObjMutateFunc(params).Return(); ok {
		return optional.Value(result)
	}
	return c.MutateWithObjMutateFunc(params)
}

func (c *RuleBasedClass) ActivelyMutateWithObjMutateFunc(params ParamChain) optional.Of[Node] {
	if c.mutator == nil {
		return optional.Empty[Node]()
	}
	return c.mutator.Active(NewCompoundObject(c, params))
}

func (c *RuleBasedClass) MutateWithObjMutateFunc(params ParamChain) optional.Of[Node] {
	newChildren := MutateParamChain(params)
	if newChildren.IsNotEmpty() {
		return optional.Value[Node](NewCompoundObject(c, newChildren.Value()))
	}
	if c.mutator == nil {
		return optional.Empty[Node]()
	}
	return c.mutator.Normal(NewCompoundObject(c, params))
}

func (c *RuleBasedClass) Name() string {
	return c.name
}

func (c *RuleBasedClass) TopLevelString() string {
	return c.String()
}

func (c *RuleBasedClass) String() string {
	return c.Name()
}

func (c *RuleBasedClass) MutoString() string {
	return c.String()
}

func (c *RuleBasedClass) Equals(d *RuleBasedClass) bool {
	return c.Name() == d.Name()
}

func (c *RuleBasedClass) LinkRule(mutator Rule) {
	c.mutator = mutator
}

func (c *RuleBasedClass) UnlinkRule() {
	c.mutator = nil
}

func (c *RuleBasedClass) Rule() Rule {
	return c.mutator
}

func NewRuleBasedClass(name string, mutator Rule) *RuleBasedClass {
	return &RuleBasedClass{name: name, mutator: mutator}
}

func NewUnlinkedRuleBasedClass(name string) *RuleBasedClass {
	return NewRuleBasedClass(name, nil)
}

var _ Node = &RuleBasedClass{}

func UnsafeNodeToRuleBasedClass(node Node) *RuleBasedClass {
	return node.(*RuleBasedClass)
}

func UnsafeClassToRuleBasedClass(class Class) *RuleBasedClass {
	return class.(*RuleBasedClass)
}
