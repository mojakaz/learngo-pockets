package genericcache

import (
	"slices"
	"sync"
	"time"
)

// Cache is key-value storage.
type Cache[K comparable, V any] struct {
	data              map[K]entryWithTimeout[V]
	mu                sync.Mutex
	ttl               time.Duration
	maxSize           int
	chronologicalKeys []K
}

// New creates a usable Cache.
func New[K comparable, V any](maxSize int, ttl time.Duration) Cache[K, V] {
	return Cache[K, V]{
		data:              make(map[K]entryWithTimeout[V]),
		ttl:               ttl,
		maxSize:           maxSize,
		chronologicalKeys: make([]K, 0, maxSize),
	}
}

// Read returns the associated value for a key,
// and a boolean to true if the key is absent.
func (c *Cache[K, V]) Read(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zeroV V

	e, ok := c.data[key]

	switch {
	case !ok:
		return zeroV, false
	case e.expires.Before(time.Now()):
		// The value has expired.
		c.deleteKeyValue(key)
		return zeroV, false
	default:
		return e.value, true
	}
}

// Upsert overrides the value for a given key.
func (c *Cache[K, V]) Upsert(key K, value V) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, alreadyPresent := c.data[key]
	switch {
	case alreadyPresent:
		c.deleteKeyValue(key)
	case len(c.data) == c.maxSize:
		c.deleteKeyValue(c.chronologicalKeys[0])
	}
	c.addKeyValue(key, value)
	// Do not return an error for the moment,
	// but it can happen in the near future.
	return nil
}

// Delete removes the entry for the given key.
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.deleteKeyValue(key)
}

type entryWithTimeout[V any] struct {
	value   V
	expires time.Time // After that time, the value is useless.
}

// addKeyValue inserts a key and its value into the cache.
func (c *Cache[K, V]) addKeyValue(key K, value V) {
	c.data[key] = entryWithTimeout[V]{
		value:   value,
		expires: time.Now().Add(c.ttl),
	}
	c.chronologicalKeys = append(c.chronologicalKeys, key)
}

// deleteKeyValue removes a key and its associated value from the cache.
func (c *Cache[K, V]) deleteKeyValue(key K) {
	c.chronologicalKeys = slices.DeleteFunc(c.chronologicalKeys, func(k K) bool { return k == key })
	delete(c.data, key)
}
