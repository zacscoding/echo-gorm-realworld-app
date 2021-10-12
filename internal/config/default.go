package config

var defaultConfig = map[string]interface{}{
	"server.port":         8080,
	"server.timeout":      "5s",
	"server.readTimeout":  "5s",
	"server.writeTimeout": "10s",
	"server.docs.enabled": true,
	"server.docs.path":    "/config/doc.html",

	"jwt.secret":         "secret-key",
	"jwt.sessionTimeout": "240h",

	"db.dataSourceName":   "root:password@tcp(127.0.0.1:3306)/local_db?charset=utf8&parseTime=True&multiStatements=true",
	"db.migrate.enable":   false,
	"db.migrate.dir":      "",
	"db.pool.maxOpen":     50,
	"db.pool.maxIdle":     5,
	"db.pool.maxLifetime": "86400s",
}