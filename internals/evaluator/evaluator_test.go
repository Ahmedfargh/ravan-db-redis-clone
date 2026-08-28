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
	expected := "[server,production,8080]\n"
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

func TestEvaluatorUnknownCommandTrace(t *testing.T) {
	eval := NewEvaluator(nil)

	_, err := eval.EvaluateQuery(`SET role (UNKNOWN_CMD key)`)
	if err == nil {
		t.Fatalf("Expected error for unknown command in subquery, got nil")
	}

	evalErr, ok := err.(*EvalError)
	if !ok {
		t.Fatalf("Expected *EvalError, got %T: %v", err, err)
	}

	errStr := evalErr.Error()
	if !strings.Contains(errStr, "[RQL Runtime Error]") {
		t.Errorf("Expected [RQL Runtime Error] banner, got:\n%s", errStr)
	}
	if !strings.Contains(errStr, "UNKNOWN_CMD") {
		t.Errorf("Expected UNKNOWN_CMD in trace, got:\n%s", errStr)
	}
	if !strings.Contains(errStr, "Step 1") || !strings.Contains(errStr, "Step 2") {
		t.Errorf("Expected multi-step evaluation trace, got:\n%s", errStr)
	}
}

func TestEvaluatorListNestedCommandTrace(t *testing.T) {
	eval := NewEvaluator(nil)

	// Sub-command inside list literal missing required arguments
	_, err := eval.EvaluateQuery(`SET tags [ "server" (GET) 8080 ]`)
	if err == nil {
		t.Fatalf("Expected error for missing arguments in list subquery, got nil")
	}

	evalErr, ok := err.(*EvalError)
	if !ok {
		t.Fatalf("Expected *EvalError, got %T: %v", err, err)
	}

	errStr := evalErr.Error()
	if !strings.Contains(errStr, "list literal") {
		t.Errorf("Expected list literal step in trace, got:\n%s", errStr)
	}
	if !strings.Contains(errStr, "wrong number of arguments") {
		t.Errorf("Expected wrong number of arguments cause, got:\n%s", errStr)
	}
}

func TestEvaluatorListAndIndexOperationsEndToEnd(t *testing.T) {
	database.InitiatDataStore()
	eval := NewEvaluator(nil)

	// 1. SET list
	res, err := eval.EvaluateQuery(`SET items ["item1" "item2" "item3"]`)
	if err != nil || res != "OK\n" {
		t.Fatalf("SET items failed: %v, res: %q", err, res)
	}

	// 2. VALUEBYINDEX on list
	res, err = eval.EvaluateQuery(`VALUEBYINDEX items 0`)
	if err != nil || res != "item1\n" {
		t.Fatalf("VALUEBYINDEX items 0 failed: %v, res: %q", err, res)
	}

	res, err = eval.EvaluateQuery(`VALUEBYINDEX items 2`)
	if err != nil || res != "item3\n" {
		t.Fatalf("VALUEBYINDEX items 2 failed: %v, res: %q", err, res)
	}

	// 3. UPDATEINDEX on list
	res, err = eval.EvaluateQuery(`UPDATEINDEX items 1 "updated_item"`)
	if err != nil || res != "OK\n" {
		t.Fatalf("UPDATEINDEX failed: %v, res: %q", err, res)
	}

	res, err = eval.EvaluateQuery(`VALUEBYINDEX items 1`)
	if err != nil || res != "updated_item\n" {
		t.Fatalf("VALUEBYINDEX items 1 failed after update: %v, res: %q", err, res)
	}

	// 4. DELINDEX on list
	res, err = eval.EvaluateQuery(`DELINDEX items 0`)
	if err != nil || res != "OK\n" {
		t.Fatalf("DELINDEX failed: %v, res: %q", err, res)
	}

	// 5. DELFROMLIST on list
	res, err = eval.EvaluateQuery(`DELFROMLIST items "item3"`)
	if err != nil || res != "OK\n" {
		t.Fatalf("DELFROMLIST failed: %v, res: %q", err, res)
	}

	// 6. GET items should remain with [updated_item]
	res, err = eval.EvaluateQuery(`GET items`)
	if err != nil || res != "[updated_item]\n" {
		t.Fatalf("GET items failed: %v, res: %q", err, res)
	}

	// 7. String literal testing with VALUEBYINDEX
	res, err = eval.EvaluateQuery(`SET greeting "Raven"`)
	if err != nil || res != "OK\n" {
		t.Fatalf("SET greeting failed: %v, res: %q", err, res)
	}

	res, err = eval.EvaluateQuery(`VALUEBYINDEX greeting 0`)
	if err != nil || res != "R\n" {
		t.Fatalf("VALUEBYINDEX greeting 0 failed: %v, res: %q", err, res)
	}
}

