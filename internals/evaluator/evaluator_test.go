package evaluator

import (
	"Raven/internals/commands"
	"Raven/internals/database"
	"strings"
	"testing"
)

func TestEvaluatorBasicAndQuotedStrings(t *testing.T) {
	database.InitiatDataStore()
	eval := NewEvaluator(nil)

	// Test SET with spaces in quotes
	res, err := eval.EvaluateQuery(`SET greeting "Hello World from Raven"`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if res != "OK\n" {
		t.Fatalf("Expected OK\\n, got %q", res)
	}

	// Test GET
	res, err = eval.EvaluateQuery(`GET greeting`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if res != "Hello World from Raven\n" {
		t.Fatalf("Expected 'Hello World from Raven\\n', got %q", res)
	}
}

func TestEvaluatorNestedCommands(t *testing.T) {
	database.InitiatDataStore()
	eval := NewEvaluator(nil)

	// Set original key
	eval.EvaluateQuery(`SET default_role "Admin"`)

	// Execute nested command: SET user_role (GET default_role)
	res, err := eval.EvaluateQuery(`SET user_role (GET default_role)`)
	if err != nil {
		t.Fatalf("Unexpected error in nested query: %v", err)
	}
	if res != "OK\n" {
		t.Fatalf("Expected OK\\n, got %q", res)
	}

	// GET user_role should now be "Admin"
	res, err = eval.EvaluateQuery(`GET user_role`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if res != "Admin\n" {
		t.Fatalf("Expected 'Admin\\n', got %q", res)
	}
}

func TestEvaluatorListLiteralWithNestedCommand(t *testing.T) {
	database.InitiatDataStore()
	eval := NewEvaluator(nil)

	// Pre-populate data
	eval.EvaluateQuery(`SET env "production"`)

	// Evaluate command with list literal containing a nested command
	// SET tags [ "server" (GET env) 8080 ]
	res, err := eval.EvaluateQuery(`SET tags [ "server" (GET env) 8080 ]`)
	if err != nil {
		t.Fatalf("Unexpected error in list query: %v", err)
	}
	if res != "OK\n" {
		t.Fatalf("Expected OK\\n, got %q", res)
	}

	// Check saved value
	res, err = eval.EvaluateQuery(`GET tags`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := "[server production 8080]\n"
	if res != expected {
		t.Fatalf("Expected %q, got %q", expected, res)
	}
}

func TestDynamicCommandRegistration(t *testing.T) {
	eval := NewEvaluator(nil)

	// Register a dynamic custom command "ECHO_UPPER" at runtime
	commands.GlobalRegistry.Register("ECHO_UPPER", func(args []string) (string, error) {
		return strings.ToUpper(strings.Join(args, " ")) + "\n", nil
	})

	res, err := eval.EvaluateQuery(`ECHO_UPPER hello raven engine`)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if res != "HELLO RAVEN ENGINE\n" {
		t.Fatalf("Expected 'HELLO RAVEN ENGINE\\n', got %q", res)
	}
}
