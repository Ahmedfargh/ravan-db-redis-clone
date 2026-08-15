package commands

import (
	"Raven/internals/database"
	"fmt"
	"strconv"
)

type IncrCommand struct{}

func (c *IncrCommand) Execute(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'incr' command")
	}
	key := args[0]
	valObj, exists := database.Data_store.GetValue(key)
	var currentVal int64 = 0

	if exists {
		strVal := fmt.Sprintf("%v", valObj.Value)
		parsed, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			return "", fmt.Errorf("ERR value is not an integer or out of range")
		}
		currentVal = parsed
	}

	newVal := currentVal + 1
	database.Data_store.SetValue(key, strconv.FormatInt(newVal, 10), 0)
	return fmt.Sprintf("(integer) %d\n", newVal), nil
}
