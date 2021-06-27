package user

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tidwall/gjson"
	"github.com/zacscoding/echo-gorm-realworld-app/config"
	"github.com/zacscoding/echo-gorm-realworld-app/database"
	"github.com/zacscoding/echo-gorm-realworld-app/logging"
	"github.com/zacscoding/echo-gorm-realworld-app/serverenv"
	userMocks "github.com/zacscoding/echo-gorm-realworld-app/user/database/mocks"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/user/model"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/hashutils"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/httputils"
	"go.uber.org/zap/zapcore"
	"testing"
)

var (
	cfg, _       = config.Load("")
	defaultUser1 = newUser(1, "user1", false)
	disabledUser = newUser(2, "disabledUser", true)
)

func newUser(id uint, username string, disable bool) *userModel.User {
	p, _ := hashutils.EncodePassword(username)
	return &userModel.User{
		ID:       id,
		Email:    fmt.Sprintf("%s@gmail.com", username),
		Name:     username,
		Password: p,
		Bio:      username + " bio",
		Image:    username + " image",
		Disabled: disable,
	}
}

func fixtures(_ *testing.T) (e *echo.Echo, h *Handler, u *userMocks.UserDB) {
	logging.SetConfig(&logging.Config{
		Encoding: "console",
		Level:    zapcore.FatalLevel,
	})
	u = &userMocks.UserDB{}

	// setup Save
	u.On("Save", mock.Anything, mock.MatchedBy(func(u *userModel.User) bool {
		return defaultUser1.Email == u.Email
	})).Return(nil)
	u.On("Save", mock.Anything, mock.MatchedBy(func(u *userModel.User) bool {
		return disabledUser.Email == u.Email
	})).Return(database.ErrKeyConflict)
	u.On("Save", mock.Anything, mock.Anything).Return(errors.New("force error"))

	// setup FindByEmail
	u.On("FindByEmail", mock.Anything, mock.MatchedBy(func(email string) bool {
		return defaultUser1.Email == email
	})).Return(defaultUser1, nil)
	u.On("FindByEmail", mock.Anything, mock.MatchedBy(func(email string) bool {
		return disabledUser.Email == email
	})).Return(nil, database.ErrRecordNotFound)
	u.On("FindByEmail", mock.Anything, mock.Anything).Return(nil, errors.New("force error"))

	// setup FindByID
	u.On("FindByID", mock.Anything, mock.MatchedBy(func(id uint) bool {
		return defaultUser1.ID == id
	})).Return(defaultUser1, nil)
	u.On("FindByID", mock.Anything, mock.MatchedBy(func(id uint) bool {
		return disabledUser.ID == id
	})).Return(nil, database.ErrRecordNotFound)
	u.On("FindByID", mock.Anything, mock.Anything).Return(nil, errors.New("force error"))

	env := serverenv.NewServerEnv(serverenv.WithUserDB(u))
	h, _ = NewHandler(env, cfg)

	e = echo.New()
	e.Validator = httputils.NewValidator()

	return e, h, u
}

func assertUserResponse(t *testing.T, res string, expected *userModel.User, isSignUp bool) {
	u := gjson.Get(res, "user")
	assert.True(t, u.Exists())
	assert.Equal(t, expected.Email, u.Get("email").String())
	assert.NotEmpty(t, u.Get("token").String())
	assert.Equal(t, expected.Name, u.Get("username").String())
	if isSignUp {
		assert.Empty(t, u.Get("bio").String())
		assert.Empty(t, u.Get("image").String())
	} else {
		assert.Equal(t, expected.Bio, u.Get("bio").String())
		assert.Equal(t, expected.Image, u.Get("image").String())
	}
}

func assertErrorResponse(t *testing.T, statusCode int, msg string, actual error) {
	httpErr, ok := actual.(*echo.HTTPError)
	assert.True(t, ok)
	assert.Equal(t, statusCode, httpErr.Code)
	e := httpErr.Message.(*httputils.Error)
	assert.Contains(t, e.Errors["body"].(string), msg)
}
