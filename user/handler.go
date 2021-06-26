package user

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/zacscoding/echo-gorm-realworld-app/api/types"
	"github.com/zacscoding/echo-gorm-realworld-app/config"
	"github.com/zacscoding/echo-gorm-realworld-app/database"
	"github.com/zacscoding/echo-gorm-realworld-app/logging"
	"github.com/zacscoding/echo-gorm-realworld-app/serverenv"
	userDB "github.com/zacscoding/echo-gorm-realworld-app/user/database"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/user/model"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/authutils"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/hashutils"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/httputils"
	"go.uber.org/zap"
	"net/http"
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

	// TODO: add profiles
	profileGroup := e.Group("/profile")
	profileGroup.Use(authMiddleware)
	profileGroup.GET("/:username", h.handleGetProfile)
	profileGroup.POST("/:username/follow", h.handleFollow)
	profileGroup.DELETE("/:username/follow", h.handleUnfollow)
}

// handleSignUp handles "POST /api/users" to register a new user.
func (h *Handler) handleSignUp(c echo.Context) error {
	var (
		logger = logging.FromContext(c.Request().Context())
		ctx    = c.Request().Context()
		req    = &SignUpRequest{}
		user   userModel.User
	)
	// Bind request
	if err := req.Bind(c, &user); err != nil {
		logger.Errorw("UserHandler_handlePostUser failed to bind register request", "err", err)
		return httputils.WrapBindError(err)
	}
	// Save given user
	if err := h.userDB.Save(ctx, &user); err != nil {
		return err
	}
	return h.responseUser(c, &user)
}

// handleSignIn handles "POST /api/users/login" to sign in an user.
func (h *Handler) handleSignIn(c echo.Context) error {
	var (
		logger = logging.FromContext(c.Request().Context())
		ctx    = c.Request().Context()
		req    = &SignInRequest{}
	)
	// Bind request
	if err := req.Bind(c); err != nil {
		logger.Errorw("UserHandler_handleSignIn failed to bind register request", "err", err)
		return httputils.WrapBindError(err)
	}
	// Find an user from given email
	user, err := h.userDB.FindByEmail(ctx, req.User.Email)
	if err != nil {
		logger.Errorw("UserHandler_handleSignIn failed to find an user", "err", err)
		if err == database.ErrRecordNotFound {
			return httputils.NewNotFoundError(fmt.Sprintf("user(%s) not found", req.User.Email))
		}
		return httputils.NewInternalServerError(err)
	}
	// Check password
	if err := hashutils.MatchesPassword(user.Password, req.User.Password); err != nil {
		logger.Errorw("UserHandler_handleSignIn failed to sign in with wrong password", "err", err)
		return httputils.NewStatusUnprocessableEntity("password mismatch")
	}
	return h.responseUser(c, user)
}

// handleCurrentUser handles "GET /api/user" to get current user.
func (h *Handler) handleCurrentUser(c echo.Context) error {
	var (
		logger = logging.FromContext(c.Request().Context())
		ctx    = c.Request().Context()
	)
	// Find current user
	user, err := h.userDB.FindByID(ctx, authutils.CurrentUser(c))
	if err != nil {
		logger.Errorw("UserHandler_handleUpdateUser failed to find an user", "err", err)
		return httputils.NewInternalServerError(err)
	}
	return h.responseUser(c, user)
}

// handleSignIn handles "PUT /api/user" to update current user.
func (h *Handler) handleUpdateUser(c echo.Context) error {
	var (
		logger = logging.FromContext(c.Request().Context())
		ctx    = c.Request().Context()
		req    = &UpdateUserRequest{}
	)
	// Find current user
	user, err := h.userDB.FindByID(ctx, authutils.CurrentUser(c))
	if err != nil {
		logger.Errorw("UserHandler_handleUpdateUser failed to find an user", "err", err)
		return httputils.NewInternalServerError(nil)
	}
	// Bind request
	if err := req.Bind(c, user); err != nil {
		logger.Errorw("UserHandler_handleUpdateUser failed to bind request", "err", err)
		return httputils.WrapBindError(err)
	}
	// Update user
	if err := h.userDB.Update(ctx, user); err != nil {
		logger.Errorw("UserHandler_handleUpdateUser failed to update an user", "err", err)
		return httputils.NewInternalServerError(err)
	}
	return h.responseUser(c, user)
}

