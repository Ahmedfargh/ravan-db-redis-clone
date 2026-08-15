package commands

import (
	"Raven/internals/database"
	"strings"
	"time"
)

type KeysCommand struct{}

func (c *KeysCommand) Execute(args []string) (string, error) {
	database.Data_store.GetValue("") // Trigger RLock
	var keysList []string
	now := uint64(time.Now().Unix())

	for k, v := range database.Data_store.Data {
		if v.Ttl > 0 && now > v.Ttl {
			continue
		}
		keysList = append(keysList, k)
	}

	if len(keysList) == 0 {
		return "(empty list or set)\n", nil
	}
	return strings.Join(keysList, "\n") + "\n", nil
}
