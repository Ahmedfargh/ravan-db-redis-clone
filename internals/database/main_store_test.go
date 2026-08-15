package database

import (
	"sync"
	"testing"
)

func TestDataStoreConcurrentAccess(t *testing.T) {
	InitiatDataStore()

	var wg sync.WaitGroup
	numRoutines := 50

	// Concurrent writes
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := "key"
			Data_store.SetValue(key, id, 0)
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Data_store.GetValue("key")
		}()
	}

	wg.Wait()

	// Verify key exists
	val, exists := Data_store.GetValue("key")
	if !exists || val == nil {
		t.Fatalf("Expected key to exist in DataStore")
	}
}

func TestGetValueInvertedLogicFix(t *testing.T) {
	InitiatDataStore()

	// Non-existing key should return nil, false (no panic!)
	val, exists := Data_store.GetValue("nonexistent_key")
	if exists {
		t.Errorf("Expected exists=false for missing key")
	}
	if val != nil {
		t.Errorf("Expected val=nil for missing key")
	}

	// Existing key should return val, true
	Data_store.SetValue("testKey", "testVal", 0)
	val, exists = Data_store.GetValue("testKey")
	if !exists {
		t.Errorf("Expected exists=true for existing key")
	}
	if val == nil || val.Value != "testVal" {
		t.Errorf("Expected value 'testVal', got %v", val)
	}
}
