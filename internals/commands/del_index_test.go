package commands

import (
	"Raven/internals/database"
	"Raven/internals/parser"
	"testing"
)

func TestDeleteValueFromIndex_Execute_ArgCount(t *testing.T) {
	database.InitiatDataStore()
	cmd := &DeleteValueFromIndex{}

	// Test with 0 args
	_, err := cmd.Execute([]string{})
	if err == nil {
		t.Fatalf("expected error for 0 args, got nil")
	}

	// Test with 1 arg
	_, err = cmd.Execute([]string{"key"})
	if err == nil {
		t.Fatalf("expected error for 1 arg, got nil")
	}

	// Test with 3 args
	_, err = cmd.Execute([]string{"key", "1", "extra"})
	if err == nil {
		t.Fatalf("expected error for 3 args, got nil")
	}
}

func TestDeleteValueFromIndex_Execute_InvalidIndex(t *testing.T) {
	database.InitiatDataStore()
	cmd := &DeleteValueFromIndex{}

	// Test with non-numeric index
	res, err := cmd.Execute([]string{"key", "abc"})
	if err == nil {
		t.Fatalf("expected error for non-numeric index, got nil")
	}
	if res != "Invalid index format" {
		t.Errorf("expected 'Invalid index format', got %q", res)
	}
}

func TestDeleteValueFromIndex_Execute_NonExistentKey(t *testing.T) {
	database.InitiatDataStore()
	cmd := &DeleteValueFromIndex{}

	// Test with missing key
	res, err := cmd.Execute([]string{"missing_key", "0"})
	if err == nil {
		t.Fatalf("expected error for missing key, got nil")
	}
	if res != "value don't exists" {
		t.Errorf("expected 'value don't exists', got %q", res)
	}
}

func TestDeleteValueFromIndex_Execute_StringLiteral(t *testing.T) {
	database.InitiatDataStore()
	cmd := &DeleteValueFromIndex{}

	// Store a StringLiteral "Raven"
	strNode := &parser.StringLiteral{Value: "Raven"}
	database.Data_store.SetValue("str_key", strNode, 0)

	// Test valid index 2 to delete 'v' (Raven -> Raen)
	res, err := cmd.Execute([]string{"str_key", "2"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "OK\n" {
		t.Errorf("expected 'OK\n', got %q", res)
	}
	if strNode.Value != "Raen" {
		t.Errorf("expected 'Raen', got %q", strNode.Value)
	}

	// Test delete first char (index 0: Raen -> aen)
	res, err = cmd.Execute([]string{"str_key", "0"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strNode.Value != "aen" {
		t.Errorf("expected 'aen', got %q", strNode.Value)
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

func TestDeleteValueFromIndex_Execute_ListLiteral(t *testing.T) {
	database.InitiatDataStore()
	cmd := &DeleteValueFromIndex{}

	// Store a ListLiteral [ "apple", "banana", 123 ]
	listNode := &parser.ListLiteral{
		Values: []parser.Node{
			&parser.StringLiteral{Value: "apple"},
			&parser.StringLiteral{Value: "banana"},
			&parser.NumberLiteral{Value: "123"},
		},
	}
	database.Data_store.SetValue("list_key", listNode, 0)

	// Test delete index 1 ("banana")
	res, err := cmd.Execute([]string{"list_key", "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "OK\n" {
		t.Errorf("expected 'OK\n', got %q", res)
	}

	// Verify length is 2 and contents are "apple" and 123
	if len(listNode.Values) != 2 {
		t.Fatalf("expected 2 elements in list, got %d", len(listNode.Values))
	}
	if listNode.Values[0].String() != "apple" {
		t.Errorf("expected element 0 to be 'apple', got %q", listNode.Values[0].String())
	}
	if listNode.Values[1].String() != "123" {
		t.Errorf("expected element 1 to be '123', got %q", listNode.Values[1].String())
	}

	// Test delete index 0 ("apple")
	res, err = cmd.Execute([]string{"list_key", "0"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(listNode.Values) != 1 || listNode.Values[0].String() != "123" {
		t.Fatalf("expected 1 element '123', got %v", listNode.Values)
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

func TestDeleteValueFromIndex_Execute_UnsupportedType(t *testing.T) {
	database.InitiatDataStore()
	cmd := &DeleteValueFromIndex{}

	// Store a raw string value
	database.Data_store.SetValue("raw_str_key", "simple string", 0)

	// Test execution on unsupported type
	res, err := cmd.Execute([]string{"raw_str_key", "0"})
	if err == nil {
		t.Fatalf("expected error for unsupported type, got nil")
	}
	if res != "Element is not a string or list literal" {
		t.Errorf("expected 'Element is not a string or list literal', got %q", res)
	}
}

func TestDeleteValueFromIndex_GlobalRegistry(t *testing.T) {
	database.InitiatDataStore()

	// Store a StringLiteral
	strNode := &parser.StringLiteral{Value: "Raven"}
	database.Data_store.SetValue("str_key", strNode, 0)

	// Execute through registry using DELINDEX
	res, err := GlobalRegistry.Execute("DELINDEX", []string{"str_key", "2"})
	if err != nil {
		t.Fatalf("unexpected registry error: %v", err)
	}
	if res != "OK\n" {
		t.Errorf("expected 'OK\n', got %q", res)
	}
	if strNode.Value != "Raen" {
		t.Errorf("expected 'Raen', got %q", strNode.Value)
	}

	// Execute through registry using DELETEINDEX alias
	res, err = GlobalRegistry.Execute("DELETEINDEX", []string{"str_key", "0"})
	if err != nil {
		t.Fatalf("unexpected registry error: %v", err)
	}
	if res != "OK\n" {
		t.Errorf("expected 'OK\n', got %q", res)
	}
	if strNode.Value != "aen" {
		t.Errorf("expected 'aen', got %q", strNode.Value)
	}
}
