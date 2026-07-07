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
	entry, is := c.m[key]
	if !is || entry.expiresAt.After(time.Now()) {
		defer c.mu.RUnlock()
		if !is {
			return nil, time.Time{}, false
		}
		return entry.value, entry.fetchedAt, true
	}
	c.mu.RUnlock()
	c.mu.Lock()
	defer c.mu.Unlock()
	if entry, is = c.m[key]; is && !entry.expiresAt.After(time.Now()) {
		delete(c.m, key)
	}
	return nil, time.Time{}, false
}

func (c *TTLCache) Set(key string, value any, ttl time.Duration) {
	entry := entry{value: value, fetchedAt: time.Now(), expiresAt: time.Now().Add(ttl)}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.m[key] = entry
}
