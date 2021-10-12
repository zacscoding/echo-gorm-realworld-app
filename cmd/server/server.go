package main

import (
	"context"
	"fmt"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/config"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/server"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/serverenv"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/logging"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func startAppServer(conf *config.Config) {
	// setup server environments.
	serverEnv, err := serverenv.SetupWith(conf)
	if err != nil {
		logging.DefaultLogger().Fatalw("failed to setup server environments", "err", err)
	}

	// setup server.
	srv, err := server.New(serverEnv, conf)
	if err != nil {
		logging.DefaultLogger().Fatalw("failed to initialize server", "err", err)
	}

	// start server.
	appsrv := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.ServerConfig.Port),
		Handler:      srv,
		ReadTimeout:  parseDuration(conf.ServerConfig.ReadTimeout, time.Second*5),
		WriteTimeout: parseDuration(conf.ServerConfig.WriteTimeout, time.Second*10),
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
