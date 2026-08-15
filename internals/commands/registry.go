package commands

import (
	"fmt"
	"strings"
	"sync"
)

// Command is the interface that all database commands must implement.
type Command interface {
	Execute(args []string) (string, error)
}

// CommandHandlerFunc adapter allows plain functions to satisfy the Command interface.
type CommandHandlerFunc func(args []string) (string, error)

func (f CommandHandlerFunc) Execute(args []string) (string, error) {
	return f(args)
}

type CommandRegistry struct {
	mu       sync.RWMutex
	commands map[string]Command
}

var GlobalRegistry = NewCommandRegistry()

func NewCommandRegistry() *CommandRegistry {
	r := &CommandRegistry{
		commands: make(map[string]Command),
	}
	r.registerDefaults()
	return r
}

func (r *CommandRegistry) Register(name string, handler any) {
	r.mu.Lock()
	defer r.mu.Unlock()

	switch h := handler.(type) {
	case Command:
		r.commands[strings.ToUpper(name)] = h
	case func([]string) (string, error):
		r.commands[strings.ToUpper(name)] = CommandHandlerFunc(h)
	default:
		panic(fmt.Sprintf("invalid command handler type: %T", handler))
	}
}

func (r *CommandRegistry) Get(name string) (Command, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cmd, exists := r.commands[strings.ToUpper(name)]
	return cmd, exists
}

func (r *CommandRegistry) Execute(name string, args []string) (string, error) {
	cmd, exists := r.Get(name)
	if !exists {
		return "", fmt.Errorf("ERR unknown command '%s'", name)
	}
	return cmd.Execute(args)
}
