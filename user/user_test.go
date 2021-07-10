package user

import (
	"errors"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tidwall/gjson"
	"github.com/zacscoding/echo-gorm-realworld-app/database"
	userMocks "github.com/zacscoding/echo-gorm-realworld-app/user/database/mocks"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/user/model"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/hashutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func (s *TestSuite) TestHandleSignUp() {
	s.u.On("Save", mock.Anything, mock.Anything).Return(nil)

	// when
	uri := "/api/users"
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, uri, toJsonReader(map[string]interface{}{
		"user": map[string]interface{}{
			"username": defaultUsers[0].Name,
			"email":    defaultUsers[0].Email,
			"password": defaultUsers[0].Name,
		},
	}))
	req.Header.Set("Content-Type", "application/json")

	s.e.ServeHTTP(rec, req)

	// then
	s.u.AssertCalled(s.T(), "Save", mock.Anything, mock.MatchedBy(func(u *userModel.User) bool {
		if hashutils.MatchesPassword(u.Password, defaultUsers[0].Name) != nil {
			return false
		}
		return u.Email == defaultUsers[0].Email &&
			u.Name == defaultUsers[0].Name
	}))
	s.Equal(http.StatusOK, rec.Code)
	assertUserResponse(s.T(), rec.Body.String(), defaultUsers[0], true)
}

func (s *TestSuite) TestHandleSignUp_BindError() {
	cases := []struct {
		name     string
		username string
		email    string
		password string
		// expected
		code int
		msg  string
	}{
		{
			name:     "missing username",
			email:    "user@gmail.com",
			password: "pass",
			code:     http.StatusUnprocessableEntity,
			msg:      "Username validation error. reason: required",
		}, {
			name:     "missing email",
			username: "user1",
			password: "pass",
			code:     http.StatusUnprocessableEntity,
			msg:      "Email validation error. reason: required",
		}, {
			name:     "missing password",
			username: "user1",
			email:    "user1@gmail.com",
			code:     http.StatusUnprocessableEntity,
			msg:      "Password validation error. reason: required",
		}, {
			name:     "invalid email",
			username: "user1",
			email:    "not email pattern",
			password: "pass",
			code:     http.StatusUnprocessableEntity,
			msg:      "Email validation error. reason: email",
		},
	}

	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			uri := "/api/users"
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, uri, toJsonReader(map[string]interface{}{
				"user": map[string]interface{}{
					"username": tc.username,
					"email":    tc.email,
					"password": tc.password,
				},
			}))
			req.Header.Set("Content-Type", "application/json")

			s.e.ServeHTTP(rec, req)

			assert.Equal(t, tc.code, rec.Code)
			assert.Contains(t, gjson.Get(rec.Body.String(), "errors.body").String(), tc.msg)
		})
	}
}

func (s *TestSuite) TestHandleSignUp_Fail() {
	cases := []struct {
		name      string
		setupMock func(m *userMocks.UserDB)
		// expected
		code int
		msg  string
	}{
		{
			name: "duplicate email",
			setupMock: func(m *userMocks.UserDB) {
				m.On("Save", mock.Anything, mock.Anything).Return(database.ErrKeyConflict)
			},
			code: http.StatusUnprocessableEntity,
			msg:  "duplicate email",
		}, {
			name: "any error",
			setupMock: func(m *userMocks.UserDB) {
				m.On("Save", mock.Anything, mock.Anything).Return(errors.New("force error"))
			},
			code: http.StatusInternalServerError,
			msg:  "force error",
		},
	}

	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			s.resetMocks()
			tc.setupMock(s.u)

			req, _ := http.NewRequest(http.MethodPost, "/api/users", toJsonReader(map[string]interface{}{
				"user": map[string]interface{}{
					"username": defaultUsers[0].Name,
					"email":    defaultUsers[0].Email,
					"password": defaultUsers[0].Name,
				},
			}))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			s.e.ServeHTTP(rec, req)

			assert.Equal(t, tc.code, rec.Code)
			assert.Contains(t, gjson.Get(rec.Body.String(), "errors.body").String(), tc.msg)
		})
	}
}

func (s *TestSuite) TestHandleSignIn() {
	s.u.On("FindByEmail", mock.Anything, defaultUsers[0].Email).Return(copyUser(defaultUsers[0]), nil)

	// when
	uri := "/api/users/login"
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, uri, toJsonReader(map[string]interface{}{
		"user": map[string]interface{}{
			"email":    defaultUsers[0].Email,
			"password": defaultUsers[0].Name,
		},
	}))
	req.Header.Set("Content-Type", "application/json")

	s.e.ServeHTTP(rec, req)

	// then
	s.u.AssertCalled(s.T(), "FindByEmail", mock.Anything, defaultUsers[0].Email)
	s.Equal(http.StatusOK, rec.Code)
	assertUserResponse(s.T(), rec.Body.String(), defaultUsers[0], false)
}

