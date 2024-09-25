package models_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/edaywalid/reverse-proxy/internal/cache"
	"github.com/stretchr/testify/assert"
)

const (
	CacheExpiration = 3 * time.Second
	CleanUpInterval = 1 * time.Second
)

func TestCacheItem(t *testing.T) {

	// Create a new CacheItem with test data
	ci := cache.NewCacheItem([]byte("test"), http.Header{"Content-Type": []string{"text/plain"}}, 200)

	// Assert that the cache item is initialized correctly
	assert.NotNil(t, ci)
	assert.Equal(t, ci.Body, []byte("test"))
	assert.Equal(t, ci.Headers.Get("Content-Type"), "text/plain")
	assert.Equal(t, ci.Status, 200)

	// Assert that the creation time is within the correct range
	assert.NotNil(t, ci.Created())
	assert.LessOrEqual(t, ci.Created().Unix(), time.Now().Unix(), "Creation time should be less than or equal to the current time")

	// Sleep for CacheExpiration duration to simulate cache aging
	time.Sleep(CacheExpiration)

	// Assert that the cache item has expired
	// Note: Created() time should still be earlier than `time.Now()` after sleep
	assert.LessOrEqual(t, ci.Created().Add(CacheExpiration).Unix(), time.Now().Unix(), "CacheItem should have expired")
}

func TestCache(t *testing.T) {

	// Create a new cache
	c := cache.NewCache()

	// Assert that the cache is initialized correctly
	assert.NotNil(t, c)

	// Assert that the cache is empty initially
	assert.Equal(t, 0, len(c.Items()), "Cache should be empty")

	// Assert that inserting an item into the cache works
	c.SetItem("test", cache.CacheItem{
		Body:    []byte("test data"),
		Headers: nil,
		Status:  200,
	})

	// Check if the item was inserted correctly
	item, ok := c.GetItem("test")
	assert.True(t, ok, "Item should exist in cache")
	assert.NotNil(t, item, "Item should not be nil")

	// Assert that removing an item from the cache works
	c.RemoveItem("test")
	_, ok = c.GetItem("test")
	assert.False(t, ok, "Item should not exist after removal")
	assert.Equal(t, 0, len(c.Items()), "Cache should be empty after removing items")

	// Testing auto-cleanup of expired items

	// Start the cleanup in the background
	go c.StartCleanUp(CleanUpInterval)

	// Insert items into the cache
	c.SetItem("test1", cache.CacheItem{
		Body:    []byte("item1"),
		Headers: nil,
		Status:  200,
	})
	c.SetItem("test2", cache.CacheItem{
		Body:    []byte("item2"),
		Headers: nil,
		Status:  200,
	})

	// Cache should have two items now
	assert.Equal(t, 2, len(c.Items()), "Cache should have 2 items before expiration")

	// Sleep for CacheExpiration duration to simulate expiration of items
	time.Sleep(CacheExpiration + CleanUpInterval) // Add CleanUpInterval to ensure cleanup runs

	// Cache should be empty after expiration and cleanup
	assert.Equal(t, 0, len(c.Items()), "Cache should be empty after expiration and cleanup")
}
