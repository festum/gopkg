package cache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type Cache struct {
	*gocache.Cache
}

// NewLimitCache create cache with limit
func NewLimitCache(expiration, interval time.Duration) *Cache {
	return &Cache{
		gocache.New(expiration, interval),
	}
}

// NewCache create cache without limit
func NewCache() *Cache {
	return &Cache{
		gocache.New(gocache.NoExpiration, gocache.NoExpiration),
	}
}

// Save value by key permanently
func (c Cache) Save(key string, msg interface{}) {
	c.Set(key, msg, gocache.DefaultExpiration)
}
