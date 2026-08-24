package commands

import (
	"Raven/internals/database"
	"strings"
)

type KeysCommand struct{}

func (c *KeysCommand) Execute(args []string) (string, error) {
	keysList := database.Data_store.GetKeys()
	if len(keysList) == 0 {
		return "(empty list or set)\n", nil
	}
	return strings.Join(keysList, "\n") + "\n", nil
}
