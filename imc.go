package imc

import (
	"sync"
)

type (
	Key   = string
	Value = string
)

type Cache interface {
	GetOrSet(key Key, valueFn func() Value) Value
	Get(key Key) (Value, bool)
}

type InMemoryCache struct {
	dataMutex sync.RWMutex
	data      map[Key]Value
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[Key]Value),
	}
}

func (cache *InMemoryCache) Get(key Key) (Value, bool) {
	cache.dataMutex.RLock()
	defer cache.dataMutex.RUnlock()

	value, found := cache.data[key]

	return value, found
}

func (cache *InMemoryCache) GetOrSet(key Key, valueFn func() Value) Value {
	value, found := cache.Get(key)
	if found {
		return value
	}

	cache.dataMutex.Lock()
	defer cache.dataMutex.Unlock()
	cache.data[key] = valueFn()

	return cache.data[key]
}
