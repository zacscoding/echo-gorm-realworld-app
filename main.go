package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/zacscoding/echo-gorm-realworld-app/config"
	"github.com/zacscoding/echo-gorm-realworld-app/logging"
	"github.com/zacscoding/echo-gorm-realworld-app/serverenv"
	"log"
)

func main() {
	// setup configs
	configFile := flag.String("config", "config.yaml", "indicates a config file")
	flag.Parse()

	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatal("failed to initialize configs. err:", err)
	}
	data, _ := json.MarshalIndent(cfg, "", "    ")
	logging.DefaultLogger().Infof("Starting a new application server. configs\n%s", string(data))

	// setup server environments.
	serverEnv, err := serverenv.SetupWith(cfg)
	if err != nil {
		logging.DefaultLogger().Fatalw("failed to setup server environments", "err", err)
	}

	// TODO : setup server
	fmt.Println(serverEnv)
}
