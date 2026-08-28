package commands

import (
	"Raven/internals/database"
	"Raven/internals/parser"
	"fmt"
	"strconv"
)

type ValueByIndex struct {
}

func (c *ValueByIndex) Execute(args []string) (string, error) {
	if len(args) != 2 {
		return "ValueByIndex Accept Only two args\n", fmt.Errorf("error in function args")
	}
	value_key := args[0]
	value_index, err := strconv.Atoi(args[1])
	if err != nil {
		return "Invalid index format\n", fmt.Errorf("strconv error: %w", err)
	}

	element, exists := database.Data_store.GetValue(value_key)
	if !exists {
		return "value don't exists\n", fmt.Errorf("not found element")
	}
	element_value, _ := element.Get()
	switch strElem := element_value.(type) {
	case *parser.StringLiteral:
		if value_index < 0 || value_index >= len(strElem.Value) {
			return "Index out of range\n", fmt.Errorf("out of bounds")
		}
		char := string(strElem.Value[value_index])
		return char + "\n", nil
	case *parser.ListLiteral:
		if value_index < 0 || value_index >= len(strElem.Values) {
			return "Index out of range\n", fmt.Errorf("out of bounds")
		}
		value_in_index := strElem.Values[value_index]
		return value_in_index.String() + "\n", nil
	default:
		return "Element is not a string literal\n", fmt.Errorf("type mismatch")
	}
}
