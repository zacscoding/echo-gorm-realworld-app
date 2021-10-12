package config

import "time"

var defaultConfig = map[string]interface{}{
	"logging.level":       -1,
	"logging.encoding":    "console",
	"logging.development": true,

	"server.port":         8080,
	"server.timeout":      5 * time.Second,
	"server.readTimeout":  5 * time.Second,
	"server.writeTimeout": 10 * time.Second,
	"server.docs.enabled": true,
	"server.docs.path":    "/config/doc.html",

	"jwt.secret":         "secret-key",
	"jwt.sessionTimeout": 240 * time.Hour,

	"db.dataSourceName":   "root:password@tcp(127.0.0.1:3306)/local_db?charset=utf8&parseTime=True&multiStatements=true",
	"db.migrate.enable":   false,
	"db.migrate.dir":      "",
	"db.pool.maxOpen":     50,
	"db.pool.maxIdle":     5,
	"db.pool.maxLifetime": 86400 * time.Second,

	"cache.enabled":            false,
	"cache.prefix":             "rewalworld-",
	"cache.type":               "redis",
	"cache.ttl":                60 * time.Second,
	"cache.redis.cluster":      false,
	"cache.redis.endpoints":    []string{"localhost:6379"},
	"cache.redis.readTimeout":  3 * time.Second,
	"cache.redis.writeTimeout": 3 * time.Second,
	"cache.redis.dialTimeout":  5 * time.Second,
	"cache.redis.poolSize":     10,
	"cache.redis.poolTimeout":  1 * time.Minute,
	"cache.redis.maxConnAge":   time.Duration(0),
	"cache.redis.idleTimeout":  5 * time.Minute,
}
