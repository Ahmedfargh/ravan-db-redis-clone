package commands

import (
	"fmt"
	"strings"
	"sync"
)

type CommandHandler func(args []string) (string, error)

type CommandRegistry struct {
	mu       sync.RWMutex
	handlers map[string]CommandHandler
}

var GlobalRegistry = NewCommandRegistry()

func NewCommandRegistry() *CommandRegistry {
	r := &CommandRegistry{
		handlers: make(map[string]CommandHandler),
	}
	r.registerDefaults()
	return r
}

func (r *CommandRegistry) Register(name string, handler CommandHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[strings.ToUpper(name)] = handler
}

func (r *CommandRegistry) Get(name string) (CommandHandler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	handler, exists := r.handlers[strings.ToUpper(name)]
	return handler, exists
}

func (r *CommandRegistry) Execute(name string, args []string) (string, error) {
	handler, exists := r.Get(name)
	if !exists {
		return "", fmt.Errorf("ERR unknown command '%s'", name)
	}
	return handler(args)
}
