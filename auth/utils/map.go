package utils

import (
	"sync"

	"github.com/IbraheemHaseeb7/pubsub"
)

type SafeMap struct {
	mu    sync.RWMutex
	data  map[string]any
}

// NewSafeMap creates and returns a new SafeMap.
func NewSafeMap() *SafeMap {
	return &SafeMap{
		data: make(map[string]any),
	}
}

// Store adds or updates a key-value pair in the map.
func (sm *SafeMap) Store(key string, value chan pubsub.PubsubMessage) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = value
}

// StoreAny adds or updates a key-value pair in the map.
func (sm *SafeMap) StoreAny(key string, value any) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.data[key] = value
}

// Load retrieves a value from the map.
func (sm *SafeMap) Load(key string) chan pubsub.PubsubMessage {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	val, _ := sm.data[key]
	finalValue, ok := val.(chan pubsub.PubsubMessage)
	if !ok {
		return nil
	}
	return finalValue
}

func (sm *SafeMap) Delete(key string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.data, key)
}

