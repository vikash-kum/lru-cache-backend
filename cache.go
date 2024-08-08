package main

import (
	"container/list"
	"sync"
	"time"
)

type CacheItem struct {
	key        string
	value      interface{}
	expiration time.Time
}

type LRUCache struct {
	capacity int
	items    map[string]*list.Element
	order    *list.List
	mu       sync.RWMutex
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		order:    list.New(),
	}
}

func (c *LRUCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, found := c.items[key]; found {
		c.order.MoveToFront(el)
		el.Value.(*CacheItem).value = value
		el.Value.(*CacheItem).expiration = time.Now().Add(ttl)
	} else {
		if c.order.Len() >= c.capacity {
			oldest := c.order.Back()
			c.order.Remove(oldest)
			delete(c.items, oldest.Value.(*CacheItem).key)
		}
		item := &CacheItem{
			key:        key,
			value:      value,
			expiration: time.Now().Add(ttl),
		}
		c.items[key] = c.order.PushFront(item)
	}
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if el, found := c.items[key]; found {
		if time.Now().After(el.Value.(*CacheItem).expiration) {
			c.mu.RUnlock()
			c.mu.Lock()
			c.order.Remove(el)
			delete(c.items, key)
			c.mu.Unlock()
			c.mu.RLock()
			return nil, false
		}
		c.order.MoveToFront(el)
		return el.Value.(*CacheItem).value, true
	}
	return nil, false
}

func (c *LRUCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, found := c.items[key]; found {
		c.order.Remove(el)
		delete(c.items, key)
	}
}
