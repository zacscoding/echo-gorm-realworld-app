package user

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/zacscoding/echo-gorm-realworld-app/config"
	"github.com/zacscoding/echo-gorm-realworld-app/serverenv"
	userMocks "github.com/zacscoding/echo-gorm-realworld-app/user/database/mocks"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/httputils"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TODO: add tests after error handler

func Test1(t *testing.T) {
	cfg, err := config.Load("")
	assert.NoError(t, err)
	env := serverenv.NewServerEnv(serverenv.WithUserDB(&userMocks.UserDB{}))
	assert.NoError(t, err)
	h, err := NewHandler(env, cfg)
	assert.NoError(t, err)

	e := echo.New()
	e.Validator = httputils.NewValidator()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err = h.handleSignUp(c)
	fmt.Println(">>>>>")
	fmt.Println("err:", err)
	fmt.Println("body:", rec.Body.String())
}
