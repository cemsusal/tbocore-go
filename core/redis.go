package core

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"strings"
	"time"
)

// RedisCache implements redis caching
type RedisCache struct {
	client *redis.Client
}

// Set a value by a given key
func (rc RedisCache) Set(key string, value interface{}, expiresIn time.Duration) (bool, error) {
	serializedValue, _ := json.Marshal(value)
	err := rc.client.Set(key, string(serializedValue), expiresIn).Err()
	return true, err
}

// Get a value by a given key
func (rc RedisCache) Get(key string) (interface{}, error) {
	var deserializedValue interface{}
	serializedValue, err := rc.client.Get(key).Result()
	json.Unmarshal([]byte(serializedValue), &deserializedValue)
	return deserializedValue, err
}

// IsSet check a keys is already set
func (rc RedisCache) IsSet(key string) (int64, error) {
	isSet, err := rc.client.Exists(key).Result()
	return isSet, err
}

// Remove deletes a value by its key
func (rc RedisCache) Remove(key string) (int64, error) {
	isSet, err := rc.client.Del(key).Result()
	return isSet, err
}

// RemoveByPattern removes all values matching pattern
func (rc RedisCache) RemoveByPattern(pattern string) (int64, error) {
	keys, err := rc.client.Keys(pattern).Result()
	if err != nil {
		panic("Could not identify the key set")
	}
	isRemoved, rmErr := rc.client.Del(strings.Join(keys, ",")).Result()
	return isRemoved, rmErr
}

// Init a redis client instance
func (rc RedisCache) Init(addr string, poolSize int, maxRetries int, password string, db int) *RedisCache {
	rc.client = setClient(addr, poolSize, maxRetries, password, db)
	return &rc
}

func setClient(addr string, poolSize int, maxRetries int, password string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:       addr,
		PoolSize:   poolSize,
		MaxRetries: maxRetries,
		Password:   password,
		DB:         db,
	})
}
