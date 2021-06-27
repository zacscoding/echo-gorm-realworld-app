package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zacscoding/echo-gorm-realworld-app/config"
	"github.com/zacscoding/echo-gorm-realworld-app/logging"
	"github.com/zacscoding/echo-gorm-realworld-app/serverenv"
	"github.com/zacscoding/echo-gorm-realworld-app/user"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/authutils"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/httputils"
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

	// TODO: temporary server code.
	e := echo.New()
	e.Use(middleware.Recover())
	e.Validator = httputils.NewValidator()
	v1 := e.Group("/api")
	authMiddleware := authutils.NewJWTMiddleware(
		map[string]struct{}{
			"/api/profile/:username": {},
		},
		cfg.JWTConfig.Secret,
	)
	userHandler, err := user.NewHandler(serverEnv, cfg)
	if err != nil {
		logging.DefaultLogger().Fatalw("failed to initialize user handler", "err", err)
	}
	userHandler.Route(v1, authMiddleware)

	if e.Start(fmt.Sprintf(":%d", cfg.ServerConfig.Port)); err != nil {
		logging.DefaultLogger().Fatalw("shutting down server", "err", err)
	}
}
