package commands

import (
	"Raven/internals/database"
	"fmt"
	"strings"
	"time"
)

type MgetCommand struct{}

func (c *MgetCommand) Execute(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'mget' command")
	}
	var results []string
	now := uint64(time.Now().Unix())

	for _, key := range args {
		valObj, exists := database.Data_store.GetValue(key)
		if !exists || (valObj.Ttl > 0 && now > valObj.Ttl) {
			results = append(results, "(nil)")
		} else {
			results = append(results, fmt.Sprintf("%v", valObj.Value))
		}
	}
	return strings.Join(results, "\n") + "\n", nil
}
