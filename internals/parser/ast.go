package parser

import (
	"fmt"
	"strings"
)

type Node interface {
	String() string
}

type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) String() string {
	return sl.Value
}

type NumberLiteral struct {
	Value string
}

func (nl *NumberLiteral) String() string {
	return nl.Value
}

type ListLiteral struct {
	Values []Node
}

func (ll *ListLiteral) string() string {
	var argsStr []string
	for _, arg := range ll.Values {
		argsStr = append(argsStr, arg.String())
	}
	return "[" + strings.Join(argsStr, ",") + "]"
}

type CommandExpr struct {
	CommandName string
	Args        []Node
}

func (ce *CommandExpr) String() string {
	var argsStr []string
	for _, arg := range ce.Args {
		argsStr = append(argsStr, arg.String())
	}
	return fmt.Sprintf("(%s %s)", ce.CommandName, strings.Join(argsStr, " "))
}
