package config

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	cfg, err := Load("")

	assert.NoError(t, err)
	// server configs
	equal(t, 8080, defaultConfig["server.port"].(int), cfg.ServerConfig.Port)
	equal(t, "5s", defaultConfig["server.timeout"].(string), cfg.ServerConfig.Timeout)
	equal(t, "5s", defaultConfig["server.readTimeout"].(string), cfg.ServerConfig.ReadTimeout)
	equal(t, "10s", defaultConfig["server.writeTimeout"].(string), cfg.ServerConfig.WriteTimeout)
	equal(t, true, defaultConfig["server.docs.enabled"].(bool), cfg.ServerConfig.Docs.Enabled)
	equal(t, "/config/doc.html", defaultConfig["server.docs.path"].(string), cfg.ServerConfig.Docs.Path)
	// jwt configs
	equal(t, "secret-key", defaultConfig["jwt.secret"].(string), cfg.JWTConfig.Secret)
	equal(t, "240h", defaultConfig["jwt.sessionTimeout"].(string), cfg.JWTConfig.SessionTimeout)
	// db configs
	equal(t, "root:password@tcp(127.0.0.1:3306)/local_db?charset=utf8&parseTime=True&multiStatements=true",
		defaultConfig["db.dataSourceName"].(string), cfg.DBConfig.DataSourceName)
	equal(t, false, defaultConfig["db.migrate.enable"].(bool), cfg.DBConfig.Migrate.Enable)
	equal(t, "", defaultConfig["db.migrate.dir"].(string), cfg.DBConfig.Migrate.Dir)
	equal(t, 50, defaultConfig["db.pool.maxOpen"].(int), cfg.DBConfig.Pool.MaxOpen)
	equal(t, 5, defaultConfig["db.pool.maxIdle"].(int), cfg.DBConfig.Pool.MaxIdle)
	equal(t, "86400s", defaultConfig["db.pool.maxLifetime"].(string), cfg.DBConfig.Pool.MaxLifetime)
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
