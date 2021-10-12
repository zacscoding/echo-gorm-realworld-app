package user

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/database"
	userMocks "github.com/zacscoding/echo-gorm-realworld-app/internal/user/database/mocks"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/internal/user/model"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/utils/authutils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func (s *TestSuite) TestHandleGetProfile() {
	cases := []struct {
		name        string
		username    string
		currentUser *userModel.User
		setupMock   func(m *userMocks.UserDB)
		assertFunc  func(t *testing.T, rec *httptest.ResponseRecorder, m *userMocks.UserDB)
	}{
		{
			name:        "profile with follow",
			username:    defaultUsers[1].Name,
			currentUser: defaultUsers[0],
			setupMock: func(m *userMocks.UserDB) {
				m.On("FindByName", mock.Anything, defaultUsers[1].Name).Return(copyUser(defaultUsers[1]), nil)
				m.On("IsFollow", mock.Anything, defaultUsers[0].ID, defaultUsers[1].ID).Return(true, nil)
			},
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, m *userMocks.UserDB) {
				s.u.AssertCalled(s.T(), "FindByName", mock.Anything, defaultUsers[1].Name)
				s.u.AssertCalled(s.T(), "IsFollow", mock.Anything, defaultUsers[0].ID, defaultUsers[1].ID)

				assert.Equal(t, http.StatusOK, rec.Code)
				assertProfileResponse(t, rec.Body.String(), defaultUsers[1], true)
			},
		}, {
			name:     "profile with anonymous",
			username: defaultUsers[1].Name,
			setupMock: func(m *userMocks.UserDB) {
				m.On("FindByName", mock.Anything, defaultUsers[1].Name).Return(copyUser(defaultUsers[1]), nil)
			},
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, m *userMocks.UserDB) {
				s.u.AssertCalled(s.T(), "FindByName", mock.Anything, defaultUsers[1].Name)
				s.u.AssertNotCalled(t, "IsFollow")

				assert.Equal(t, http.StatusOK, rec.Code)
				assertProfileResponse(t, rec.Body.String(), defaultUsers[1], false)
			},
		}, {
			name:     "user not found",
			username: defaultUsers[1].Name,
			setupMock: func(m *userMocks.UserDB) {
				m.On("FindByName", mock.Anything, defaultUsers[1].Name).Return(nil, database.ErrRecordNotFound)
			},
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, m *userMocks.UserDB) {
				s.u.AssertCalled(s.T(), "FindByName", mock.Anything, defaultUsers[1].Name)
				s.u.AssertNotCalled(t, "IsFollow")
				assertErrorResponse(t, rec, http.StatusNotFound, "not found")
			},
		},
	}

	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			if tc.setupMock != nil {
				s.resetMocks()
				tc.setupMock(s.u)
			}
			req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/profiles/%s", tc.username), nil)
			req.Header.Set("Content-Type", "application/json")
			if tc.currentUser != nil {
				token, _ := s.h.makeJWTToken(tc.currentUser)
				authutils.SetAuthToken(req, token)
			}
			rec := httptest.NewRecorder()

			s.e.ServeHTTP(rec, req)

			tc.assertFunc(t, rec, s.u)
		})
	}
}

func (s *TestSuite) TestHandleFollow() {
	cases := []struct {
		name        string
		username    string
		currentUser *userModel.User
		setupMock   func(m *userMocks.UserDB)
		assertFunc  func(t *testing.T, rec *httptest.ResponseRecorder, m *userMocks.UserDB)
	}{
		{
			name:        "success",
			username:    defaultUsers[1].Name,
			currentUser: defaultUsers[0],
			setupMock: func(m *userMocks.UserDB) {
				m.On("FindByName", mock.Anything, defaultUsers[1].Name).Return(copyUser(defaultUsers[1]), nil)
				m.On("Follow", mock.Anything, defaultUsers[0].ID, defaultUsers[1].ID).Return(nil)
			},
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, m *userMocks.UserDB) {
				s.u.AssertCalled(s.T(), "FindByName", mock.Anything, defaultUsers[1].Name)
				s.u.AssertCalled(s.T(), "Follow", mock.Anything, defaultUsers[0].ID, defaultUsers[1].ID)

				assert.Equal(t, http.StatusOK, rec.Code)
				assertProfileResponse(t, rec.Body.String(), defaultUsers[1], true)
			},
		}, {
			name:        "fail not found",
			username:    defaultUsers[1].Name,
			currentUser: defaultUsers[0],
			setupMock: func(m *userMocks.UserDB) {
				m.On("FindByName", mock.Anything, defaultUsers[1].Name).Return(nil, database.ErrRecordNotFound)
			},
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, m *userMocks.UserDB) {
				s.u.AssertCalled(s.T(), "FindByName", mock.Anything, defaultUsers[1].Name)

				assertErrorResponse(t, rec, http.StatusNotFound, "not found")
			},
		}, {
			name:        "fail already follow",
			username:    defaultUsers[1].Name,
			currentUser: defaultUsers[0],
			setupMock: func(m *userMocks.UserDB) {
				m.On("FindByName", mock.Anything, defaultUsers[1].Name).Return(copyUser(defaultUsers[1]), nil)
				m.On("Follow", mock.Anything, defaultUsers[0].ID, defaultUsers[1].ID).Return(database.ErrKeyConflict)
			},
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, m *userMocks.UserDB) {
				s.u.AssertCalled(s.T(), "FindByName", mock.Anything, defaultUsers[1].Name)
				s.u.On("Follow", mock.Anything, defaultUsers[0].ID, defaultUsers[1].ID)

				assertErrorResponse(t, rec, http.StatusUnprocessableEntity, "user already following")
			},
		}, {
			name:     "fail no auth",
			username: defaultUsers[1].Name,
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, m *userMocks.UserDB) {
				assertErrorResponse(t, rec, http.StatusUnauthorized, "auth required")
			},
		},
	}

	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			if tc.setupMock != nil {
				s.resetMocks()
				tc.setupMock(s.u)
			}
			req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/api/profiles/%s/follow", tc.username), nil)
			req.Header.Set("Content-Type", "application/json")
			if tc.currentUser != nil {
				token, _ := s.h.makeJWTToken(tc.currentUser)
				authutils.SetAuthToken(req, token)
			}
			rec := httptest.NewRecorder()

			s.e.ServeHTTP(rec, req)

			tc.assertFunc(t, rec, s.u)
		})
	}
}

