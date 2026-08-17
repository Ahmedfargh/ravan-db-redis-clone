package database

import (
	"sync"
)

type DataValue struct {
	mu    sync.RWMutex
	Value any
	Ttl   uint64
}

func NewDataValueObject(value any, ttl uint64) *DataValue {
	return &DataValue{
		Value: value,
		Ttl:   ttl,
	}
}

// Get safely reads value and ttl under value-level RLock
func (dv *DataValue) Get() (any, uint64) {
	dv.mu.RLock()
	defer dv.mu.RUnlock()
	return dv.Value, dv.Ttl
}

// Set safely updates value and ttl under value-level Lock
func (dv *DataValue) Set(value any, ttl uint64) {
	dv.mu.Lock()
	defer dv.mu.Unlock()
	dv.Value = value
	dv.Ttl = ttl
}

// Lock locks the individual DataValue for exclusive mutation
func (dv *DataValue) Lock() {
	dv.mu.Lock()
}

// Unlock unlocks the individual DataValue
func (dv *DataValue) Unlock() {
	dv.mu.Unlock()
}

// RLock locks the individual DataValue for shared reading
func (dv *DataValue) RLock() {
	dv.mu.RLock()
}

// RUnlock unlocks the individual DataValue after shared reading
func (dv *DataValue) RUnlock() {
	dv.mu.RUnlock()
}
