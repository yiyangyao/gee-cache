package cache

import (
	"gee-cache/src/byteview"
	"gee-cache/src/lru"
	"sync"
)

type Cache struct {
	mutex      sync.Mutex
	lru        *lru.Cache
	CacheBytes int64
}

func (c *Cache) Add(key string, value byteview.ByteView) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.CacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *Cache) Get(key string) (value byteview.ByteView, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(byteview.ByteView), ok
	}

	return
}
