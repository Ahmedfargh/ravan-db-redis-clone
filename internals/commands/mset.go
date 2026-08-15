package commands

import (
	"Raven/internals/database"
	"fmt"
)

type MsetCommand struct{}

func (c *MsetCommand) Execute(args []string) (string, error) {
	if len(args) < 2 || len(args)%2 != 0 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'mset' command")
	}
	for i := 0; i < len(args); i += 2 {
		key := args[i]
		val := args[i+1]
		database.Data_store.SetValue(key, val, 0)
	}
	return "OK\n", nil
}
