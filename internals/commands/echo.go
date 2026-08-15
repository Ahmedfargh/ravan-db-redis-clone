package commands

import (
	"fmt"
	"strings"
)

type EchoCommand struct{}

func (c *EchoCommand) Execute(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("ERR wrong number of arguments for 'echo' command")
	}
	return strings.Join(args, " ") + "\n", nil
}
