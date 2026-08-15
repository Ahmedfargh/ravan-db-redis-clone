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
