package core

import (
	"fmt"
)

type TheCache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) (bool, error)
	IsSet(key string) (bool, error)
	Remove(key string) (bool, error)
	RemoveByPattern(pattern string) (bool, error)
}

// CacheManager manages cache
type CacheManager struct {
	Cache TheCache
}

type acquire func() interface{}

// Get gets or sets a value in cache
func (c *CacheManager) Get(key string, expiresInSeconds int, acquire acquire) interface{} {
	isSet, err := c.Cache.IsSet(key)
	if err != nil {
		fmt.Println("Cache is set error")
	}
	if isSet {
		cacheResult, cacheErr := c.Cache.Get(key)
		if cacheErr != nil {
			fmt.Println("Cache is set error")
		}
		return cacheResult
	}
	result := acquire()

	c.Cache.Set(key, result)
	return result
}

// NewCacheManager initiates new instance for the cache manager
func NewCacheManager(theCache TheCache) *CacheManager {
	return &CacheManager{Cache: theCache}
}
