package storages

import (
	"time"
)

type SessionStorage struct {
	redis *RedisStorage
}

func NewSessionStorage(redis *RedisStorage) *SessionStorage {
	return &SessionStorage{redis: redis}
}

func (s *SessionStorage) GetValue(key string) (string, error) {
	return s.redis.Get(key)
}

func (s *SessionStorage) GetTTL(key string) (time.Duration, error) {
	return s.redis.TTL(key)
}

func (s *SessionStorage) Save(key string, value interface{}, expiration time.Duration) (string, error) {
	return s.redis.Set(key, value, expiration)
}

func (s *SessionStorage) SearchKeys(pattern string) ([]string, error) {
	return s.redis.Keys(pattern)
}

func (s *SessionStorage) DeleteKeys(keys ...string) (int64, error) {
	return s.redis.Delete(keys...)
}
