package cache

import "time"

type entry struct {
	value     any
	expiresAt time.Time
	fetchedAt time.Time
}

type TTLCache struct {
	m (map[string]entry)
}

func New() *TTLCache {
	cache := TTLCache{}
	cache.m = make(map[string]entry)
	return &cache
}

func (c *TTLCache) Get(key string) (value any, fetchedAt time.Time, ok bool) {
	entry, is := c.m[key]
	if !is {
		return nil, time.Time{}, false
	}
	if entry.expiresAt.After(time.Now()) {
		return entry.value, entry.fetchedAt, true
	}
	delete(c.m, key)
	return nil, time.Time{}, false
}

func (c *TTLCache) Set(key string, value any, ttl time.Duration) {
	entry := entry{value: value, fetchedAt: time.Now(), expiresAt: time.Now().Add(ttl)}
	c.m[key] = entry
}
