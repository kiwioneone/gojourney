package gojourney

import "github.com/coocood/freecache"

type Storage interface {
	Get(key []byte) (value []byte, err error)
	Set(key, value []byte) error
}

type LocalCache struct {
	cache *freecache.Cache
}

func NewLocalCache(size int) *LocalCache {
	return &LocalCache{cache: freecache.NewCache(size)}
}

func (c *LocalCache) Get(key []byte) (value []byte, err error) {
	return c.cache.Get(key)
}

func (c *LocalCache) Set(key, value []byte) error {
	return c.cache.Set(key, value, 0)
}
