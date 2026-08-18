package commands

import (
	"Raven/internals/database"
	"Raven/internals/parser"
	"fmt"
	"strconv"
)

type UpdateListIndexValue struct{}

func (c *UpdateListIndexValue) Execute(args []string) (string, error) {
	if len(args) != 3 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'updateindex' command required 3")
	}
	key := args[0]
	valueIndex, err := strconv.Atoi(args[1])
	if err != nil {
		return "Invalid index format", fmt.Errorf("strconv error: %w", err)
	}
	newVal := args[2]

	dataValue, exists := database.Data_store.GetValue(key)
	if !exists {
		return "value don't exists", fmt.Errorf("not found element")
	}

	dataValue.Lock()
	defer dataValue.Unlock()

	switch dataType := dataValue.Value.(type) {
	case *parser.StringLiteral:
		if valueIndex < 0 || valueIndex >= len(dataType.Value) {
			return "Index out of range", fmt.Errorf("out of bounds")
		}
		dataType.Value = dataType.Value[:valueIndex] + newVal + dataType.Value[valueIndex+1:]
		return "OK\n", nil

	case *parser.ListLiteral:
		if valueIndex < 0 || valueIndex >= len(dataType.Values) {
			return "Index out of range", fmt.Errorf("out of bounds")
		}
		var newNode parser.Node
		if _, err := strconv.ParseFloat(newVal, 64); err == nil {
			newNode = &parser.NumberLiteral{Value: newVal}
		} else {
			newNode = &parser.StringLiteral{Value: newVal}
		}
		dataType.Values[valueIndex] = newNode
		return "OK\n", nil

	default:
		return "Element is not a string or list literal", fmt.Errorf("type mismatch")
	}
}
