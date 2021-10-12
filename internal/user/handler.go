package user

import (
	"github.com/labstack/echo/v4"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/config"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/serverenv"
	userDB "github.com/zacscoding/echo-gorm-realworld-app/internal/user/database"
	"time"
)

type Handler struct {
	cfg         *config.Config
	userDB      userDB.UserDB
	jwtSecret   []byte
	jwtDuration time.Duration
}

// NewHandler returns a new Handle from given serverenv.ServerEnv and config.Config.
func NewHandler(env *serverenv.ServerEnv, conf *config.Config) (*Handler, error) {
	return &Handler{
		cfg:         conf,
		userDB:      userDB.NewUserDB(conf, env.GetDB()),
		jwtSecret:   []byte(conf.JWTConfig.Secret),
		jwtDuration: conf.JWTConfig.SessionTimeout,
	}, nil
}

// Route configures route given "/api" echo.Group to "/api/users/**, /api/profile/**" paths.
func (h *Handler) Route(e *echo.Group, authMiddleware echo.MiddlewareFunc) {
	// anonymous
	anonymousUserGroup := e.Group("/users")
	anonymousUserGroup.POST("/login", h.handleSignIn)
	anonymousUserGroup.POST("", h.handleSignUp)

	// auth required
	userGroup := e.Group("/user")
	userGroup.Use(authMiddleware)
	userGroup.GET("", h.handleCurrentUser)
	userGroup.PUT("", h.handleUpdateUser)

	profileGroup := e.Group("/profiles")
	profileGroup.Use(authMiddleware)
	profileGroup.GET("/:username", h.handleGetProfile)
	profileGroup.POST("/:username/follow", h.handleFollow)
	profileGroup.DELETE("/:username/follow", h.handleUnfollow)
}
