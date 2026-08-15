package commands

import (
	"Raven/internals/database"
	"fmt"
	"strconv"
	"time"
)

type ExpireCommand struct{}

func (c *ExpireCommand) Execute(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'expire' command")
	}
	key := args[0]
	seconds, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		return "", fmt.Errorf("ERR value is not an integer or out of range")
	}

	valObj, exists := database.Data_store.GetValue(key)
	if !exists {
		return "(integer) 0\n", nil
	}

	ttl := uint64(time.Now().Unix()) + seconds
	database.Data_store.SetValue(key, valObj.Value, ttl)
	return "(integer) 1\n", nil
}
