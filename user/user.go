package user

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/zacscoding/echo-gorm-realworld-app/api/types"
	"github.com/zacscoding/echo-gorm-realworld-app/database"
	"github.com/zacscoding/echo-gorm-realworld-app/logging"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/user/model"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/authutils"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/hashutils"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/httputils"
	"net/http"
)

// handleSignUp handles "POST /api/users" to register a new user.
func (h *Handler) handleSignUp(c echo.Context) error {
	var (
		ctx    = c.Request().Context()
		logger = logging.FromContext(ctx)
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
		if err == database.ErrKeyConflict {
			return httputils.NewStatusUnprocessableEntity(fmt.Sprintf("duplicate email: %s", req.User.Email))
		}
		return httputils.NewInternalServerError(err)
	}
	return h.responseUser(c, &user)
}

// handleSignIn handles "POST /api/users/login" to sign in an user.
func (h *Handler) handleSignIn(c echo.Context) error {
	var (
		ctx    = c.Request().Context()
		logger = logging.FromContext(ctx)
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
		ctx    = c.Request().Context()
		logger = logging.FromContext(ctx)
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
		ctx    = c.Request().Context()
		logger = logging.FromContext(ctx)
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
