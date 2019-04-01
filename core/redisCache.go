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

// Get(key string) (interface{}, error)
// 	Set(key string, value interface{}) (bool, error)
// 	IsSet(key string) (bool, error)
// 	Remove(key string) (bool, error)
// 	RemoveByPattern(pattern string) (bool, error)

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
func (rc RedisCache) IsSet(key string) (bool, error) {
	isSet, err := rc.client.Exists(key).Result()
	return isSet > 0, err
}

// Remove deletes a value by its key
func (rc RedisCache) Remove(key string) (bool, error) {
	isRemoved, err := rc.client.Del(key).Result()
	return isRemoved > 0, err
}

// RemoveByPattern removes all values matching pattern
func (rc RedisCache) RemoveByPattern(pattern string) (bool, error) {
	keys, err := rc.client.Keys(pattern).Result()
	if err != nil {
		panic("Could not identify the key set")
	}
	isRemoved, rmErr := rc.client.Del(strings.Join(keys, ",")).Result()
	return isRemoved > 0, rmErr
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

// NewRedisCache manager
func NewRedisCache(config *Config) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Host,
		Password: config.Redis.Password, // no password set
		DB:       config.Redis.Db,       // use default DB
	})
	return &RedisCache{client: client}
}
