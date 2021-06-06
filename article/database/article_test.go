package database

import (
	"context"
	"github.com/stretchr/testify/suite"
	"github.com/zacscoding/echo-gorm-realworld-app/article/model"
	"github.com/zacscoding/echo-gorm-realworld-app/database"
	"github.com/zacscoding/echo-gorm-realworld-app/logging"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/user/model"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"testing"
)

type Suite struct {
	suite.Suite
	db         ArticleDB
	originDB   *gorm.DB
	dbTeardown database.CloseFunc
	u1, u2, u3 *userModel.User
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
	s.db = NewArticleDB(s.originDB)
}

func (s *Suite) TearDownSuite() {
	s.dbTeardown()
}

func (s *Suite) SetupTest() {
	err := database.DeleteRecordAll(s.T(), s.originDB, []string{
		model.TableNameArticleFavorite, "user_id > 0",
		model.TableNameArticleTag, "article_id > 0",
		model.TableNameArticle, "article_id > 0",
		model.TableNameTag, "tag_id > 0",
		userModel.TableNameFollow, "user_id > 0",
		userModel.TableNameUser, "user_id > 0",
	})
	s.NoError(err)

	s.u1 = &userModel.User{Email: "user1@email.com", Name: "user1", Password: "user1password", Bio: "user1 bio", Image: "user1 image"}
	s.u2 = &userModel.User{Email: "user2@email.com", Name: "user2", Password: "user2password", Bio: "user2 bio", Image: "user2 image"}
	s.u3 = &userModel.User{Email: "user3@email.com", Name: "user3", Password: "user3password", Bio: "user3 bio", Image: "user3 image"}
	s.NoError(s.originDB.Create([]*userModel.User{s.u1, s.u2, s.u3}).Error)
}

func (s *Suite) TestSave() {
	s.NoError(s.originDB.Create([]*model.Tag{{Name: "exist1"}}).Error)
	a := &model.Article{
		Slug:        "article1",
		Title:       "article1",
		Description: "description",
		Body:        "body",
		Author:      *s.u1,
		Tags: []*model.Tag{
			{Name: "exist1"}, {Name: "newTag1"},
		},
	}

	// when
	err := s.db.Save(context.TODO(), a)
	// then
	s.NoError(err)
	find, err := s.db.FindBySlug(context.TODO(), nil, a.Slug)
	s.NoError(err)
	s.Equal(a.Slug, find.Slug)
	s.Equal(a.Title, find.Title)
	s.Equal(a.Description, find.Description)
	s.Equal(a.Body, find.Body)
	s.Equal(a.Author.ID, find.AuthorID)
	var tagCount int64
	s.NoError(s.originDB.Model(new(model.Tag)).Count(&tagCount).Error)
	s.EqualValues(2, tagCount)
}
