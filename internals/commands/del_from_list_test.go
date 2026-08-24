package commands

import (
	"Raven/internals/database"
	"Raven/internals/parser"
	"testing"
)

func TestDeleteFromList_Execute_ArgCount(t *testing.T) {
	database.InitiatDataStore()
	cmd := &DeleteFromList{}

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
}

func TestDeleteFromList_Execute_NonExistentKey(t *testing.T) {
	database.InitiatDataStore()
	cmd := &DeleteFromList{}

	// Test with missing key
	res, err := cmd.Execute([]string{"missing_key", "banana"})
	if err == nil {
		t.Fatalf("expected error for missing key, got nil")
	}
	if res != "value don't exists" {
		t.Errorf("expected 'value don't exists', got %q", res)
	}
}

func TestDeleteFromList_Execute_ListLiteral(t *testing.T) {
	database.InitiatDataStore()
	cmd := &DeleteFromList{}

	// Store a ListLiteral [ "apple", "banana", "orange", "banana", 123 ]
	listNode := &parser.ListLiteral{
		Values: []parser.Node{
			&parser.StringLiteral{Value: "apple"},
			&parser.StringLiteral{Value: "banana"},
			&parser.StringLiteral{Value: "orange"},
			&parser.StringLiteral{Value: "banana"},
			&parser.NumberLiteral{Value: "123"},
		},
	}
	database.Data_store.SetValue("list_key", listNode, 0)

	// Delete all occurrences of "banana"
	res, err := cmd.Execute([]string{"list_key", "banana"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "OK\n" {
		t.Errorf("expected 'OK\n', got %q", res)
	}

	// Verify remaining elements: [ "apple", "orange", 123 ]
	if len(listNode.Values) != 3 {
		t.Fatalf("expected 3 elements in list, got %d", len(listNode.Values))
	}
	if listNode.Values[0].String() != "apple" {
		t.Errorf("expected element 0 to be 'apple', got %q", listNode.Values[0].String())
	}
	if listNode.Values[1].String() != "orange" {
		t.Errorf("expected element 1 to be 'orange', got %q", listNode.Values[1].String())
	}
	if listNode.Values[2].String() != "123" {
		t.Errorf("expected element 2 to be '123', got %q", listNode.Values[2].String())
	}

	// Delete multiple items at once ("apple", "123")
	res, err = cmd.Execute([]string{"list_key", "apple", "123"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(listNode.Values) != 1 || listNode.Values[0].String() != "orange" {
		t.Fatalf("expected 1 element 'orange', got %v", listNode.Values)
	}
}

func TestDeleteFromList_Execute_StringLiteral(t *testing.T) {
	database.InitiatDataStore()
	cmd := &DeleteFromList{}

	// Store a StringLiteral "Hello World"
	strNode := &parser.StringLiteral{Value: "Hello World"}
	database.Data_store.SetValue("str_key", strNode, 0)

	// Delete "l" -> "Heo Word"
	res, err := cmd.Execute([]string{"str_key", "l"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != "OK\n" {
		t.Errorf("expected 'OK\n', got %q", res)
	}
	if strNode.Value != "Heo Word" {
		t.Errorf("expected 'Heo Word', got %q", strNode.Value)
	}
}

func TestDeleteFromList_Execute_UnsupportedType(t *testing.T) {
	database.InitiatDataStore()
	cmd := &DeleteFromList{}

	// Store a raw string value
	database.Data_store.SetValue("raw_str_key", "simple string", 0)

	// Test execution on unsupported type
	res, err := cmd.Execute([]string{"raw_str_key", "simple"})
	if err == nil {
		t.Fatalf("expected error for unsupported type, got nil")
	}
	if res != "Element is not a string or list literal" {
		t.Errorf("expected 'Element is not a string or list literal', got %q", res)
	}
}

func TestDeleteFromList_GlobalRegistry(t *testing.T) {
	database.InitiatDataStore()

	// Store a ListLiteral
	listNode := &parser.ListLiteral{
		Values: []parser.Node{
			&parser.StringLiteral{Value: "item1"},
			&parser.StringLiteral{Value: "item2"},
		},
	}
	database.Data_store.SetValue("items_key", listNode, 0)

	// Execute through registry using DELFROMLIST
	res, err := GlobalRegistry.Execute("DELFROMLIST", []string{"items_key", "item1"})
	if err != nil {
		t.Fatalf("unexpected registry error: %v", err)
	}
	if res != "OK\n" {
		t.Errorf("expected 'OK\n', got %q", res)
	}
	if len(listNode.Values) != 1 || listNode.Values[0].String() != "item2" {
		t.Fatalf("expected 1 element 'item2', got %v", listNode.Values)
	}
}