func (h *Handler) handleGetProfile(c echo.Context) error {
	var (
		logger   = logging.FromContext(c.Request().Context())
		ctx      = c.Request().Context()
		username = c.Param("username")
	)

	user, err := h.getUserByUsername(ctx, logger, username)
	if err != nil {
		return err
	}

	currentUserID := authutils.CurrentUser(c)
	if currentUserID != 0 && currentUserID != user.ID {
		isFollow, err := h.userDB.IsFollow(ctx, currentUserID, user.ID)
		if err != nil {
			logger.Errorw("UserHandler_handleGetProfile failed to find the following relation",
				"userID", currentUserID, "followID", user.ID, "err", err)
			return httputils.NewInternalServerError(err)
		}
		user.Following = isFollow
	}
	return c.JSON(http.StatusOK, types.ToUserProfile(user))
}

func (h *Handler) handleFollow(c echo.Context) error {
	var (
		logger   = logging.FromContext(c.Request().Context())
		ctx      = c.Request().Context()
		username = c.Param("username")
	)
	user, err := h.getUserByUsername(ctx, logger, username)
	if err != nil {
		return err
	}
	currentUserID := authutils.CurrentUser(c)
	if err := h.userDB.Follow(ctx, currentUserID, user.ID); err != nil {
		logger.Errorw("UserHandler_handleFollow", "failed to update follow relation", "userID", currentUserID,
			"followID", user.ID, "err", err)
		return httputils.NewInternalServerError(err)
	}
	user.Following = true
	return c.JSON(http.StatusOK, types.ToUserProfile(user))
}

func (h *Handler) handleUnfollow(c echo.Context) error {
	var (
		logger   = logging.FromContext(c.Request().Context())
		ctx      = c.Request().Context()
		username = c.Param("username")
	)
	user, err := h.getUserByUsername(ctx, logger, username)
	if err != nil {
		return err
	}
	currentUserID := authutils.CurrentUser(c)
	if err := h.userDB.UnFollow(ctx, currentUserID, user.ID); err != nil {
		logger.Errorw("UserHandler_handleFollow failed to update unfollow relation", "userID", currentUserID, "followID", user.ID, "err", err)
		if err == database.ErrRecordNotFound {
			return httputils.NewNotFoundError(fmt.Sprintf("user already unfollowing user(%s)", username))
		}
		return httputils.NewInternalServerError(err)
	}
	user.Following = false
	return c.JSON(http.StatusOK, types.ToUserProfile(user))
}

func (h *Handler) getUserByUsername(ctx context.Context, logger *zap.SugaredLogger, username string) (*userModel.User, error) {
	user, err := h.userDB.FindByName(ctx, username)
	if err != nil {
		logger.Errorw("UserHandler_getUserByUsername failed to find an user", "err", err)
		if err == database.ErrRecordNotFound {
			return nil, httputils.NewNotFoundError(fmt.Sprintf("user(%s) not found", username))
		}
		return nil, httputils.NewInternalServerError(err)
	}
	return user, nil
}

func (h *Handler) responseUser(c echo.Context, user *userModel.User) error {
	logger := logging.FromContext(c.Request().Context())
	// Make JWT token
	token, err := h.makeJWTToken(user)
	if err != nil {
		logger.Errorw("UserHandler_responseUser failed to generate JWT token", "err", err)
		return httputils.NewInternalServerError(err)
	}
	return c.JSON(http.StatusOK, types.ToUserResponse(user, token))
}

func (h *Handler) makeJWTToken(u *userModel.User) (string, error) {
	return authutils.MakeJWTToken(u.ID, h.jwtSecret, h.jwtDuration)
}