func (s *TestSuite) TestHandleSignIn_BindingError() {
	cases := []struct {
		name     string
		email    string
		password string
		// expected
		code int
		msg  string
	}{
		{
			name:     "missing email",
			password: "pass",
			code:     http.StatusUnprocessableEntity,
			msg:      "Email validation error. reason: required",
		}, {
			name:     "not email pattern",
			email:    "notemail",
			password: "pass",
			code:     http.StatusUnprocessableEntity,
			msg:      "Email validation error. reason: email",
		}, {
			name:  "missing password",
			email: defaultUsers[0].Email,
			code:  http.StatusUnprocessableEntity,
			msg:   "Password validation error. reason: required",
		},
	}

	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			uri := "/api/users/login"
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, uri, toJsonReader(map[string]interface{}{
				"user": map[string]interface{}{
					"email":    tc.email,
					"password": tc.password,
				},
			}))
			req.Header.Set("Content-Type", "application/json")

			s.e.ServeHTTP(rec, req)

			assert.Equal(t, tc.code, rec.Code)
			assert.Contains(t, gjson.Get(rec.Body.String(), "errors.body").String(), tc.msg)
		})
	}
}

func (s *TestSuite) TestHandleSignIn_Fail() {
	cases := []struct {
		name      string
		email     string
		password  string
		setupMock func(m *userMocks.UserDB)
		// expected
		code int
		msg  string
	}{
		{
			name:     "user not found",
			email:    defaultUsers[0].Email,
			password: "invalid password",
			setupMock: func(m *userMocks.UserDB) {
				m.On("FindByEmail", mock.Anything, defaultUsers[0].Email).Return(nil, database.ErrRecordNotFound)
			},
			code: http.StatusNotFound,
			msg:  "not found",
		}, {
			name:     "mismatch password",
			email:    defaultUsers[0].Email,
			password: "invalid password",
			setupMock: func(m *userMocks.UserDB) {
				m.On("FindByEmail", mock.Anything, defaultUsers[0].Email).Return(copyUser(defaultUsers[0]), nil)
			},
			code: http.StatusUnprocessableEntity,
			msg:  "password mismatch",
		}, {
			name:     "any error",
			email:    defaultUsers[0].Email,
			password: "invalid password",
			setupMock: func(m *userMocks.UserDB) {
				m.On("FindByEmail", mock.Anything, defaultUsers[0].Email).Return(nil, errors.New("force error"))
			},
			code: http.StatusInternalServerError,
			msg:  "force error",
		},
	}

	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			s.resetMocks()
			tc.setupMock(s.u)
			uri := "/api/users/login"
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, uri, toJsonReader(map[string]interface{}{
				"user": map[string]interface{}{
					"email":    tc.email,
					"password": tc.password,
				},
			}))
			req.Header.Set("Content-Type", "application/json")

			s.e.ServeHTTP(rec, req)

			assert.Equal(t, tc.code, rec.Code)
			assert.Contains(t, gjson.Get(rec.Body.String(), "errors.body").String(), tc.msg)
		})
	}
}

func (s *TestSuite) TestHandleCurrentUser() {
	s.u.On("FindByID", mock.Anything, defaultUsers[0].ID).Return(copyUser(defaultUsers[0]), nil)

	// when
	uri := "/api/user"
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	token, _ := s.h.makeJWTToken(defaultUsers[0])
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	s.e.ServeHTTP(rec, req)

	// then
	s.u.AssertCalled(s.T(), "FindByID", mock.Anything, defaultUsers[0].ID)
	s.Equal(http.StatusOK, rec.Code)
	assertUserResponse(s.T(), rec.Body.String(), defaultUsers[0], false)
}

func (s *TestSuite) TestHandleCurrentUser_BindError() {
	// TODO
}

func (s *TestSuite) TestHandleCurrentUser_Fail() {
	// TODO
}

func (s *TestSuite) TestHandleUpdateUser() {
	updatedUser := &userModel.User{}
	copier.Copy(updatedUser, defaultUsers[0])
	s.u.On("FindByID", mock.Anything, defaultUsers[0].ID).Return(copyUser(defaultUsers[0]), nil)
	s.u.On("Update", mock.Anything, mock.Anything).Return(nil)
	updatedUser.Name = "updated-user"
	updatedUser.Bio = "updated-bio"
	updatedUser.Image = "updated-image"

	// when
	uri := "/api/user"
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, uri, toJsonReader(map[string]interface{}{
		"user": map[string]interface{}{
			"username": updatedUser.Name,
			"bio":      updatedUser.Bio,
			"image":    updatedUser.Image,
		},
	}))
	token, _ := s.h.makeJWTToken(updatedUser)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	s.e.ServeHTTP(rec, req)

	// then
	s.u.AssertCalled(s.T(), "FindByID", mock.Anything, defaultUsers[0].ID)
	s.u.AssertCalled(s.T(), "Update", mock.Anything, mock.Anything)
	s.Equal(http.StatusOK, rec.Code)
	assertUserResponse(s.T(), rec.Body.String(), updatedUser, false)
}

func (s *TestSuite) TestHandleUpdateUser_BindError() {
	// TODO
}

func (s *TestSuite) TestHandleUpdateUser_Fail() {
	// TODO
}
