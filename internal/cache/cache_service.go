package cache

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/edaywalid/reverse-proxy/pkg/utils"
)

type CacheService struct {
	cache *Cache
}

func NewCacheService(cache *Cache) *CacheService {
	cs := &CacheService{
		cache: cache,
	}
	go cs.cache.StartCleanUp(utils.CleanupInterval)
	return cs
}

func (s *CacheService) Get(key string) (*CacheItem, bool) {

	item, ok := s.cache.GetItem(key)
	if !ok || time.Since(item.Created()) > utils.CacheExpiration {
		return nil, false
	}

	return &item, ok
}

func (s *CacheService) Set(key string, resp *http.Response) {
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	item := NewCacheItem(body, resp.Header, resp.StatusCode)
	s.cache.SetItem(key, *item)

	resp.Body = io.NopCloser(bytes.NewReader(body))
}
