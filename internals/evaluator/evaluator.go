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
	res, err := e.Evaluate(cmdExpr)
	if err != nil {
		if evalErr, ok := err.(*EvalError); ok {
			evalErr.Query = query
		}
		return "", err
	}
	return res, nil
}

func (e *Evaluator) Evaluate(node parser.Node) (string, error) {
	switch n := node.(type) {
	case *parser.StringLiteral:
		return n.Value, nil

	case *parser.NumberLiteral:
		return n.Value, nil

	case *parser.CommandExpr:
		var evaluatedArgs []string
		for i, argNode := range n.Args {
			val, err := e.Evaluate(argNode)
			if err != nil {
				return "", wrapWithFrame(err, fmt.Sprintf("Evaluating argument %d of command '%s'", i+1, n.CommandName), argNode.String())
			}
			// Trim trailing newline when nested result is passed as argument
			val = strings.TrimSuffix(val, "\n")
			evaluatedArgs = append(evaluatedArgs, val)
		}

		res, err := e.registry.Execute(n.CommandName, evaluatedArgs)
		if err != nil {
			return "", wrapWithFrame(err, fmt.Sprintf("Executing command handler '%s'", n.CommandName), n.String())
		}
		return res, nil

	case *parser.ListLiteral:
		var list_items []string
		for i, list_item := range n.Values {
			val, err := e.Evaluate(list_item)
			if err != nil {
				return "", wrapWithFrame(err, fmt.Sprintf("Evaluating element %d in list literal", i+1), list_item.String())
			}
			// Trim trailing newline when nested result is passed as argument
			val = strings.TrimSuffix(val, "\n")
			list_items = append(list_items, val)
		}
		return fmt.Sprintf("[%s]", strings.Join(list_items, " ")), nil

	default:
		return "", wrapWithFrame(fmt.Errorf("unknown AST node type: %T", node), "Evaluating AST node", "")
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

