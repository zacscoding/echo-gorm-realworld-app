package database

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zacscoding/echo-gorm-realworld-app/database"
	"github.com/zacscoding/echo-gorm-realworld-app/logging"
	"github.com/zacscoding/echo-gorm-realworld-app/user/model"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"math/rand"
	"testing"
	"time"
)

var (
	defaultUser  = newTestUser("default@gmail.com", false)
	defaultUser2 = newTestUser("default2@gmail.com", false)
	disabledUser = newTestUser("disable@gmail.com", true)
)

type Suite struct {
	suite.Suite
	db         UserDB
	originDB   *gorm.DB
	dbTeardown database.CloseFunc
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) SetupSuite() {
	logging.SetConfig(&logging.Config{
		Encoding:    "console",
		Level:       zapcore.FatalLevel,
		Development: false,
	})
	s.originDB, s.dbTeardown = database.NewTestDatabase(s.T(), true)
	s.db = NewUserDB(s.originDB)
}

func (s *Suite) TearDownSuite() {
	s.dbTeardown()
}

func (s *Suite) SetupTest() {
	err := database.DeleteRecordAll(s.T(), s.originDB, []string{
		model.TableNameFollow, "user_id > 0 AND follow_id > 0",
		model.TableNameUser, "user_id > 0",
	})
	s.NoError(err)

	users := []*model.User{defaultUser, defaultUser2, disabledUser}
	s.NoError(s.originDB.Create(users).Error)
}

func (s *Suite) TestSave() {
	cases := []struct {
		name string
		u    model.User
		msg  string
	}{
		{
			name: "success",
			u: model.User{
				Email:     "user1@gmail.com",
				Name:      "zac",
				Password:  "password!",
				Bio:       "test bio",
				Image:     "https://cdn/1.png",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Disabled:  false,
			},
		}, {
			name: "duplicate",
			u: model.User{
				Email:     defaultUser.Email,
				Name:      "zac",
				Password:  "password!",
				Bio:       "test bio",
				Image:     "https://cdn/1.png",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Disabled:  false,
			},
			msg: "conflict key",
		},
	}

	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			err := s.db.Save(context.TODO(), &tc.u)
			if tc.msg != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), err.Error())
				return
			}
			assert.NoError(t, err)
		})
	}
}

func (s *Suite) TestFindByEmail() {
	cases := []struct {
		name     string
		email    string
		assertFn func(t *testing.T, u *model.User)
		msg      string
	}{
		{
			name:  "exist",
			email: defaultUser.Email,
			assertFn: func(t *testing.T, u *model.User) {
				assert.Greater(t, u.ID, uint(0))
				assert.Equal(t, defaultUser.Name, u.Name)
				assert.Equal(t, defaultUser.Password, u.Password)
				assert.Equal(t, defaultUser.Bio, u.Bio)
				assert.Equal(t, defaultUser.Image, u.Image)
				assert.WithinDuration(t, defaultUser.CreatedAt, u.CreatedAt, time.Minute)
				assert.WithinDuration(t, defaultUser.UpdatedAt, u.UpdatedAt, time.Minute)
				assert.False(t, u.Disabled)
			},
		}, {
			name:  "not found",
			email: "notfound@gmail.com",
			msg:   database.ErrRecordNotFound.Error(),
		},
	}

	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			u, err := s.db.FindByEmail(context.TODO(), tc.email)
			if tc.msg != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.msg)
				return
			}
			tc.assertFn(t, u)
		})
	}
}

func (s *Suite) TestUpdate() {
	cases := []struct {
		name     string
		update   *model.User
		assertFn func(t *testing.T, find *model.User, findErr error)
		msg      string
	}{
		{
			name: "success",
			update: &model.User{
				ID:       defaultUser.ID,
				Email:    "updated@email.com",
				Name:     "updatedName",
				Password: "updatedPassword",
				Bio:      "updatedBio",
				Image:    "updatedImage",
				Disabled: false,
			},
			assertFn: func(t *testing.T, find *model.User, findErr error) {
				assert.NoError(t, findErr)
				assert.Equal(t, "updated@email.com", find.Email)
				assert.Equal(t, "updatedName", find.Name)
				assert.Equal(t, "updatedPassword", find.Password)
				assert.Equal(t, "updatedBio", find.Bio)
				assert.Equal(t, "updatedImage", find.Image)
				assert.False(t, find.Disabled)
			},
		}, {
			name: "not found",
			update: &model.User{
				ID:    100000,
				Email: "updated@email.com",
			},
			msg: database.ErrRecordNotFound.Error(),
		}, {
			name: "duplicate email",
			update: &model.User{
				ID:    defaultUser2.ID,
				Email: disabledUser.Email,
			},
			msg: database.ErrKeyConflict.Error(),
		},
	}

	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			err := s.db.Update(context.TODO(), tc.update)
			if tc.msg != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.msg)
				return
			}

			assert.NoError(t, err)
			find, err := s.db.FindByEmail(context.TODO(), tc.update.Email)
			tc.assertFn(t, find, err)
		})
	}
}

func (s *Suite) TestFollow() {
	u1, u2 := defaultUser, defaultUser2

	err := s.db.Follow(context.TODO(), u1.ID, u2.ID)
	s.NoError(err)

	err = s.db.Follow(context.TODO(), u1.ID, u2.ID)
	s.Equal(database.ErrKeyConflict, err)

	err = s.db.Follow(context.TODO(), u1.ID, 10000)
	s.Equal(database.ErrFKConstraint, err)
}

func (s *Suite) TestUnFollow() {
	u1, u2 := defaultUser, defaultUser2
	err := s.db.Follow(context.TODO(), u1.ID, u2.ID)
	s.NoError(err)

	err = s.db.UnFollow(context.TODO(), u1.ID, u2.ID)
	assert.NoError(s.T(), err)

	err = s.db.UnFollow(context.TODO(), u1.ID, u2.ID)
	fmt.Println(err)
}

func newTestUser(email string, disabled bool) *model.User {
	randIdx := rand.Intn(100)
	return &model.User{
		Email:     email,
		Name:      fmt.Sprintf("zac-%d", randIdx),
		Password:  fmt.Sprintf("password%d", randIdx),
		Bio:       fmt.Sprintf("working during %d days...", randIdx),
		Image:     fmt.Sprintf("https://mycdn/profile-%d.png", randIdx),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Disabled:  disabled,
	}
}
