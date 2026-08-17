package database

import (
	"sync"
)

type DataStore struct {
	mu   sync.RWMutex
	Data map[string]*DataValue
}

var Data_store *DataStore

func InitiatDataStore() {
	Data_store = &DataStore{
		Data: make(map[string]*DataValue),
	}
}

// GetValue looks up the DataValue pointer from the map
func (data_store *DataStore) GetValue(key string) (*DataValue, bool) {
	data_store.mu.RLock()
	valObj, exists := data_store.Data[key]
	data_store.mu.RUnlock()

	if !exists {
		return nil, false
	}
	return valObj, true
}

// SetValue updates an existing key using per-value locking, or inserts a new key
func (data_store *DataStore) SetValue(key string, value any, ttl uint64) (*DataValue, bool) {
	// Fast-path: Check if key already exists under shared RLock
	data_store.mu.RLock()
	existingVal, exists := data_store.Data[key]
	data_store.mu.RUnlock()

	if exists {
		// Mutate value at the individual DataValue lock level without blocking other keys!
		existingVal.Set(value, ttl)
		return existingVal, true
	}

	// Slow-path: Key is new, acquire map lock to insert entry
	data_store.mu.Lock()
	defer data_store.mu.Unlock()

	// Double check after acquiring write lock
	if existingVal, exists := data_store.Data[key]; exists {
		existingVal.Set(value, ttl)
		return existingVal, true
	}

	newVal := NewDataValueObject(value, ttl)
	data_store.Data[key] = newVal
	return newVal, true
}

// DelValue removes a key entry from the map
func (data_store *DataStore) DelValue(key string) bool {
	data_store.mu.Lock()
	defer data_store.mu.Unlock()

	_, exists := data_store.Data[key]
	if exists {
		delete(data_store.Data, key)
		return true
	}
	return false
}
