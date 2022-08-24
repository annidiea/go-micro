package cache

import (
	"context"
	"github.com/eko/gocache/v2/store"
	"time"
)

var CacheManager *Cache

type CacheInterface interface {
	Get(ctx context.Context, key interface{}) (interface{}, error)
	Set(ctx context.Context, key, object interface{}, options *store.Options) error
	Delete(ctx context.Context, key interface{}) error
	Clear(ctx context.Context) error
}

// 这是为自己项目的缓存而设计的；
type Cache struct {
	ctx   context.Context
	cache CacheInterface
}

// 属于自己当前的项目需要的options
type Options struct {
	Cost       int64
	Expiration time.Duration
	Tags       []string
}

func NewCache(cache CacheInterface) *Cache {
	ctx := context.Background()
	return &Cache{
		cache: cache,
		ctx:   ctx,
	}
}

func (c *Cache) Set(key, value interface{}, options *Options) error {
	option := store.Options(*options)
	return c.cache.Set(c.ctx, marshal(key), value, &option)
}

func (c *Cache) Get(key interface{}) (interface{}, error) {
	return c.cache.Get(c.ctx, key)
}

func (c *Cache) Delete(key interface{}) error {
	return c.cache.Delete(c.ctx, key)
}

func (c *Cache) Clear() error {
	return c.cache.Clear(c.ctx)
}

func marshal(v interface{}) (msg []byte) {
	switch data := v.(type) {
	case string:
		return []byte(data)
	case []byte:
		return data
	}
	return nil
}
