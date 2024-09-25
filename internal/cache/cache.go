package cache

import (
	"net/http"
	"sync"
	"time"

	"github.com/edaywalid/reverse-proxy/pkg/utils"
)

type CacheItem struct {
	/*
		store the response body as a byte slice
		because the response body is an io.ReadCloser
		and it can be read only once
	*/
	Body    []byte
	Headers http.Header
	Status  int
	created time.Time
}

type Cache struct {
	items map[string]CacheItem
	lock  sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}

func NewCacheItem(body []byte, headers http.Header, status int) *CacheItem {
	return &CacheItem{
		Body:    body,
		Headers: headers,
		Status:  status,
		created: time.Now(),
	}
}

func (ci *CacheItem) Created() time.Time {
	return ci.created
}

func (c *Cache) GetItem(key string) (CacheItem, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	item, ok := c.items[key]
	return item, ok
}

func (c *Cache) SetItem(key string, item CacheItem) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	c.items[key] = item
}

func (c *Cache) RemoveItem(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.items, key)
}

func (c *Cache) cleanUpExpired() {
	c.lock.Lock()
	defer c.lock.Unlock()

	for key, item := range c.items {
		if time.Since(item.Created()) > utils.CacheExpiration {
			delete(c.items, key)
		}
	}
}

func (c *Cache) StartCleanUp(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)
			c.cleanUpExpired()
		}
	}()
}

func (c *Cache) Items() map[string]CacheItem {
	return c.items
}
