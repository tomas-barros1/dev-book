package cache

import (
	"errors"
	"time"
)

type noopCache struct{}

func NewNoop() Cache {
	return &noopCache{}
}

var errNoopCacheMiss = errors.New("noop cache: miss")

func (n *noopCache) Get(key string, dest interface{}) error {
	return errNoopCacheMiss
}

func (n *noopCache) Set(key string, value interface{}, ttl time.Duration) error {
	return nil
}

func (n *noopCache) Del(key string) error {
	return nil
}

func (n *noopCache) Close() error {
	return nil
}
