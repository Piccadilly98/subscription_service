package cache

import (
	"sync"
	"time"
)

type Cache struct {
	cache map[string]time.Time
	mu    sync.RWMutex
	ttl   time.Duration
}

func NewCache(duration time.Duration) *Cache {
	return &Cache{
		cache: make(map[string]time.Time),
		ttl:   duration,
	}
}

func (c *Cache) CheckBySubID(id string) bool {
	now := time.Now()
	c.mu.RLock()
	t, ok := c.cache[id]
	c.mu.RUnlock()
	if ok {
		if t.Add(c.ttl).After(now) {
			return true
		}
		c.mu.Lock()
		delete(c.cache, id)
		c.mu.Unlock()
	}
	return false
}

func (c *Cache) AddToCacheByID(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[id] = time.Now()
}

func (c *Cache) DeleteByID(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, id)
}
