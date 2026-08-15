package commands

import (
	"Raven/internals/database"
	"fmt"
	"time"
)

type ExistsCommand struct{}

func (c *ExistsCommand) Execute(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'exists' command")
	}
	existsCount := 0
	for _, key := range args {
		valObj, exists := database.Data_store.GetValue(key)
		if exists {
			if valObj.Ttl > 0 && uint64(time.Now().Unix()) > valObj.Ttl {
				database.Data_store.DelValue(key)
			} else {
				existsCount++
			}
		}
	}
	return fmt.Sprintf("(integer) %d\n", existsCount), nil
}
