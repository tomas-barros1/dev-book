package cache

import "time"

type Cache interface {
	Get(key string, dest interface{}) error
	Set(key string, value interface{}, ttl time.Duration) error
	Del(key string) error
	Close() error
}
