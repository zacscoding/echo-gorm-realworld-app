package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/config"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/logging"
)

// NewCache creates a new redis.UniversalClient and cache.Cache.
func NewCache(conf *config.Config) (redis.UniversalClient, *cache.Cache, error) {
	if !conf.CacheConfig.Enabled {
		return nil, nil, fmt.Errorf("disabled cache in config")
	}

	var (
		redisConf = conf.CacheConfig.RedisConfig
		cli       redis.UniversalClient
	)
	if redisConf.Cluster {
		cli = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:         redisConf.Endpoints,
			ReadTimeout:   redisConf.ReadTimeout,
			WriteTimeout:  redisConf.WriteTimeout,
			DialTimeout:   redisConf.DialTimeout,
			PoolSize:      redisConf.PoolSize,
			PoolTimeout:   redisConf.PoolTimeout,
			MaxConnAge:    redisConf.MaxConnAge,
			IdleTimeout:   redisConf.IdleTimeout,
			ReadOnly:      true, // read on slave nodes.
			RouteRandomly: true, // read on masster or slave nodes.
		})
	} else {
		cli = redis.NewClient(&redis.Options{
			Addr:         redisConf.Endpoints[0],
			ReadTimeout:  redisConf.ReadTimeout,
			WriteTimeout: redisConf.WriteTimeout,
			DialTimeout:  redisConf.DialTimeout,
			PoolSize:     redisConf.PoolSize,
			PoolTimeout:  redisConf.PoolTimeout,
			MaxConnAge:   redisConf.MaxConnAge,
			IdleTimeout:  redisConf.IdleTimeout,
		})
	}
	// check ping.
	if err := cli.Ping(context.Background()).Err(); err != nil {
		logging.DefaultLogger().Infow("failed to ping redis", "err", err)
	} else {
		logging.DefaultLogger().Info("connected to redis")
	}

	return cli, cache.New(&cache.Options{
		Redis: cli,
	}), nil
}
