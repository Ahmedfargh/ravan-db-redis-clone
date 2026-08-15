package commands

import (
	"Raven/internals/database"
	"fmt"
)

type DelCommand struct{}

func (c *DelCommand) Execute(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'del' command")
	}
	deletedCount := 0
	for _, key := range args {
		if database.Data_store.DelValue(key) {
			deletedCount++
		}
	}
	return fmt.Sprintf("(integer) %d\n", deletedCount), nil
}
