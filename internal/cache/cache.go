package cache

import (
	"sync"
	"time"
)

type entry struct {
	value     any
	expiresAt time.Time
	fetchedAt time.Time
}

type TTLCache struct {
	m  (map[string]entry)
	mu sync.RWMutex
}

func New() *TTLCache {
	cache := TTLCache{}
	cache.m = make(map[string]entry)
	return &cache
}

func (c *TTLCache) Get(key string) (value any, fetchedAt time.Time, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, is := c.m[key]
	if !is {
		return nil, time.Time{}, false
	}
	if entry.expiresAt.After(time.Now()) {
		c.mu.RLock()
		defer c.mu.RUnlock()
		if entry == c.m[key] {
			return entry.value, entry.fetchedAt, true
		}
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.m, key)
	return nil, time.Time{}, false
}

func (c *TTLCache) Set(key string, value any, ttl time.Duration) {
	entry := entry{value: value, fetchedAt: time.Now(), expiresAt: time.Now().Add(ttl)}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = entry
}
