package commands

import (
	"Raven/internals/database"
	"Raven/internals/parser"
	"testing"
)

func TestUpdateListIndexValue_Execute_ArgCount(t *testing.T) {
	database.InitiatDataStore()
	cmd := &UpdateListIndexValue{}

	// Test with 0 args
	_, err := cmd.Execute([]string{})
	if err == nil {
		t.Fatalf("expected error for 0 args, got nil")
	}

	// Test with 2 args
	_, err = cmd.Execute([]string{"key", "1"})
	if err == nil {
		t.Fatalf("expected error for 2 args, got nil")
	}

	// Test with 4 args
	_, err = cmd.Execute([]string{"key", "1", "val", "extra"})
	if err == nil {
		t.Fatalf("expected error for 4 args, got nil")
	}
}

func TestUpdateListIndexValue_Execute_InvalidIndex(t *testing.T) {
	database.InitiatDataStore()
	cmd := &UpdateListIndexValue{}

	// Test with non-numeric index
	res, err := cmd.Execute([]string{"key", "abc", "newval"})
	if err == nil {
		t.Fatalf("expected error for non-numeric index, got nil")
	}
	if res != "Invalid index format" {
		t.Errorf("expected 'Invalid index format', got %q", res)
	}
}

func TestUpdateListIndexValue_Execute_NonExistentKey(t *testing.T) {
	database.InitiatDataStore()
	cmd := &UpdateListIndexValue{}

	// Test with missing key
	res, err := cmd.Execute([]string{"missing_key", "0", "newval"})
	if err == nil {
		t.Fatalf("expected error for missing key, got nil")
	}
	if res != "value don't exists" {
		t.Errorf("expected 'value don't exists', got %q", res)
	}
}

func TestUpdateListIndexValue_Execute_StringLiteral(t *testing.T) {
	database.InitiatDataStore()
	cmd := &UpdateListIndexValue{}

	// Store a StringLiteral
	strNode := &parser.StringLiteral{Value: "Raven"}
	database.Data_store.SetValue("str_key", strNode, 0)

	// Test valid index 2 to change 'v' to 'X' (Raven -> RaXen)
	res, err := cmd.Execute([]string{"str_key", "2", "X"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "OK\n" {
		t.Errorf("expected 'OK\n', got %q", res)
	}
	if strNode.Value != "RaXen" {
		t.Errorf("expected 'RaXen', got %q", strNode.Value)
	}

	// Test out of bounds index (positive)
	res, err = cmd.Execute([]string{"str_key", "5", "X"})
	if err == nil {
		t.Fatalf("expected error for out of bounds positive index, got nil")
	}
	if res != "Index out of range" {
		t.Errorf("expected 'Index out of range', got %q", res)
	}

	// Test out of bounds index (negative)
	res, err = cmd.Execute([]string{"str_key", "-1", "X"})
	if err == nil {
		t.Fatalf("expected error for out of bounds negative index, got nil")
	}
	if res != "Index out of range" {
		t.Errorf("expected 'Index out of range', got %q", res)
	}
}

func TestUpdateListIndexValue_Execute_ListLiteral(t *testing.T) {
	database.InitiatDataStore()
	cmd := &UpdateListIndexValue{}

	// Store a ListLiteral [ "apple", 123 ]
	listNode := &parser.ListLiteral{
		Values: []parser.Node{
			&parser.StringLiteral{Value: "apple"},
			&parser.NumberLiteral{Value: "123"},
		},
	}
	database.Data_store.SetValue("list_key", listNode, 0)

	// Test update index 0 to "banana" (string literal)
	res, err := cmd.Execute([]string{"list_key", "0", "banana"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "OK\n" {
		t.Errorf("expected 'OK\n', got %q", res)
	}

	// Verify the element at index 0 is updated to a StringLiteral
	strNode, ok := listNode.Values[0].(*parser.StringLiteral)
	if !ok {
		t.Fatalf("expected updated element to be *parser.StringLiteral, got %T", listNode.Values[0])
	}
	if strNode.Value != "banana" {
		t.Errorf("expected 'banana', got %q", strNode.Value)
	}

	// Test update index 1 to "456" (number literal)
	res, err = cmd.Execute([]string{"list_key", "1", "456"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "OK\n" {
		t.Errorf("expected 'OK\n', got %q", res)
	}

	// Verify the element at index 1 is updated to a NumberLiteral
	numNode, ok := listNode.Values[1].(*parser.NumberLiteral)
	if !ok {
		t.Fatalf("expected updated element to be *parser.NumberLiteral, got %T", listNode.Values[1])
	}
	if numNode.Value != "456" {
		t.Errorf("expected '456', got %q", numNode.Value)
	}

	// Test out of bounds index (positive)
	res, err = cmd.Execute([]string{"list_key", "2", "orange"})
	if err == nil {
		t.Fatalf("expected error for out of bounds positive index, got nil")
	}
	if res != "Index out of range" {
		t.Errorf("expected 'Index out of range', got %q", res)
	}

	// Test out of bounds index (negative)
	res, err = cmd.Execute([]string{"list_key", "-1", "orange"})
	if err == nil {
		t.Fatalf("expected error for out of bounds negative index, got nil")
	}
	if res != "Index out of range" {
		t.Errorf("expected 'Index out of range', got %q", res)
	}
}

func TestUpdateListIndexValue_Execute_UnsupportedType(t *testing.T) {
	database.InitiatDataStore()
	cmd := &UpdateListIndexValue{}

	// Store a raw string value
	database.Data_store.SetValue("raw_str_key", "simple string", 0)

	// Test execution on unsupported type
	res, err := cmd.Execute([]string{"raw_str_key", "0", "X"})
	if err == nil {
		t.Fatalf("expected error for unsupported type, got nil")
	}
	if res != "Element is not a string or list literal" {
		t.Errorf("expected 'Element is not a string or list literal', got %q", res)
	}
}

func TestUpdateListIndexValue_GlobalRegistry(t *testing.T) {
	database.InitiatDataStore()

	// Store a StringLiteral
	strNode := &parser.StringLiteral{Value: "Raven"}
	database.Data_store.SetValue("str_key", strNode, 0)

	// Execute through registry
	res, err := GlobalRegistry.Execute("UPDATEINDEX", []string{"str_key", "2", "X"})
	if err != nil {
		t.Fatalf("unexpected registry error: %v", err)
	}
	if res != "OK\n" {
		t.Errorf("expected 'OK\n', got %q", res)
	}
	if strNode.Value != "RaXen" {
		t.Errorf("expected 'RaXen', got %q", strNode.Value)
	}
}
