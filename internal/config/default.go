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
}
