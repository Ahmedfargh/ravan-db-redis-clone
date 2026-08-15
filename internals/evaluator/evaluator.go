package evaluator

import (
	"Raven/internals/commands"
	"Raven/internals/parser"
	"fmt"
	"strings"
)

func EvaluateQuery(query string) (string, error) {
	cmdExpr, err := parser.ParseQuery(query)
	if err != nil {
		return "", err
	}
	return Evaluate(cmdExpr)
}

func Evaluate(node parser.Node) (string, error) {
	switch n := node.(type) {
	case *parser.StringLiteral:
		return n.Value, nil

	case *parser.NumberLiteral:
		return n.Value, nil

	case *parser.CommandExpr:
		var evaluatedArgs []string
		for _, argNode := range n.Args {
			val, err := Evaluate(argNode)
			if err != nil {
				return "", err
			}
			// Trim trailing newline when nested result is passed as argument
			val = strings.TrimSuffix(val, "\n")
			evaluatedArgs = append(evaluatedArgs, val)
		}

		res, err := commands.GlobalRegistry.Execute(n.CommandName, evaluatedArgs)
		if err != nil {
			return "", err
		}
		return res, nil

	default:
		return "", fmt.Errorf("unknown AST node type: %T", node)
	}
}
