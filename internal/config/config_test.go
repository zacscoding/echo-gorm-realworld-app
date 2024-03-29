package config

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	cfg, err := Load("")

	assert.NoError(t, err)
	// logging configs
	equal(t, -1, defaultConfig["logging.level"].(int), cfg.LoggingConfig.Level)
	equal(t, "console", defaultConfig["logging.encoding"].(string), cfg.LoggingConfig.Encoding)
	equal(t, true, defaultConfig["logging.development"].(bool), cfg.LoggingConfig.Development)

	// server configs
	equal(t, 8080, defaultConfig["server.port"].(int), cfg.ServerConfig.Port)
	equal(t, 5*time.Second, defaultConfig["server.timeout"].(time.Duration), cfg.ServerConfig.Timeout)
	equal(t, 5*time.Second, defaultConfig["server.readTimeout"].(time.Duration), cfg.ServerConfig.ReadTimeout)
	equal(t, 10*time.Second, defaultConfig["server.writeTimeout"].(time.Duration), cfg.ServerConfig.WriteTimeout)
	equal(t, true, defaultConfig["server.docs.enabled"].(bool), cfg.ServerConfig.Docs.Enabled)
	equal(t, "/config/doc.html", defaultConfig["server.docs.path"].(string), cfg.ServerConfig.Docs.Path)
	// jwt configs
	equal(t, "secret-key", defaultConfig["jwt.secret"].(string), cfg.JWTConfig.Secret)
	equal(t, 240*time.Hour, defaultConfig["jwt.sessionTimeout"].(time.Duration), cfg.JWTConfig.SessionTimeout)
	// db configs
	equal(t, "root:password@tcp(127.0.0.1:3306)/local_db?charset=utf8&parseTime=True&multiStatements=true",
		defaultConfig["db.dataSourceName"].(string), cfg.DBConfig.DataSourceName)
	equal(t, false, defaultConfig["db.migrate.enable"].(bool), cfg.DBConfig.Migrate.Enable)
	equal(t, "", defaultConfig["db.migrate.dir"].(string), cfg.DBConfig.Migrate.Dir)
	equal(t, 50, defaultConfig["db.pool.maxOpen"].(int), cfg.DBConfig.Pool.MaxOpen)
	equal(t, 5, defaultConfig["db.pool.maxIdle"].(int), cfg.DBConfig.Pool.MaxIdle)
	equal(t, 86400*time.Second, defaultConfig["db.pool.maxLifetime"].(time.Duration), cfg.DBConfig.Pool.MaxLifetime)
	// redis configs
	equal(t, false, defaultConfig["cache.enabled"].(bool), cfg.CacheConfig.Enabled)
	equal(t, "rewalworld-", defaultConfig["cache.prefix"].(string), cfg.CacheConfig.Prefix)
	equal(t, "redis", defaultConfig["cache.type"].(string), cfg.CacheConfig.Type)
	equal(t, 60*time.Second, defaultConfig["cache.ttl"].(time.Duration), cfg.CacheConfig.TTL)
	equal(t, false, defaultConfig["cache.redis.cluster"].(bool), cfg.CacheConfig.RedisConfig.Cluster)
	equal(t, []string{"localhost:6379"}, defaultConfig["cache.redis.endpoints"].([]string), cfg.CacheConfig.RedisConfig.Endpoints)
	equal(t, 3*time.Second, defaultConfig["cache.redis.readTimeout"].(time.Duration), cfg.CacheConfig.RedisConfig.ReadTimeout)
	equal(t, 3*time.Second, defaultConfig["cache.redis.writeTimeout"].(time.Duration), cfg.CacheConfig.RedisConfig.WriteTimeout)
	equal(t, 5*time.Second, defaultConfig["cache.redis.dialTimeout"].(time.Duration), cfg.CacheConfig.RedisConfig.DialTimeout)
	equal(t, 10, defaultConfig["cache.redis.poolSize"].(int), cfg.CacheConfig.RedisConfig.PoolSize)
	equal(t, 1*time.Minute, defaultConfig["cache.redis.poolTimeout"].(time.Duration), cfg.CacheConfig.RedisConfig.PoolTimeout)
	equal(t, 0, defaultConfig["cache.redis.maxConnAge"].(time.Duration), cfg.CacheConfig.RedisConfig.MaxConnAge)
	equal(t, 5*time.Minute, defaultConfig["cache.redis.idleTimeout"].(time.Duration), cfg.CacheConfig.RedisConfig.IdleTimeout)
}

func equal(t *testing.T, expected, defaultValue, actualValue interface{}) {
	assert.EqualValues(t, expected, defaultValue)
	assert.EqualValues(t, expected, actualValue)
}

func TestLoadWithEnv(t *testing.T) {
	// given
	err := os.Setenv(fmt.Sprintf("%sSERVER_PORT", EnvPrefix), "4000")
	assert.NoError(t, err)

	// when
	cfg, err := Load("")

	// then
	assert.NoError(t, err)
	assert.Equal(t, 4000, cfg.ServerConfig.Port)
}

func TestLoadWithConfigFile(t *testing.T) {
	// given
	err := os.Setenv(fmt.Sprintf("%sSERVER_PORT", EnvPrefix), "4000")
	assert.NoError(t, err)

	// when
	cfg, err := Load("test-config.yaml")

	// then
	assert.NoError(t, err)
	assert.Equal(t, 5000, cfg.ServerConfig.Port)
}

func TestMarshalJSON(t *testing.T) {
	cfg, err := Load("")
	assert.NoError(t, err)
	data, err := json.Marshal(cfg)
	assert.NoError(t, err)

	var m map[string]interface{}
	assert.NoError(t, json.Unmarshal(data, &m))

	assert.Equal(t, "root:****@tcp(127.0.0.1:3306)/local_db?charset=utf8&parseTime=True&multiStatements=true", m["db.dataSourceName"])
	assert.Equal(t, "****", m["jwt.secret"])
}
