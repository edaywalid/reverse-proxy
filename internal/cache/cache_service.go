package cache

import (
	"io"
	"net/http"
	"time"

	"github.com/edaywalid/reverse-proxy/pkg/utils"
)

type cacheService struct {
	cache *Cache
}

func NewCacheService(cache *Cache) *cacheService {
	return &cacheService{
		cache: cache,
	}
}

func (s *cacheService) Get(key string) (*CacheItem, bool) {

	item, ok := s.cache.GetItem(key)
	if !ok || time.Since(item.Created()) > utils.CacheExpiration {
		return nil, false
	}

	return &item, ok
}

func (s *cacheService) Set(key string, resp *http.Response) {
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	item := NewCacheItem(body, resp.Header, resp.StatusCode)
	s.cache.SetItem(key, *item)

	resp.Body = io.NopCloser(bytes.NewReader(body))
}
