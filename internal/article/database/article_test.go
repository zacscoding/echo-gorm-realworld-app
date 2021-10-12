package database

import (
	"context"
	"github.com/gosimple/slug"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/article/model"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/config"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/database"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/internal/user/model"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/logging"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"testing"
	"time"
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
	cfg, _ := config.Load("")
	logging.SetConfig(&logging.Config{
		Encoding:    "console",
		Level:       zapcore.FatalLevel,
		Development: false,
	})
	s.originDB, s.dbTeardown = database.NewTestDatabase(s.T(), true)
	s.db = NewArticleDB(cfg, s.originDB)
}

func (s *Suite) TearDownSuite() {
	s.dbTeardown()
}

func (s *Suite) SetupTest() {
	err := database.DeleteRecordAll(s.T(), s.originDB, []string{
		model.TableNameComment, "comment_id > 0",
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
	existTag := model.Tag{Name: "exist1"}
	s.NoError(s.originDB.Create(&existTag).Error)
	a := newArticle("article1", "description", "body", *s.u1, []string{existTag.Name, "newTag1"})
	now := time.Now()

	// when
	err := s.db.Save(context.TODO(), a)
	// then
	s.NoError(err)
	find, err := s.db.FindBySlug(context.TODO(), nil, a.Slug)
	s.NoError(err)
	s.Equal(a.Slug, find.Slug)
	s.Equal(slug.Make(a.Title), find.Slug)
	s.Equal(a.Title, find.Title)
	s.Equal(a.Description, find.Description)
	s.Equal(a.Body, find.Body)
	s.Equal(a.Author.ID, find.AuthorID)
	s.WithinDuration(now, find.CreatedAt, time.Minute)
	s.WithinDuration(now, find.CreatedAt, time.Minute)
	s.False(find.DeletedAt.Valid)
	var tagCount int64
	s.NoError(s.originDB.Model(new(model.Tag)).Count(&tagCount).Error)
	s.EqualValues(2, tagCount)
}

func (s *Suite) TestSaveFail() {
	exist := newArticle("article1", "", "", *s.u1, nil)
	s.NoError(s.db.Save(context.TODO(), exist))

	cases := []struct {
		name    string
		article *model.Article
		msg     string
	}{
		{
			name:    "duplicate",
			article: newArticle(exist.Title, exist.Description, exist.Body, exist.Author, nil),
			msg:     database.ErrKeyConflict.Error(),
		},
	}

	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			err := s.db.Save(context.TODO(), tc.article)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.msg)
		})
	}
}

func (s *Suite) TestUpdate() {
	exist := newArticle("article1", "description", "body", *s.u1, []string{"tag1", "tag2"})
	s.NoError(s.db.Save(context.TODO(), exist))

	update := &model.Article{
		ID:          exist.ID,
		Title:       "newslug",
		Description: "updated description",
		Body:        "updated body",
		Author:      exist.Author,
		Tags:        exist.Tags,
	}

	// when
	err := s.db.Update(context.TODO(), &exist.Author, update)
	// then
	s.NoError(err)
	var find model.Article
	s.NoError(s.originDB.First(&find, "article_id = ?", update.ID).Error)
	s.Equal(update.Slug, find.Slug)
	s.Equal(slug.Make(update.Slug), find.Slug)
	s.Equal(update.Title, find.Title)
	s.Equal(update.Description, find.Description)
	s.Equal(update.Body, find.Body)
}

func (s *Suite) TestUpdateFail() {
	articles := []*model.Article{
		newArticle("article1", "description", "body", *s.u1, nil),
		newArticle("article2", "description", "body", *s.u1, nil),
	}
	s.NoError(s.originDB.Create(&articles).Error)

	cases := []struct {
		name   string
		user   *userModel.User
		update *model.Article
		msg    string
	}{
		{
			name: "duplicate slug",
			user: s.u1,
			update: &model.Article{
				ID:    articles[0].ID,
				Title: articles[1].Title,
			},
			msg: database.ErrKeyConflict.Error(),
		}, {
			name: "not found by author",
			user: s.u2,
			update: &model.Article{
				ID:    articles[0].ID,
				Title: articles[0].Title,
			},
			msg: database.ErrRecordNotFound.Error(),
		}, {
			name: "not found by article id",
			user: s.u1,
			update: &model.Article{
				ID:   articles[1].ID + 10,
				Slug: articles[0].Slug,
			},
			msg: database.ErrRecordNotFound.Error(),
		},
	}

	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			err := s.db.Update(context.TODO(), tc.user, tc.update)
			// then
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.msg)
		})
	}
}

func (s *Suite) TestDeleteBySlug() {
	exist := newArticle("article1", "description", "body", *s.u1, []string{"tag1", "tag2"})
	s.NoError(s.db.Save(context.TODO(), exist))
	now := time.Now()
	// when
	err := s.db.DeleteBySlug(context.TODO(), &exist.Author, exist.Slug)
	// then
	s.NoError(err)
	var find model.Article
	s.NoError(s.originDB.Unscoped().First(&find, "article_id = ?", exist.ID).Error)
	s.WithinDuration(now, find.DeletedAt.Time, time.Minute)
}

func (s *Suite) TestDeleteBySlugFail() {
	exist := newArticle("article1", "description", "body", *s.u1, []string{"tag1", "tag2"})
	s.NoError(s.db.Save(context.TODO(), exist))

	cases := []struct {
		name string
		user *userModel.User
		slug string
		msg  string
	}{
		{
			name: "not found by slug",
			user: s.u1,
			slug: "not_exist_slug",
			msg:  database.ErrRecordNotFound.Error(),
		}, {
			name: "not found by author",
			user: s.u2,
			slug: exist.Slug,
			msg:  database.ErrRecordNotFound.Error(),
		},
	}

	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			err := s.db.DeleteBySlug(context.TODO(), tc.user, tc.slug)
			// then
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.msg)
		})
	}
}

func (s *Suite) TestFavoriteArticle() {
	// TODO
}

func (s *Suite) TestFavoriteArticle_Fail() {
	// TODO
}

func (s *Suite) TestUnFavoriteArticle() {
	// TODO
}

func (s *Suite) TestUnFavoriteArticle_Fail() {
	// TODO
}

func (s *Suite) TestFindTags() {
	// TODO
}

func newArticle(title, description, body string, author userModel.User, tagValues []string) *model.Article {
	var tags []*model.Tag
	for _, value := range tagValues {
		tags = append(tags, &model.Tag{Name: value})
	}

	return &model.Article{
		Title:       title,
		Description: description,
		Body:        body,
		Author:      author,
		Tags:        tags,
	}
}
