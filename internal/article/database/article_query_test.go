package database

import (
	"context"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	model2 "github.com/zacscoding/echo-gorm-realworld-app/internal/article/model"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/config"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/database"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/internal/user/model"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/logging"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm"
	"testing"
	"time"
)

type QuerySuite struct {
	suite.Suite
	db         ArticleQueryDB
	originDB   *gorm.DB
	dbTeardown database.CloseFunc
	users      []*userModel.User
}

func TestQuerySuite(t *testing.T) {
	suite.Run(t, new(QuerySuite))
}

func (s *QuerySuite) SetupSuite() {
	cfg, _ := config.Load("")
	logging.SetConfig(&logging.Config{
		Encoding:    "console",
		Level:       zapcore.FatalLevel,
		Development: false,
	})
	s.originDB, s.dbTeardown = database.NewTestDatabase(s.T(), true)
	s.db = NewArticleDB(cfg, s.originDB)
	sqlDB, err := s.originDB.DB()
	s.NoError(err)

	// Setup fixtures
	fixtures, err := testfixtures.New(
		testfixtures.Database(sqlDB),
		testfixtures.DangerousSkipTestDatabaseCheck(), // disable to check contains "test" in database
		testfixtures.Dialect("mysql"),
		testfixtures.Files(
			"fixtures/users.yaml",
			"fixtures/tags.yaml",
			"fixtures/articles.yaml",
			"fixtures/article_tags.yaml",
			"fixtures/article_favorites.yaml",
		),
	)
	s.NoError(err)
	fixtures.Load() // TODO: assert no error. currently error occur.
	var users []*userModel.User
	s.NoError(s.originDB.Find(&users).Error)
	s.users = users
}

func (s *QuerySuite) TearDownSuite() {
	s.dbTeardown()
}

func (s *QuerySuite) TestFindBySlug() {
	cases := []struct {
		name string
		slug string
		u    *userModel.User
		// expected
		exist          bool
		tagValues      string
		favorited      bool
		favoritesCount int
	}{
		{
			name:           "exist with favorited",
			slug:           "user1article1",
			u:              s.users[1], //user2
			exist:          true,
			tagValues:      "tag1 tag2",
			favorited:      true,
			favoritesCount: 2,
		}, {
			name:           "exist with no favorited",
			slug:           "user1article2",
			u:              s.users[2], //user3
			exist:          true,
			tagValues:      "tag1 tag2",
			favorited:      false,
			favoritesCount: 1,
		}, {
			name:           "exist without user",
			slug:           "user1article1",
			exist:          true,
			tagValues:      "tag1 tag2",
			favorited:      false,
			favoritesCount: 2,
		}, {
			name:  "no record",
			slug:  "user1article100",
			exist: false,
		},
	}
	ctx := context.TODO()
	for _, tc := range cases {
		s.T().Run(tc.name, func(t *testing.T) {
			a, err := s.db.FindBySlug(ctx, tc.u, tc.slug)
			if !tc.exist {
				assert.Nil(t, a)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), database.ErrRecordNotFound.Error())
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.slug, a.Slug)
			assertArticleTableValues(t, s.originDB, a)
			assertArticleExtraFields(t, a, tc.tagValues, tc.favorited, tc.favoritesCount)
		})
	}
}

func (s *QuerySuite) TestFindArticlesByQuery() {
	var (
		user  = s.users[2]
		query = model2.ArticleQuery{
			Tag:         "tag1",
			Author:      "user1",
			FavoritedBy: "user2",
		}
		limit = 2
	)
	// first iteration
	articles, err := s.db.FindArticlesByQuery(context.TODO(), user, query, 0, limit)
	s.NoError(err)
	s.Len(articles.Articles, 2)
	s.EqualValues(articles.ArticlesCount, 4)
	assertArticleTableValues(s.T(), s.originDB, articles.Articles[0])
	assertArticleExtraFields(s.T(), articles.Articles[0], "tag1 tag2", true, 2)
	assertArticleTableValues(s.T(), s.originDB, articles.Articles[1])
	assertArticleExtraFields(s.T(), articles.Articles[1], "tag1 tag3", false, 1)

	// second iteration
	articles, err = s.db.FindArticlesByQuery(context.TODO(), user, query, limit, limit)
	s.NoError(err)
	s.Len(articles.Articles, 2)
	s.EqualValues(articles.ArticlesCount, 4)
	assertArticleTableValues(s.T(), s.originDB, articles.Articles[0])
	assertArticleExtraFields(s.T(), articles.Articles[0], "tag1", false, 1)
	assertArticleTableValues(s.T(), s.originDB, articles.Articles[1])
	assertArticleExtraFields(s.T(), articles.Articles[1], "tag1 tag2", true, 2)
}

func (s *QuerySuite) TestFindArticlesByAuthors() {
	var (
		user    = s.users[0]
		authors = []uint{s.users[1].ID, s.users[2].ID} // user2, user3
		limit   = 2
	)

	// first iteration
	articles, err := s.db.FindArticlesByAuthors(context.TODO(), user, authors, 0, limit)
	s.NoError(err)
	s.Len(articles.Articles, 2)
	s.EqualValues(articles.ArticlesCount, 3)
	assertArticleTableValues(s.T(), s.originDB, articles.Articles[0])
	assertArticleExtraFields(s.T(), articles.Articles[0], "tag3", true, 2)
	assertArticleTableValues(s.T(), s.originDB, articles.Articles[1])
	assertArticleExtraFields(s.T(), articles.Articles[1], "tag2", false, 1)

	// second iteration
	articles, err = s.db.FindArticlesByAuthors(context.TODO(), user, authors, limit, limit)
	s.NoError(err)
	s.Len(articles.Articles, 1)
	s.EqualValues(articles.ArticlesCount, 3)
	assertArticleTableValues(s.T(), s.originDB, articles.Articles[0])
	assertArticleExtraFields(s.T(), articles.Articles[0], "tag1 tag3", true, 2)
}

func assertArticleTableValues(t *testing.T, db *gorm.DB, a *model2.Article) {
	var find model2.Article
	assert.NoError(t, db.First(&find, "article_id = ?", a.ID).Error)
	assert.Equal(t, find.Slug, a.Slug)
	assert.Equal(t, find.Title, a.Title)
	assert.Equal(t, find.Body, a.Body)
	assert.Equal(t, find.AuthorID, a.AuthorID)
	assert.WithinDuration(t, find.CreatedAt, a.CreatedAt, time.Minute)
	assert.WithinDuration(t, find.UpdatedAt, a.UpdatedAt, time.Minute)
	assert.Equal(t, find.DeletedAt.Time, a.DeletedAt.Time)
}

func assertArticleExtraFields(t *testing.T, a *model2.Article, tagValues string, favorited bool, favoritesCount int) {
	for _, tag := range a.Tags {
		assert.Contains(t, tagValues, tag.Name)
	}
	assert.Equal(t, favorited, a.Favorited)
	assert.EqualValues(t, favoritesCount, a.FavoritesCount)
}