func (s *TestSuite) TestHandleUnFollow() {
	cases := []struct {
		name        string
		username    string
		currentUser *userModel.User
		setupMock   func(m *userMocks.UserDB)
		assertFunc  func(t *testing.T, rec *httptest.ResponseRecorder, m *userMocks.UserDB)
	}{
		{
			name:        "success",
			username:    defaultUsers[1].Name,
			currentUser: defaultUsers[0],
			setupMock: func(m *userMocks.UserDB) {
				m.On("FindByName", mock.Anything, defaultUsers[1].Name).Return(copyUser(defaultUsers[1]), nil)
				m.On("UnFollow", mock.Anything, defaultUsers[0].ID, defaultUsers[1].ID).Return(nil)
			},
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, m *userMocks.UserDB) {
				s.u.AssertCalled(t, "FindByName", mock.Anything, defaultUsers[1].Name)
				s.u.AssertCalled(t, "UnFollow", mock.Anything, defaultUsers[0].ID, defaultUsers[1].ID)

				assert.Equal(t, http.StatusOK, rec.Code)
				assertProfileResponse(t, rec.Body.String(), defaultUsers[1], false)
			},
		}, {
			name:        "fail not found",
			username:    defaultUsers[1].Name,
			currentUser: defaultUsers[0],
			setupMock: func(m *userMocks.UserDB) {
				m.On("FindByName", mock.Anything, defaultUsers[1].Name).Return(nil, database.ErrRecordNotFound)
			},
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, m *userMocks.UserDB) {
				s.u.AssertCalled(t, "FindByName", mock.Anything, defaultUsers[1].Name)

				assertErrorResponse(t, rec, http.StatusNotFound, "not found")
			},
		}, {
			name:        "fail already unfollow",
			username:    defaultUsers[1].Name,
			currentUser: defaultUsers[0],
			setupMock: func(m *userMocks.UserDB) {
				m.On("FindByName", mock.Anything, defaultUsers[1].Name).Return(copyUser(defaultUsers[1]), nil)
				m.On("UnFollow", mock.Anything, defaultUsers[0].ID, defaultUsers[1].ID).Return(database.ErrRecordNotFound)
			},
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, m *userMocks.UserDB) {
				s.u.AssertCalled(t, "FindByName", mock.Anything, defaultUsers[1].Name)
				s.u.AssertCalled(t, "UnFollow", mock.Anything, defaultUsers[0].ID, defaultUsers[1].ID)

				assertErrorResponse(t, rec, http.StatusUnprocessableEntity, "user already unfollowing")
			},
		}, {
			name:     "fail no auth",
			username: defaultUsers[1].Name,
			assertFunc: func(t *testing.T, rec *httptest.ResponseRecorder, m *userMocks.UserDB) {
				assertErrorResponse(t, rec, http.StatusUnauthorized, "auth required")
			},
		},
	}

	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			if tc.setupMock != nil {
				s.resetMocks()
				tc.setupMock(s.u)
			}
			req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/profiles/%s/follow", tc.username), nil)
			req.Header.Set("Content-Type", "application/json")
			if tc.currentUser != nil {
				token, _ := s.h.makeJWTToken(tc.currentUser)
				authutils.SetAuthToken(req, token)
			}
			rec := httptest.NewRecorder()

			s.e.ServeHTTP(rec, req)

			tc.assertFunc(t, rec, s.u)
		})
	}
}
