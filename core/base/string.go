package base

import "fmt"

type String struct {
	value string
}

func NewString(value string) String {
	return String{value: value}
}

func (s String) NodeType() NodeType {
	return NodeTypeString
}

func (s String) Value() string {
	return s.value
}

func (s String) String() string {
	return fmt.Sprintf("%#v", s.value)
}

func UnsafeNodeToString(n Node) String {
	return n.(String)
}
