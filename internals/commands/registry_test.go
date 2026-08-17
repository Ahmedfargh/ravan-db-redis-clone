package commands

import (
	"Raven/internals/database"
	"strings"
	"testing"
)

func TestBuiltinCommands(t *testing.T) {
	database.InitiatDataStore()

	// PING
	res, err := GlobalRegistry.Execute("PING", nil)
	if err != nil || res != "PONG\n" {
		t.Fatalf("Expected PONG, got %q", res)
	}

	// ECHO
	res, err = GlobalRegistry.Execute("ECHO", []string{"hello", "world"})
	if err != nil || res != "hello world\n" {
		t.Fatalf("Expected 'hello world\\n', got %q", res)
	}

	// SET & GET
	GlobalRegistry.Execute("SET", []string{"count", "10"})
	res, _ = GlobalRegistry.Execute("GET", []string{"count"})
	if res != "10\n" {
		t.Fatalf("Expected 10\\n, got %q", res)
	}

	// INCR & DECR
	res, _ = GlobalRegistry.Execute("INCR", []string{"count"})
	if res != "(integer) 11\n" {
		t.Fatalf("Expected (integer) 11\\n, got %q", res)
	}

	res, _ = GlobalRegistry.Execute("DECR", []string{"count"})
	if res != "(integer) 10\n" {
		t.Fatalf("Expected (integer) 10\\n, got %q", res)
	}

	// MSET & MGET
	GlobalRegistry.Execute("MSET", []string{"k1", "v1", "k2", "v2"})
	res, _ = GlobalRegistry.Execute("MGET", []string{"k1", "k2", "missing"})
	expectedMget := "v1\nv2\n(nil)\n"
	if res != expectedMget {
		t.Fatalf("Expected %q, got %q", expectedMget, res)
	}

	// EXISTS
	res, _ = GlobalRegistry.Execute("EXISTS", []string{"k1", "k2", "nonexistent"})
	if res != "(integer) 2\n" {
		t.Fatalf("Expected (integer) 2\\n, got %q", res)
	}

	// DEL
	res, _ = GlobalRegistry.Execute("DEL", []string{"k1", "k2"})
	if res != "(integer) 2\n" {
		t.Fatalf("Expected (integer) 2\\n, got %q", res)
	}

	res, _ = GlobalRegistry.Execute("GET", []string{"k1"})
	if res != "(nil)\n" {
		t.Fatalf("Expected (nil)\\n, got %q", res)
	}

	// EXPIRE & TTL
	GlobalRegistry.Execute("SET", []string{"temp_key", "val"})
	GlobalRegistry.Execute("EXPIRE", []string{"temp_key", "10"})
	res, _ = GlobalRegistry.Execute("TTL", []string{"temp_key"})
	if !strings.HasPrefix(res, "(integer)") {
		t.Fatalf("Expected integer TTL, got %q", res)
	}

	// KEYS
	res, _ = GlobalRegistry.Execute("KEYS", []string{"*"})
	if !strings.Contains(res, "count") {
		t.Fatalf("Expected keys list to contain 'count', got %q", res)
	}
}

func TestUnknownCommand(t *testing.T) {
	_, err := GlobalRegistry.Execute("INVALID_CMD", nil)
	if err == nil {
		t.Fatalf("Expected error for unknown command")
	}
}
