package commands

import (
	"Raven/internals/database"
	"Raven/internals/parser"
	"fmt"
	"strconv"
)

type DeleteValueFromIndex struct{}

func (c *DeleteValueFromIndex) Execute(args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'delindex' command required 2")
	}
	key := args[0]
	valueIndex, err := strconv.Atoi(args[1])
	if err != nil {
		return "Invalid index format", fmt.Errorf("strconv error: %w", err)
	}

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
		dataType.Value = dataType.Value[:valueIndex] + dataType.Value[valueIndex+1:]
		return "OK\n", nil

	case *parser.ListLiteral:
		if valueIndex < 0 || valueIndex >= len(dataType.Values) {
			return "Index out of range", fmt.Errorf("out of bounds")
		}
		dataType.Values = append(dataType.Values[:valueIndex], dataType.Values[valueIndex+1:]...)
		return "OK\n", nil

	default:
		return "Element is not a string or list literal", fmt.Errorf("type mismatch")
	}
}
