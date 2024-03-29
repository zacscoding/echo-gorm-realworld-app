package server

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/article"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/config"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/serverenv"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/user"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/logging"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/utils/authutils"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/utils/httputils"
)

type Server struct {
	*echo.Echo
	articleHandler *article.Handler
	userHandler    *user.Handler
}

// New returns a new Server from given
func New(env *serverenv.ServerEnv, conf *config.Config) (*Server, error) {
	// Setup echo and middlewares.
	e := echo.New()
	e.Use(middleware.Recover(), middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid.NewString()
		},
		RequestIDHandler: func(c echo.Context, requestID string) {
			logger := logging.DefaultLogger().With("x-request-id", requestID)
			ctx := logging.WithLogger(c.Request().Context(), logger)
			c.SetRequest(c.Request().WithContext(ctx))
		},
	}))
	e.Validator = httputils.NewValidator()
	v1 := e.Group("/api")
	authMiddleware := authutils.NewJWTMiddleware(
		map[string]struct{}{
			"/api/profiles/:username":      {},
			"/api/articles":                {},
			"/api/articles/:slug":          {},
			"/api/articles/:slug/comments": {},
		},
		conf.JWTConfig.Secret,
	)

	// Setup handlers and route.
	userHandler, err := user.NewHandler(env, conf)
	if err != nil {
		return nil, errors.Wrap(err, "initialize user handlers")
	}
	userHandler.Route(v1, authMiddleware)

	articleHandler, err := article.NewHandler(env, conf)
	if err != nil {
		return nil, errors.Wrap(err, "initialize article handlers")
	}
	articleHandler.Route(v1, authMiddleware)

	// Serve api docs if enabled.
	if conf.ServerConfig.Docs.Enabled {
		e.Static("/docs", conf.ServerConfig.Docs.Path)
	}

	return &Server{
		Echo:           e,
		userHandler:    userHandler,
		articleHandler: articleHandler,
	}, nil
}
