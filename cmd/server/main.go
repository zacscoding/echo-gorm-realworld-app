package main

import (
	"encoding/json"
	"flag"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/config"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/logging"
	"log"
)

func main() {
	// setup configs
	configFile := flag.String("config", "config.yaml", "indicates a config file")
	flag.Parse()

	conf, err := config.Load(*configFile)
	if err != nil {
		log.Fatal("failed to initialize configs. err:", err)
	}
	data, _ := json.MarshalIndent(conf, "", "    ")
	logging.DefaultLogger().Infof("Starting a new application server. configs\n%s", string(data))

	startAppServer(conf)
}
