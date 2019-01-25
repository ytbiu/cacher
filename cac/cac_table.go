package cac

import (
	"sync"
	"time"
)

type cacTable struct {
	name string
	sync.RWMutex
	key2Cache          map[string]*cac
	nextExpireCac      *cac
	checkIntervalTimer *time.Timer
}

func NewCacTable(name string) *cacTable {
	return &cacTable{
		name:      name,
		key2Cache: make(map[string]*cac),
	}
}

func (c *cacTable) Add(key string, val interface{}, expire time.Duration, cbs ...func()) {
	c.RLock()
	cache, find := c.key2Cache[key]
	c.RUnlock()
	if find {
		c.Lock()
		if cache.val != val {
			cache.put(val)
		}
		cache.reset(expire)
		c.Unlock()
		return
	}
	newCache := newCac(key, val, expire, c.Delete, cbs...)
	c.Lock()
	c.key2Cache[key] = newCache
	c.Unlock()

}

func (c *cacTable) Delete(key string) {
	c.Lock()
	defer c.Unlock()

	delete(c.key2Cache, key)
}

func (c *cacTable) Get(key string) (interface{}, bool) {
	c.RLock()
	cache, find := c.key2Cache[key]
	c.RUnlock()

	if find {
		return cache.val, true
	}
	return nil, false
}

func (c *cacTable) Reset(key string, expire time.Duration) bool {
	c.RLock()
	cache, find := c.key2Cache[key]
	c.RUnlock()

	if find {
		c.Lock()
		cache.reset(expire)
		c.Unlock()
	}
	return find
}
