package commands

import (
	"Raven/internals/database"
	"Raven/internals/parser"
	"fmt"
	"strings"
)

type DeleteFromList struct{}

func (c *DeleteFromList) Execute(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'delfromlist' command")
	}
	key := args[0]
	itemsToDelete := args[1:]

	dataValue, exists := database.Data_store.GetValue(key)
	if !exists {
		return "value don't exists", fmt.Errorf("not found element")
	}

	dataValue.Lock()
	defer dataValue.Unlock()

	toDeleteMap := make(map[string]bool)
	for _, item := range itemsToDelete {
		toDeleteMap[item] = true
	}

	switch dataType := dataValue.Value.(type) {
	case *parser.ListLiteral:
		var newValues []parser.Node
		for _, node := range dataType.Values {
			if !toDeleteMap[node.String()] {
				newValues = append(newValues, node)
			}
		}
		dataType.Values = newValues
		return "OK\n", nil

	case *parser.StringLiteral:
		for _, item := range itemsToDelete {
			if item != "" {
				dataType.Value = strings.ReplaceAll(dataType.Value, item, "")
			}
		}
		return "OK\n", nil

	default:
		return "Element is not a string or list literal", fmt.Errorf("type mismatch")
	}
}
