package user

import (
	"github.com/labstack/echo/v4"
	"github.com/zacscoding/echo-gorm-realworld-app/config"
	"github.com/zacscoding/echo-gorm-realworld-app/serverenv"
	userDB "github.com/zacscoding/echo-gorm-realworld-app/user/database"
	"time"
)

type Handler struct {
	cfg         *config.Config
	userDB      userDB.UserDB
	jwtSecret   []byte
	jwtDuration time.Duration
}

// NewHandler returns a new Handle from given serverenv.ServerEnv and config.Config.
func NewHandler(env *serverenv.ServerEnv, cfg *config.Config) (*Handler, error) {
	jwtDuration, err := time.ParseDuration(cfg.JWTConfig.SessionTimeout)
	if err != nil {
		return nil, err
	}
	return &Handler{
		cfg:         cfg,
		userDB:      env.GetUserDB(),
		jwtSecret:   []byte(cfg.JWTConfig.Secret),
		jwtDuration: jwtDuration,
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

	profileGroup := e.Group("/profile")
	profileGroup.Use(authMiddleware)
	profileGroup.GET("/:username", h.handleGetProfile)
	profileGroup.POST("/:username/follow", h.handleFollow)
	profileGroup.DELETE("/:username/unfollow", h.handleUnfollow)
}
