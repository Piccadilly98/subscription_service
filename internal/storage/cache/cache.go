package cache

import (
	"fmt"
	"sync"
)

type Cache struct {
	cache map[string]bool
	mu    sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		cache: make(map[string]bool),
	}
}

func (c *Cache) CheckBySubID(id string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.cache[id]
}

func (c *Cache) AddToCacheByID(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[id] = true
	fmt.Printf("id: %s, status: %v\n", id, true)
}

func (c *Cache) DeleteByID(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, id)
}
