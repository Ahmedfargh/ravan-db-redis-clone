package commands

import (
	"strings"
)

type PingCommand struct{}

func (c *PingCommand) Execute(args []string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " ") + "\n", nil
	}
	return "PONG\n", nil
}
