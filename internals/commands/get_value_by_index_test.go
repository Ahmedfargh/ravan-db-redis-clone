package commands

import (
	"Raven/internals/database"
	"Raven/internals/parser"
	"testing"
)

func TestValueByIndex_Execute_ArgCount(t *testing.T) {
	database.InitiatDataStore()
	cmd := &ValueByIndex{}

	// Test with 0 args
	res, err := cmd.Execute([]string{})
	if err == nil {
		t.Fatalf("expected error for 0 args, got nil")
	}
	if res != "ValueByIndex Accept Only two args" {
		t.Errorf("expected 'ValueByIndex Accept Only two args', got %q", res)
	}

	// Test with 1 arg
	res, err = cmd.Execute([]string{"key"})
	if err == nil {
		t.Fatalf("expected error for 1 arg, got nil")
	}

	// Test with 3 args
	res, err = cmd.Execute([]string{"key", "1", "extra"})
	if err == nil {
		t.Fatalf("expected error for 3 args, got nil")
	}
}

func TestValueByIndex_Execute_InvalidIndex(t *testing.T) {
	database.InitiatDataStore()
	cmd := &ValueByIndex{}

	// Test with non-numeric index
	res, err := cmd.Execute([]string{"key", "abc"})
	if err == nil {
		t.Fatalf("expected error for non-numeric index, got nil")
	}
	if res != "Invalid index format" {
		t.Errorf("expected 'Invalid index format', got %q", res)
	}
}

func TestValueByIndex_Execute_NonExistentKey(t *testing.T) {
	database.InitiatDataStore()
	cmd := &ValueByIndex{}

	// Test with missing key
	res, err := cmd.Execute([]string{"missing_key", "0"})
	if err == nil {
		t.Fatalf("expected error for missing key, got nil")
	}
	if res != "value don't exists" {
		t.Errorf("expected 'value don't exists', got %q", res)
	}
}

func TestValueByIndex_Execute_StringLiteral(t *testing.T) {
	database.InitiatDataStore()
	cmd := &ValueByIndex{}

	// Store a StringLiteral
	database.Data_store.SetValue("str_key", &parser.StringLiteral{Value: "Raven"}, 0)

	// Test valid index 0
	res, err := cmd.Execute([]string{"str_key", "0"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "R" {
		t.Errorf("expected 'R', got %q", res)
	}

	// Test valid index 4
	res, err = cmd.Execute([]string{"str_key", "4"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "n" {
		t.Errorf("expected 'n', got %q", res)
	}

	// Test out of bounds index (positive)
	res, err = cmd.Execute([]string{"str_key", "5"})
	if err == nil {
		t.Fatalf("expected error for out of bounds positive index, got nil")
	}
	if res != "Index out of range" {
		t.Errorf("expected 'Index out of range', got %q", res)
	}

	// Test out of bounds index (negative)
	res, err = cmd.Execute([]string{"str_key", "-1"})
	if err == nil {
		t.Fatalf("expected error for out of bounds negative index, got nil")
	}
	if res != "Index out of range" {
		t.Errorf("expected 'Index out of range', got %q", res)
	}
}

func TestValueByIndex_Execute_ListLiteral(t *testing.T) {
	database.InitiatDataStore()
	cmd := &ValueByIndex{}

	// Store a ListLiteral [ "apple", 123 ]
	listNode := &parser.ListLiteral{
		Values: []parser.Node{
			&parser.StringLiteral{Value: "apple"},
			&parser.NumberLiteral{Value: "123"},
		},
	}
	database.Data_store.SetValue("list_key", listNode, 0)

	// Test valid index 0
	res, err := cmd.Execute([]string{"list_key", "0"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "apple" {
		t.Errorf("expected 'apple', got %q", res)
	}

	// Test valid index 1
	res, err = cmd.Execute([]string{"list_key", "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "123" {
		t.Errorf("expected '123', got %q", res)
	}

	// Test out of bounds index (positive)
	res, err = cmd.Execute([]string{"list_key", "2"})
	if err == nil {
		t.Fatalf("expected error for out of bounds positive index, got nil")
	}
	if res != "Index out of range" {
		t.Errorf("expected 'Index out of range', got %q", res)
	}

	// Test out of bounds index (negative)
	res, err = cmd.Execute([]string{"list_key", "-1"})
	if err == nil {
		t.Fatalf("expected error for out of bounds negative index, got nil")
	}
	if res != "Index out of range" {
		t.Errorf("expected 'Index out of range', got %q", res)
	}
}

func TestValueByIndex_Execute_UnsupportedType(t *testing.T) {
	database.InitiatDataStore()
	cmd := &ValueByIndex{}

	// Store a raw string value
	database.Data_store.SetValue("raw_str_key", "simple string", 0)

	// Test execution on unsupported type
	res, err := cmd.Execute([]string{"raw_str_key", "0"})
	if err == nil {
		t.Fatalf("expected error for unsupported type, got nil")
	}
	if res != "Element is not a string literal" {
		t.Errorf("expected 'Element is not a string literal', got %q", res)
	}
}

func TestValueByIndex_GlobalRegistry(t *testing.T) {
	database.InitiatDataStore()

	// Store a StringLiteral
	database.Data_store.SetValue("str_key", &parser.StringLiteral{Value: "Raven"}, 0)

	// Execute through registry
	res, err := GlobalRegistry.Execute("VALUEBYINDEX", []string{"str_key", "2"})
	if err != nil {
		t.Fatalf("unexpected registry error: %v", err)
	}
	if res != "v" {
		t.Errorf("expected 'v', got %q", res)
	}
}
