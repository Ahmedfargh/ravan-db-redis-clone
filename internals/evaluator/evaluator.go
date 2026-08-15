package evaluator

import (
	"Raven/internals/commands"
	"Raven/internals/parser"
	"fmt"
	"strings"
)

type Evaluator struct {
	registry *commands.CommandRegistry
}

func NewEvaluator(registry *commands.CommandRegistry) *Evaluator {
	if registry == nil {
		registry = commands.GlobalRegistry
	}
	return &Evaluator{registry: registry}
}

var DefaultEvaluator = NewEvaluator(commands.GlobalRegistry)

func (e *Evaluator) EvaluateQuery(query string) (string, error) {
	p, err := parser.NewParserFromInput(query)
	if err != nil {
		return "", err
	}
	cmdExpr, err := p.Parse()
	if err != nil {
		return "", err
	}
	return e.Evaluate(cmdExpr)
}

func (e *Evaluator) Evaluate(node parser.Node) (string, error) {
	switch n := node.(type) {
	case *parser.StringLiteral:
		return n.Value, nil

	case *parser.NumberLiteral:
		return n.Value, nil

	case *parser.CommandExpr:
		var evaluatedArgs []string
		for _, argNode := range n.Args {
			val, err := e.Evaluate(argNode)
			if err != nil {
				return "", err
			}
			// Trim trailing newline when nested result is passed as argument
			val = strings.TrimSuffix(val, "\n")
			evaluatedArgs = append(evaluatedArgs, val)
		}

		res, err := e.registry.Execute(n.CommandName, evaluatedArgs)
		if err != nil {
			return "", err
		}
		return res, nil

	default:
		return "", fmt.Errorf("unknown AST node type: %T", node)
	}
}

// EvaluateQuery is a package-level helper that delegates to DefaultEvaluator.
func EvaluateQuery(query string) (string, error) {
	return DefaultEvaluator.EvaluateQuery(query)
}

// Evaluate is a package-level helper that delegates to DefaultEvaluator.
func Evaluate(node parser.Node) (string, error) {
	return DefaultEvaluator.Evaluate(node)
}
