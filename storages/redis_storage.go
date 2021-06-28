package storages

import (
	"context"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

type RedisStorage struct {
	redis   *redis.Client
	prefix  string
	context context.Context
}

func NewRedisStorage(address string, prefix string) *RedisStorage {
	storage := &RedisStorage{
		context: context.Background(),
		redis: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: "",
			DB:       0,
		}),
		prefix: prefix,
	}

	return storage
}

func (s *RedisStorage) Get(key string) (string, error) {
	fullKey := s.buildKey(key)
	return s.redis.Get(s.context, fullKey).Result()
}

func (s *RedisStorage) Has(key string) bool {
	if _, err := s.redis.Exists(s.context, s.buildKey(key)).Result(); err != nil {
		return true
	}

	return false
}

func (s *RedisStorage) TTL(key string) (time.Duration, error) {
	fullKey := s.buildKey(key)
	return s.redis.TTL(s.context, fullKey).Result()
}

func (s *RedisStorage) Set(key string, value interface{}, expiration time.Duration) (string, error) {
	fullKey := s.buildKey(key)
	return s.redis.Set(s.context, fullKey, value, expiration).Result()
}

func (s *RedisStorage) Keys(pattern string) ([]string, error) {
	fullPattern := s.buildKey(pattern)
	keys, error := s.redis.Keys(s.context, fullPattern).Result()
	result := make([]string, len(keys))
	for i, key := range keys {
		result[i] = strings.Replace(key, s.prefix, "", len(s.prefix))
	}

	return result, error
}

func (s *RedisStorage) Delete(keys ...string) (int64, error) {
	for i, key := range keys {
		keys[i] = s.buildKey(key)
	}

	return s.redis.Del(s.context, keys...).Result()
}

func (s *RedisStorage) buildKey(key string) string {
	return s.prefix + key
}
