package commands

import (
	"Raven/internals/database"
	"fmt"
	"time"
)

type GetCommand struct{}

func (c *GetCommand) Execute(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'get' command")
	}
	key := args[0]
	valObj, exists := database.Data_store.GetValue(key)
	if !exists {
		return "(nil)\n", nil
	}

	// Check TTL expiration
	if valObj.Ttl > 0 && uint64(time.Now().Unix()) > valObj.Ttl {
		database.Data_store.DelValue(key)
		return "(nil)\n", nil
	}

	return fmt.Sprintf("%v\n", valObj.Value), nil
}
