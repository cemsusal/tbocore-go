package core

import (
	"fmt"
)

type theCache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) (bool, error)
	IsSet(key string) (bool, error)
	Remove(key string) (bool, error)
	RemoveByPattern(pattern string) (bool, error)
}

// CacheManager manages cache
type CacheManager struct {
	cache theCache
}

type acquire func() interface{}

// Get gets or sets a value in cache
func (c *CacheManager) Get(key string, expiresInSeconds int, acquire acquire) interface{} {
	isSet, err := c.cache.IsSet(key)
	if err != nil {
		fmt.Println("Cache is set error")
	}
	if isSet {
		cacheResult, cacheErr := c.cache.Get(key)
		if cacheErr != nil {
			fmt.Println("Cache is set error")
		}
		return cacheResult
	}
	result := acquire()

	c.cache.Set(key, result)
	return result
}
