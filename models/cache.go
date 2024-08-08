package models

import (
	"container/list"
	"lru-cache-gin/utils"
	"sync"
	"time"
)

type CacheItem struct {
	Key        string
	Value      interface{}
	Expiration time.Time
}

type LRUCache struct {
	Capacity int
	Items    map[string]*list.Element
	Order    *list.List
	Mu       sync.RWMutex
}

func NewLRUCache(capacity int) *LRUCache {
	cache := &LRUCache{
		Capacity: capacity,
		Items:    make(map[string]*list.Element),
		Order:    list.New(),
	}
	go cache.startEvictionRoutine()
	return cache
}

func (c *LRUCache) Set(key string, value interface{}, ttl time.Duration, hub *utils.WebSocketHub) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if el, found := c.Items[key]; found {
		c.Order.MoveToFront(el)
		el.Value.(*CacheItem).Value = value
		el.Value.(*CacheItem).Expiration = time.Now().Add(ttl)
	} else {
		if c.Order.Len() >= c.Capacity {
			oldest := c.Order.Back()
			c.Order.Remove(oldest)
			delete(c.Items, oldest.Value.(*CacheItem).Key)
		}
		item := &CacheItem{
			Key:        key,
			Value:      value,
			Expiration: time.Now().Add(ttl),
		}
		c.Items[key] = c.Order.PushFront(item)
	}
	c.BroadcastCacheUpdate(hub)
}

func (c *LRUCache) Get(key string, hub *utils.WebSocketHub) (interface{}, bool) {
	c.Mu.RLock()
	defer c.Mu.RUnlock()

	if el, found := c.Items[key]; found {
		if time.Now().After(el.Value.(*CacheItem).Expiration) {
			return nil, false
		}
		c.Order.MoveToFront(el)
		return el.Value.(*CacheItem).Value, true
	}
	return nil, false
}

func (c *LRUCache) Delete(key string, hub *utils.WebSocketHub) {
	c.Mu.Lock()
	defer c.Mu.Unlock()

	if el, found := c.Items[key]; found {
		c.Order.Remove(el)
		delete(c.Items, key)
		c.BroadcastCacheUpdate(hub)
	}
}

func (c *LRUCache) BroadcastCacheUpdate(hub *utils.WebSocketHub) {
	cacheSnapshot := make(map[string]interface{})
	for k, v := range c.Items {
		item := v.Value.(*CacheItem)
		if time.Now().Before(item.Expiration) {
			remainingSeconds := int(item.Expiration.Sub(time.Now()).Seconds())
			cacheSnapshot[k] = map[string]interface{}{
				"value":      item.Value,
				"expiration": remainingSeconds,
			}
		}
	}

	hub.Broadcast <- cacheSnapshot
}

func (c *LRUCache) startEvictionRoutine() {
	for {
		time.Sleep(1 * time.Second)
		c.Mu.Lock()
		for k, el := range c.Items {
			if time.Now().After(el.Value.(*CacheItem).Expiration) {
				c.Order.Remove(el)
				delete(c.Items, k)
			}
		}
		c.Mu.Unlock()
	}
}
