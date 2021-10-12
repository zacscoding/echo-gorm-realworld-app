package user

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/database"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/internal/user/model"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/api/types"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/utils/authutils"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/utils/httputils"
	"net/http"
)

// handleGetProfile handles "GET /api/user/profile/:username" to get given user's profile.
func (h *Handler) handleGetProfile(c echo.Context) error {
	var (
		ctx      = c.Request().Context()
		username = c.Param("username")
	)

	user, err := h.getUserByUsername(ctx, username)
	if err != nil {
		return err
	}

	currentUserID := authutils.CurrentUser(c)
	if currentUserID != 0 && currentUserID != user.ID {
		isFollow, err := h.userDB.IsFollow(ctx, currentUserID, user.ID)
		if err != nil {
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
		username = c.Param("username")
	)
	user, err := h.getUserByUsername(ctx, username)
	if err != nil {
		return err
	}
	currentUserID := authutils.CurrentUser(c)
	if err := h.userDB.Follow(ctx, currentUserID, user.ID); err != nil {
		if err == database.ErrKeyConflict {
			return httputils.NewStatusUnprocessableEntity(fmt.Sprintf("user already following %s", user.Name))
		}
		return httputils.NewInternalServerError(err)
	}
	user.Following = true
	return c.JSON(http.StatusOK, types.ToUserProfile(user))
}

// handleUnfollow handles "DELETE /api/user/profile/:username/follow" to delete the following relation.
func (h *Handler) handleUnfollow(c echo.Context) error {
	var (
		ctx      = c.Request().Context()
		username = c.Param("username")
	)
	user, err := h.getUserByUsername(ctx, username)
	if err != nil {
		return err
	}
	currentUserID := authutils.CurrentUser(c)
	if err := h.userDB.UnFollow(ctx, currentUserID, user.ID); err != nil {
		if err == database.ErrRecordNotFound {
			return httputils.NewStatusUnprocessableEntity(fmt.Sprintf("user already unfollowing user(%s)", username))
		}
		return httputils.NewInternalServerError(err)
	}
	user.Following = false
	return c.JSON(http.StatusOK, types.ToUserProfile(user))
}

func (h *Handler) getUserByUsername(ctx context.Context, username string) (*userModel.User, error) {
	user, err := h.userDB.FindByName(ctx, username)
	if err != nil {
		if err == database.ErrRecordNotFound {
			return nil, httputils.NewNotFoundError(fmt.Sprintf("user(%s) not found", username))
		}
		return nil, httputils.NewInternalServerError(err)
	}
	return user, nil
}
