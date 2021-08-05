package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/zacscoding/echo-gorm-realworld-app/config"
	"github.com/zacscoding/echo-gorm-realworld-app/logging"
	"github.com/zacscoding/echo-gorm-realworld-app/server"
	"github.com/zacscoding/echo-gorm-realworld-app/serverenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	// setup server.
	srv, err := server.New(serverEnv, cfg)
	if err != nil {
		logging.DefaultLogger().Fatalw("failed to initialize server", "err", err)
	}

	// start server.
	appsrv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.ServerConfig.Port),
		Handler:      srv,
		ReadTimeout:  parseDuration(cfg.ServerConfig.ReadTimeout, time.Second*5),
		WriteTimeout: parseDuration(cfg.ServerConfig.WriteTimeout, time.Second*10),
	}

	go func() {
		if err := appsrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.DefaultLogger().Fatal(err)
		}
	}()

	// wait for termination signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	<-quit

	logging.DefaultLogger().Info("Shutting down app server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logging.DefaultLogger().Fatal(err)
	}
	logging.DefaultLogger().Info("Terminate application")
}

func parseDuration(v string, defaultDuration time.Duration) time.Duration {
	d, err := time.ParseDuration(v)
	if err != nil {
		return defaultDuration
	}
	return d
}
