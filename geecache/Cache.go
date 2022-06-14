package geecache

import (
	lru2 "DisCache/lru"
	"sync"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru2.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru2.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) Get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Lock()
	if c.lru == nil {
		return ByteView{}, false
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
