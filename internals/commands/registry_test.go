package commands

import (
	"Raven/internals/database"
	"testing"
)

func TestBuiltinCommands(t *testing.T) {
	database.InitiatDataStore()

	// PING
	res, err := GlobalRegistry.Execute("PING", nil)
	if err != nil || res != "PONG\n" {
		t.Fatalf("Expected PONG, got %q", res)
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
}
