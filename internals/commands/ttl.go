package commands

import (
	"Raven/internals/database"
	"fmt"
	"time"
)

type TtlCommand struct{}

func (c *TtlCommand) Execute(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'ttl' command")
	}
	key := args[0]
	valObj, exists := database.Data_store.GetValue(key)
	if !exists {
		return "(integer) -2\n", nil // Key does not exist
	}

	if valObj.Ttl == 0 {
		return "(integer) -1\n", nil // Key exists but has no associated expire
	}

	now := uint64(time.Now().Unix())
	if now > valObj.Ttl {
		database.Data_store.DelValue(key)
		return "(integer) -2\n", nil
	}

	remaining := int64(valObj.Ttl - now)
	return fmt.Sprintf("(integer) %d\n", remaining), nil
}
