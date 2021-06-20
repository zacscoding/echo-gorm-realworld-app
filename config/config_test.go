package config

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	cfg, err := Load("")

	assert.NoError(t, err)
	// server configs
	assert.Equal(t, defaultConfig["server.port"].(int), cfg.ServerConfig.Port)
	assert.Equal(t, defaultConfig["server.timeout"].(string), cfg.ServerConfig.Timeout)
	assert.Equal(t, defaultConfig["server.readTimeout"].(string), cfg.ServerConfig.ReadTimeout)
	assert.Equal(t, defaultConfig["server.writeTimeout"].(string), cfg.ServerConfig.WriteTimeout)
	// jwt configs
	assert.Equal(t, defaultConfig["jwt.secret"].(string), cfg.JWTConfig.Secret)
	assert.Equal(t, defaultConfig["jwt.sessionTimeout"].(string), cfg.JWTConfig.SessionTimeout)
	// db configs
	assert.Equal(t, defaultConfig["db.dataSourceName"].(string), cfg.DBConfig.DataSourceName)
	assert.Equal(t, defaultConfig["db.migrate.enable"].(bool), cfg.DBConfig.Migrate.Enable)
	assert.Equal(t, defaultConfig["db.migrate.dir"].(string), cfg.DBConfig.Migrate.Dir)
	assert.Equal(t, defaultConfig["db.pool.maxOpen"].(int), cfg.DBConfig.Pool.MaxOpen)
	assert.Equal(t, defaultConfig["db.pool.maxIdle"].(int), cfg.DBConfig.Pool.MaxIdle)
	assert.Equal(t, defaultConfig["db.pool.maxLifetime"].(int), cfg.DBConfig.Pool.MaxLifetime)
}

func TestLoadWithEnv(t *testing.T) {
	// given
	err := os.Setenv("REALWORLD_APP_SERVER_PORT", "4000")
	assert.NoError(t, err)

	// when
	cfg, err := Load("")

	// then
	assert.NoError(t, err)
	assert.Equal(t, 4000, cfg.ServerConfig.Port)
}

func TestLoadWithConfigFile(t *testing.T) {
	// given
	err := os.Setenv("REALWORLD_APP_SERVER_PORT", "4000")
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
