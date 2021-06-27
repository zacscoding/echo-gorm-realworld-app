package user

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/zacscoding/echo-gorm-realworld-app/api/types"
	"github.com/zacscoding/echo-gorm-realworld-app/database"
	"github.com/zacscoding/echo-gorm-realworld-app/logging"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/user/model"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/authutils"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/httputils"
	"go.uber.org/zap"
	"net/http"
)

// handleGetProfile handles "GET /api/user/profile/:username" to get given user's profile.
func (h *Handler) handleGetProfile(c echo.Context) error {
	var (
		ctx      = c.Request().Context()
		logger   = logging.FromContext(ctx)
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

// handleFollow handles "POST /api/user/profile/:username/follow" to update the following relation.
func (h *Handler) handleFollow(c echo.Context) error {
	var (
		ctx      = c.Request().Context()
		logger   = logging.FromContext(ctx)
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

// handleUnfollow handles "DELETE /api/user/profile/:username/follow" to delete the following relation.
func (h *Handler) handleUnfollow(c echo.Context) error {
	var (
		ctx      = c.Request().Context()
		logger   = logging.FromContext(ctx)
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
