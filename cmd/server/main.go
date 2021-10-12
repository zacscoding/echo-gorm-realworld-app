package main

import (
	"encoding/json"
	"flag"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/config"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/logging"
	"go.uber.org/zap/zapcore"
	"log"
)

func main() {
	// setup configs
	configFile := flag.String("config", "fixtures/config/config.yaml", "indicates a config file")
	flag.Parse()

	conf, err := config.Load(*configFile)
	if err != nil {
		log.Fatal("failed to initialize configs. err:", err)
	}
	logging.SetConfig(&logging.Config{
		Level:       zapcore.Level(conf.LoggingConfig.Level),
		Encoding:    conf.LoggingConfig.Encoding,
		Development: conf.LoggingConfig.Development,
	})
	data, _ := json.MarshalIndent(conf, "", "    ")
	logging.DefaultLogger().Infof("Starting a new application server. configs\n%s", string(data))

	startAppServer(conf)
}
