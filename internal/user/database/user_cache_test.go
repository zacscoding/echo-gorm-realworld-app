package database

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/cache"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/config"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/user/database/mocks"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/internal/user/model"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/logging"
	"go.uber.org/zap/zapcore"
	"testing"
)

type CacheSuite struct {
	suite.Suite
	cacheDB    UserDB
	cacheClose cache.CloseFunc
	dbMock     *mocks.UserDB
}

func TestCacheSuite(t *testing.T) {
	suite.Run(t, new(CacheSuite))
}

func (s *CacheSuite) SetupTest() {
	conf, _ := config.Load("")
	logging.SetConfig(&logging.Config{
		Encoding:    "console",
		Level:       zapcore.FatalLevel,
		Development: false,
	})

	// setup database mock
	s.dbMock = &mocks.UserDB{}

	cli, _, closeFn := cache.NewTestCache(s.T())
	s.cacheDB = NewUserCacheDB(conf, cli, s.dbMock)
	s.cacheClose = closeFn
}

func (s *CacheSuite) TearDownTest() {
	if s.cacheClose != nil {
		s.cacheClose()
	}
}

func (s *CacheSuite) TestSave() {
	u := defaultUser
	s.dbMock.On("Save", mock.Anything, mock.Anything).Return(nil)

	err := s.cacheDB.Save(context.TODO(), u)

	s.NoError(err)
	s.dbMock.AssertCalled(s.T(), "Save", mock.Anything, u)
	_, err = s.cacheDB.FindByID(context.TODO(), u.ID)
	s.NoError(err)
	s.dbMock.AssertNotCalled(s.T(), "FindByID", mock.Anything, mock.Anything)
}

func (s *CacheSuite) TestUpdateCacheIfExist() {
	u := defaultUser
	s.dbMock.On("Save", mock.Anything, mock.Anything).Return(nil)
	s.dbMock.On("Update", mock.Anything, mock.Anything).Return(nil)
	s.NoError(s.cacheDB.Save(context.TODO(), u))

	updateUser := userModel.User{
		ID:       u.ID,
		Email:    "updated@email.com",
		Name:     u.Name,
		Password: u.Password,
		Bio:      u.Bio,
		Image:    u.Image,
	}

	err := s.cacheDB.Update(context.TODO(), &updateUser)

	s.NoError(err)
	s.dbMock.AssertCalled(s.T(), "Update", mock.Anything, &updateUser)
	find, err := s.cacheDB.FindByID(context.TODO(), updateUser.ID)
	s.NoError(err)
	s.Equal(updateUser.Email, find.Email)
	s.dbMock.AssertNotCalled(s.T(), "FindByID", mock.Anything, mock.Anything)
}

func (s *CacheSuite) TestUpdateIgnoreIfNotExist() {
	u := defaultUser
	s.dbMock.On("Update", mock.Anything, mock.Anything).Return(nil)
	s.dbMock.On("FindByID", mock.Anything, u.ID).Return(u, nil)

	err := s.cacheDB.Update(context.TODO(), u)

	s.NoError(err)
	s.cacheDB.FindByID(context.TODO(), u.ID)
	s.dbMock.AssertCalled(s.T(), "FindByID", mock.Anything, u.ID)
}

func (s *CacheSuite) TestFindByNameNoCache() {
	username := "user1"
	s.dbMock.On("FindByName", mock.Anything, username).Return(defaultUser, nil)

	_, err := s.cacheDB.FindByName(context.TODO(), username)

	s.NoError(err)
	s.Empty(s.getCacheKeys())
	s.dbMock.AssertCalled(s.T(), "FindByName", mock.Anything, username)
}

func (s *CacheSuite) TestFindByEmailNoCache() {
	email := "user1@gmail.com"
	s.dbMock.On("FindByEmail", mock.Anything, email).Return(defaultUser, nil)

	_, err := s.cacheDB.FindByEmail(context.TODO(), email)

	s.NoError(err)
	s.Empty(s.getCacheKeys())
	s.dbMock.AssertCalled(s.T(), "FindByEmail", mock.Anything, email)
}

func (s *CacheSuite) TestFollowNoCache() {
	userID, followerID := uint(1), uint(2)
	s.dbMock.On("Follow", mock.Anything, userID, followerID).Return(nil)

	err := s.cacheDB.Follow(context.TODO(), userID, followerID)

	s.NoError(err)
	s.Empty(s.getCacheKeys())
	s.dbMock.AssertCalled(s.T(), "Follow", mock.Anything, userID, followerID)
}

func (s *CacheSuite) TestIsFollowNoCache() {
	userID, followerID := uint(1), uint(2)
	s.dbMock.On("IsFollow", mock.Anything, userID, followerID).Return(false, nil)

	_, err := s.cacheDB.IsFollow(context.TODO(), userID, followerID)

	s.NoError(err)
	s.Empty(s.getCacheKeys())
	s.dbMock.AssertCalled(s.T(), "IsFollow", mock.Anything, userID, followerID)
}

func (s *CacheSuite) TestIsFollowsNoCache() {
	userID, followerIDs := uint(1), []uint{2, 3}
	s.dbMock.On("IsFollows", mock.Anything, userID, followerIDs).Return(map[uint]bool{
		2: true,
	}, nil)

	_, err := s.cacheDB.IsFollows(context.TODO(), userID, followerIDs)

	s.NoError(err)
	s.Empty(s.getCacheKeys())
	s.dbMock.AssertCalled(s.T(), "IsFollows", mock.Anything, userID, followerIDs)
}

func (s *CacheSuite) TestUnFollowNoCache() {
	userID, followerID := uint(1), uint(2)
	s.dbMock.On("UnFollow", mock.Anything, userID, followerID).Return(nil)

	err := s.cacheDB.UnFollow(context.TODO(), userID, followerID)

	s.NoError(err)
	s.Empty(s.getCacheKeys())
	s.dbMock.AssertCalled(s.T(), "UnFollow", mock.Anything, userID, followerID)
}

func (s *CacheSuite) TestFindFollowerIDsNoCache() {
	userID := uint(1)
	s.dbMock.On("FindFollowerIDs", mock.Anything, userID).Return([]uint{2}, nil)

	_, err := s.cacheDB.FindFollowerIDs(context.TODO(), userID)

	s.NoError(err)
	s.Empty(s.getCacheKeys())
	s.dbMock.AssertCalled(s.T(), "FindFollowerIDs", mock.Anything, userID)
}

func (s *CacheSuite) getCacheKeys() []string {
	cli := s.cacheDB.(*userCache).cli
	return cli.Keys(context.Background(), "*").Val()
}
