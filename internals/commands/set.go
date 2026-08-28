package commands

import (
	"Raven/internals/database"
	"Raven/internals/parser"
	"fmt"
	"strconv"
	"time"
)

type SetCommand struct{}

func (c *SetCommand) Execute(args []string) (string, error) {
	if len(args) < 2 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'set' command")
	}
	key := args[0]
	val := args[1]
	var ttl uint64 = 0
	if len(args) >= 3 {
		parsedTTL, err := strconv.ParseUint(args[2], 10, 64)
		if err == nil {
			if parsedTTL > 0 {
				ttl = uint64(time.Now().Unix()) + parsedTTL
			}
		}
	}
	parsedVal, err := parser.ParseValue(val)
	if err != nil || parsedVal == nil {
		database.Data_store.SetValue(key, val, ttl)
	} else {
		database.Data_store.SetValue(key, parsedVal, ttl)
	}
	return "OK\n", nil
}
