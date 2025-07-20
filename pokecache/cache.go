package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entry map[string]cacheEntry
	mu    sync.Mutex
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		entry: make(map[string]cacheEntry),
	}

	go c.reapLoop(interval)

	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entry[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, ok := c.entry[key]
	if !ok {
		return nil, false
	}

	// Create a copy of the value to prevent external modifications
	valCopy := make([]byte, len(data.val))
	copy(valCopy, data.val)

	return valCopy, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.mu.Lock()
		for key, data := range c.entry {
			if time.Since(data.createdAt) > interval {
				delete(c.entry, key)
			}
		}
		c.mu.Unlock()
	}
}
