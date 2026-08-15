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

func (data_store *DataStore) GetValue(key string) (*DataValue, bool) {
	data_store.mu.RLock()
	defer data_store.mu.RUnlock()

	value, exists := data_store.Data[key]
	if !exists {
		return nil, false
	}
	return value, true
}

func (data_store *DataStore) SetValue(key string, value any, ttl uint64) (*DataValue, bool) {
	data_store.mu.Lock()
	defer data_store.mu.Unlock()

	data_value := NewDataValueObject(value, ttl)
	data_store.Data[key] = data_value
	return data_value, true
}

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
