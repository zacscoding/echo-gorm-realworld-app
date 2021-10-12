package cache

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"testing"
)

type CloseFunc func() error

// NewTestCache starts a redis server based on inmemory(miniredis) and
// returns redis.UniversalClient,
func NewTestCache(tb testing.TB) (redis.UniversalClient, *cache.Cache, CloseFunc) {
	s := miniredis.RunT(tb)
	cli := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{s.Addr()},
	})
	c := cache.New(&cache.Options{Redis: cli})

	return cli, c, func() error {
		s.Close()
		return nil
	}
}
