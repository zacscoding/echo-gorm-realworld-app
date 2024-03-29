package database

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/config"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/database"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/user/model"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/logging"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"math"
	"sync/atomic"
	"testing"
	"time"
)

var idx = int32(0)

var (
	defaultUser  = newTestUser("default@gmail.com", false)
	defaultUser2 = newTestUser("default2@gmail.com", false)
	defaultUser3 = newTestUser("default3@gmail.com", false)
	defaultUser4 = newTestUser("default4@gmail.com", false)
	defaultUser5 = newTestUser("default5@gmail.com", false)
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
	cfg, _ := config.Load("")
	logging.SetConfig(&logging.Config{
		Encoding:    "console",
		Level:       zapcore.FatalLevel,
		Development: false,
	})
	s.originDB, s.dbTeardown = database.NewTestDatabase(s.T(), true)
	s.db = NewUserDB(cfg, s.originDB)
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

	users := []*model.User{defaultUser, defaultUser2, defaultUser3, defaultUser4, defaultUser5, disabledUser}
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

func (s *Suite) TestFindByID() {
	cases := []struct {
		name     string
		id       uint
		assertFn func(t *testing.T, u *model.User)
		msg      string
	}{
		{
			name: "exist",
			id:   defaultUser.ID,
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
			name: "not found",
			id:   math.MaxInt8,
			msg:  database.ErrRecordNotFound.Error(),
		},
	}

	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			u, err := s.db.FindByID(context.TODO(), tc.id)
			if tc.msg != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.msg)
				return
			}
			tc.assertFn(t, u)
		})
	}
}

func (s *Suite) TestFindByUsername() {
	cases := []struct {
		name     string
		username string
		assertFn func(t *testing.T, u *model.User)
		msg      string
	}{
		{
			name:     "exist",
			username: defaultUser.Name,
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
			name:     "not found",
			username: "notfounduser@@",
			msg:      database.ErrRecordNotFound.Error(),
		},
	}

	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			u, err := s.db.FindByName(context.TODO(), tc.username)
			if tc.msg != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.msg)
				return
			}
			tc.assertFn(t, u)
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

func (s *Suite) TestIsFollow() {
	u1, u2 := defaultUser, defaultUser2
	s.NoError(s.db.Follow(context.TODO(), u1.ID, u2.ID))

	following, err := s.db.IsFollow(context.TODO(), u1.ID, u2.ID)
	s.NoError(err)
	s.True(following)

	following, err = s.db.IsFollow(context.TODO(), u2.ID, u1.ID)
	s.NoError(err)
	s.False(following)
}

func (s *Suite) TestIsFollows() {
	u1, u2, u3, u4, u5 := defaultUser, defaultUser2, defaultUser3, defaultUser4, defaultUser5
	// u1 follows u2, u4, u5
	s.NoError(s.db.Follow(context.TODO(), u1.ID, u2.ID))
	s.NoError(s.db.Follow(context.TODO(), u1.ID, u4.ID))
	s.NoError(s.db.Follow(context.TODO(), u1.ID, u5.ID))
	s.NoError(s.db.Follow(context.TODO(), u2.ID, u5.ID))
	s.NoError(s.db.Follow(context.TODO(), u2.ID, u1.ID))

	m, err := s.db.IsFollows(context.TODO(), u1.ID, []uint{
		u2.ID, u3.ID, u4.ID,
	})

	s.NoError(err)
	s.Len(m, 3)
	s.True(m[u2.ID])
	s.False(m[u3.ID])
	s.True(m[u4.ID])
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

func (s *Suite) TestFindFollowerIDs() {
	u1, u2, u3, u4, u5 := defaultUser, defaultUser2, defaultUser3, defaultUser4, defaultUser5
	// u1 follows u2, u4, u5
	s.NoError(s.db.Follow(context.TODO(), u1.ID, u2.ID))
	s.NoError(s.db.Follow(context.TODO(), u1.ID, u4.ID))
	s.NoError(s.db.Follow(context.TODO(), u1.ID, u5.ID))
	s.NoError(s.db.Follow(context.TODO(), u2.ID, u5.ID))
	s.NoError(s.db.Follow(context.TODO(), u2.ID, u1.ID))
	s.NoError(s.db.Follow(context.TODO(), u3.ID, u1.ID))

	followers, err := s.db.FindFollowerIDs(context.TODO(), u1.ID)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), followers, 3)
	assert.Contains(s.T(), followers, u2.ID)
	assert.Contains(s.T(), followers, u4.ID)
	assert.Contains(s.T(), followers, u5.ID)

	followers, err = s.db.FindFollowerIDs(context.TODO(), u5.ID)

	assert.NoError(s.T(), err)
	assert.Empty(s.T(), followers)
}

func newTestUser(email string, disabled bool) *model.User {
	idx := atomic.AddInt32(&idx, 1)
	return &model.User{
		Email:     email,
		Name:      fmt.Sprintf("zac-%d", idx),
		Password:  fmt.Sprintf("password%d", idx),
		Bio:       fmt.Sprintf("working during %d days...", idx),
		Image:     fmt.Sprintf("https://mycdn/profile-%d.png", idx),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Disabled:  disabled,
	}
}
