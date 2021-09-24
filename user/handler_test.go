package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"
	"github.com/zacscoding/echo-gorm-realworld-app/config"
	"github.com/zacscoding/echo-gorm-realworld-app/logging"
	userMocks "github.com/zacscoding/echo-gorm-realworld-app/user/database/mocks"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/user/model"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/authutils"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/hashutils"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/httputils"
	"go.uber.org/zap/zapcore"
	"io"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	defaultUsers []*userModel.User
)

type TestSuite struct {
	suite.Suite
	e *echo.Echo
	h *Handler
	u *userMocks.UserDB
}

func TestRunSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) SetupSuite() {
	logging.SetConfig(&logging.Config{
		Encoding: "console",
		Level:    zapcore.FatalLevel,
	})
	cfg, err := config.Load("")
	s.NoError(err)

	e := echo.New()
	e.Validator = httputils.NewValidator()
	apiGroup := e.Group("/api")

	u := &userMocks.UserDB{}
	h := Handler{
		cfg:         cfg,
		userDB:      u,
		jwtSecret:   []byte(cfg.JWTConfig.Secret),
		jwtDuration: time.Hour,
	}
	h.Route(apiGroup, authutils.NewJWTMiddleware(map[string]struct{}{
		"/api/profiles/:username": {},
	}, cfg.JWTConfig.Secret))

	s.e = e
	s.h = &h
	s.u = u
}

func (s *TestSuite) SetupTest() {
	s.resetMocks()
	defaultUsers = []*userModel.User{}
	for i := 1; i <= 5; i++ {
		defaultUsers = append(defaultUsers, newUser(uint(i), fmt.Sprintf("user-%d", i), false))
	}
}

func (s *TestSuite) resetMocks() {
	u := &userMocks.UserDB{}
	s.h.userDB = u
	s.u = u
}

func assertUserResponse(t *testing.T, res string, expected *userModel.User, isSignUp bool) {
	a := assert.New(t)
	u := gjson.Get(res, "user")
	a.True(u.Exists())
	a.Equal(expected.Email, u.Get("email").String())
	a.NotEmpty(u.Get("token").String())
	a.Equal(expected.Name, u.Get("username").String())
	if isSignUp {
		a.Empty(u.Get("bio").String())
		a.Empty(u.Get("image").String())
	} else {
		a.Equal(expected.Bio, u.Get("bio").String())
		a.Equal(expected.Image, u.Get("image").String())
	}
}

func assertProfileResponse(t *testing.T, res string, expected *userModel.User, follow bool) {
	a := assert.New(t)
	p := gjson.Get(res, "profile")
	a.True(p.Exists())
	a.Equal(expected.Name, p.Get("username").String())
	a.Equal(expected.Bio, p.Get("bio").String())
	a.Equal(expected.Image, p.Get("image").String())
	a.Equal(follow, p.Get("following").Bool())
}

func assertErrorResponse(t *testing.T, actual *httptest.ResponseRecorder, code int, msg string) {
	a := assert.New(t)
	a.Equal(code, actual.Code)
	a.Contains(gjson.Get(actual.Body.String(), "errors.body").String(), msg)
}

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

func copyUser(u *userModel.User) *userModel.User {
	copied := &userModel.User{}
	copier.Copy(copied, u)
	return copied
}

func toJsonReader(m map[string]interface{}) io.Reader {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(b)
}

func keyValuesToMapIfNotEmpty(keyValues ...string) (map[string]interface{}, error) {
	if len(keyValues)%2 != 0 {
		return nil, errors.New("require key value pairs")
	}

	m := make(map[string]interface{}, len(keyValues)/2)
	for i := 0; i < len(keyValues); i += 2 {
		k := keyValues[i]
		v := keyValues[i+1]
		if v != "" {
			m[k] = v
		}
	}
	return m, nil
}
