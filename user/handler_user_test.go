package user

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/authutils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleSignUp(t *testing.T) {
	cases := []struct {
		name        string
		requestBody string
		assertFunc  func(t *testing.T, rec *httptest.ResponseRecorder, err error)
	}{
		{
			name:        "success",
			requestBody: `{"user": {"username": "user1", "email": "user1@gmail.com", "password": "user1"}}`,
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				assert.NoError(t, err)
				assertUserResponse(t, rec.Body.String(), defaultUser1, true)
			},
		}, {
			name:        "duplicate email",
			requestBody: `{"user": {"username": "disabledUser", "email": "disabledUser@gmail.com", "password": "disabledUser"}}`,
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				assertErrorResponse(t, http.StatusUnprocessableEntity, "duplicate email", err)
			},
		},
	}

	for _, tc := range cases {
		// given
		e, h, _ := fixtures(t)
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.requestBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// when
		err := h.handleSignUp(c)

		// then
		tc.assertFunc(t, rec, err)
	}
}

func TestHandleSignIn(t *testing.T) {
	cases := []struct {
		name        string
		requestBody string
		assertFunc  func(t *testing.T, rec *httptest.ResponseRecorder, err error)
	}{
		{
			name:        "success",
			requestBody: `{"user": {"email": "user1@gmail.com", "password": "user1"}}`,
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				assert.NoError(t, err)
				assertUserResponse(t, rec.Body.String(), defaultUser1, false)
			},
		}, {
			name:        "mismatch password",
			requestBody: `{"user": {"email": "user1@gmail.com", "password": "user1@"}}`,
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				assertErrorResponse(t, http.StatusUnprocessableEntity, "password mismatch", err)
			},
		}, {
			name:        "not found",
			requestBody: `{"user": {"email": "disabledUser@gmail.com", "password": "disabledUser"}}`,
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				assertErrorResponse(t, http.StatusNotFound, "user(disabledUser@gmail.com) not found", err)
			},
		}, {
			name:        "internal server error",
			requestBody: `{"user": {"email": "user10@gmail.com", "password": "user10"}}`,
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				assertErrorResponse(t, http.StatusInternalServerError, "force error", err)
			},
		}, {
			name:        "validation no email",
			requestBody: `{"user": {"password": "user10"}}`,
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				assertErrorResponse(t, http.StatusUnprocessableEntity, "Email validation error. reason: required", err)
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			e, h, _ := fixtures(t)
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// when
			err := h.handleSignIn(c)

			// then
			tc.assertFunc(t, rec, err)
		})
	}
}

func TestHandleCurrentUser(t *testing.T) {
	cases := []struct {
		name          string
		currentUserID uint
		authToken     string
		assertFunc    func(t *testing.T, rec *httptest.ResponseRecorder, err error)
	}{
		{
			name:          "success",
			currentUserID: defaultUser1.ID,
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				assert.NoError(t, err)
				assertUserResponse(t, rec.Body.String(), defaultUser1, false)
			},
		}, {
			name:          "internal server error",
			currentUserID: disabledUser.ID,
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, err error) {
				assertErrorResponse(t, http.StatusInternalServerError, "record not found", err)
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// given
			e, h, _ := fixtures(t)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.authToken != "" {
				req.Header.Set("Authorization", "Bearer "+tc.authToken)
			}

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if tc.currentUserID != 0 {
				token := new(jwt.Token)
				token.Claims = &authutils.JWTClaims{
					UserID:         tc.currentUserID,
					StandardClaims: jwt.StandardClaims{},
				}
				c.Set("user", token)
			}

			// when
			err := h.handleCurrentUser(c)

			// then
			tc.assertFunc(t, rec, err)
		})
	}
}
