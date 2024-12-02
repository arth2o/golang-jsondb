package engine

import "time"

type Engine interface {
    Set(key string, value []byte) error
    Get(key string) ([]byte, error)
    Delete(key string) error
    GetByPattern(pattern string) ([]map[string]interface{}, error)
    SetWithTTL(key string, value []byte, ttl time.Duration) error
    TTL(key string) (time.Duration, error)
}