package database

import (
	"context"
	"database/sql"
	model2 "github.com/zacscoding/echo-gorm-realworld-app/internal/article/model"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/config"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/database"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/internal/user/model"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/logging"
	"gorm.io/gorm"
)

//go:generate mockery --name ArticleDB --filename article_mock.go
type ArticleDB interface {
	ArticleQueryDB
	CommentDB

	// Save saves a given article a and saves tags in article a.
	// database.ErrKeyConflict will return if duplicate emails.
	Save(ctx context.Context, a *model2.Article) error

	// Update updates a given model.Article from articleID and authorID.
	// title, description, body will be updated.
	// database.ErrRecordNotFound will be returned if not exists.
	// database.ErrKeyConflict will be returned if duplicate slug
	Update(ctx context.Context, user *userModel.User, a *model2.Article) error

	// DeleteBySlug deletes an article matched by user's id and slug.
	// database.ErrRecordNotFound will be returned if zero row affected.
	DeleteBySlug(ctx context.Context, user *userModel.User, slug string) error

	// FavoriteArticle updates the relation of article and favorites.
	// database.ErrKeyConflict will be returned if article not exists or already favorited.
	FavoriteArticle(ctx context.Context, user *userModel.User, articleID uint) error

	// UnFavoriteArticle deletes the relation of article and favorite.
	// database.ErrRecordNotFound will be returned if zero row affected.
	UnFavoriteArticle(ctx context.Context, user *userModel.User, articleID uint) error

	// FindTags returns tags all
	FindTags(ctx context.Context) ([]*model2.Tag, error)
}

type ArticleQueryDB interface {
	// FindBySlug returns a model.Article with Author, Tags, FavoritesCount if exists.
	// Favorited field will be setted if provide user.
	// database.ErrRecordNotFound will be returned if not exists.
	FindBySlug(ctx context.Context, user *userModel.User, slug string) (*model2.Article, error)

	// FindArticlesByQuery returns ([]*model.Articles, total count, error) from given queries.
	// each articles contains Author, Tags, FavoritesCount and Favorited(if provide user).
	FindArticlesByQuery(ctx context.Context, user *userModel.User, query model2.ArticleQuery, offset, limit int) (*model2.Articles, error)

	// FindArticlesByAuthors returns ([]*model.Articles, total count, error) from given author ids.
	// each articles contains Author, Tags, FavoritesCount and Favorited.
	FindArticlesByAuthors(ctx context.Context, user *userModel.User, authors []uint, offset, limit int) (*model2.Articles, error)
}

type CommentDB interface {
	// SaveComment saves a given comment c.
	// database.ErrFKConstraint will be returned if not exist article id or author id.
	SaveComment(ctx context.Context, c *model2.Comment) error

	// FindCommentsByArticleID returns ([]*model.Comments, error) from given article id.
	FindCommentsByArticleID(ctx context.Context, articleID uint) ([]*model2.Comment, error)

	// DeleteCommentByID deletes a comment matched by user'id and comment id.
	// database.ErrRecordNotFound will be returned if zero row affected.
	DeleteCommentByID(ctx context.Context, user *userModel.User, articleID, commentID uint) error
}

// NewArticleDB creates a new ArticleDB with given gorm.DB
func NewArticleDB(_ *config.Config, db *gorm.DB) ArticleDB {
	return &articleDB{
		db: db,
	}
}

type articleDB struct {
	db *gorm.DB
}

func (adb *articleDB) Save(ctx context.Context, a *model2.Article) error {
	logger := logging.FromContext(ctx)
	logger.Debugw("ArticleDB_Save try to save an article", "article", a)

	var (
		opts = &sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		}
	)
	if err := database.RunInTx(ctx, adb.db, opts, func(txDb *gorm.DB) error {
		// find tags and creates if not exists.
		for _, tag := range a.Tags {
			if err := txDb.WithContext(ctx).FirstOrCreate(tag, "name = ?", tag.Name).Error; err != nil {
				return err
			}
		}
		return txDb.WithContext(ctx).Create(a).Error
	}); err != nil {
		logger.Errorw("ArticleDB_Save failed to save an article", "err", err)
		return database.WrapError(err)
	}
	return nil
}

func (adb *articleDB) Update(ctx context.Context, user *userModel.User, a *model2.Article) error {
	logger := logging.FromContext(ctx)
	logger.Debugw("ArticleDB_Update try to update an article", "article", a)

	result := adb.db.WithContext(ctx).
		Model(a).
		Select("Slug", "Title", "Description", "Body").
		Where("article_id = ? AND author_id = ?", a.ID, user.ID).
		Updates(a)
	if result.Error != nil {
		logger.Errorw("ArticleDB_Update failed to update an article", "err", result.Error)
		return database.WrapError(result.Error)
	}
	if result.RowsAffected != 1 {
		logger.Error("ArticleDB_Update failed to update an article. zero rows affected")
		return database.WrapError(gorm.ErrRecordNotFound)
	}
	return nil
}

func (adb *articleDB) DeleteBySlug(ctx context.Context, user *userModel.User, slug string) error {
	logger := logging.FromContext(ctx)
	if user == nil {
		logger.Error("ArticleDB_DeleteBySlug no user")
		return database.WrapError(gorm.ErrRecordNotFound)
	}
	logger.Debugw("ArticleDB_DeleteBySlug try to delete an article", "userID", user.ID, "slug", slug)

	result := adb.db.WithContext(ctx).Where("slug = ? AND author_id = ?", slug, user.ID).Delete(&model2.Article{})
	if result.Error != nil {
		logger.Errorw("ArticleDB_DeleteBySlug failed to delete", "err", result.Error)
		return database.WrapError(result.Error)
	}
	if result.RowsAffected != 1 {
		logger.Error("UserDB_UnFollow failed to delete the article. zero rows affected")
		return database.WrapError(gorm.ErrRecordNotFound)
	}
	return nil
}

func (adb *articleDB) FavoriteArticle(ctx context.Context, user *userModel.User, articleID uint) error {
	logger := logging.FromContext(ctx)
	if user == nil {
		logger.Error("ArticleDB_FavoriteArticle no user")
		return database.WrapError(gorm.ErrRecordNotFound)
	}
	logger.Debugw("ArticleDB_FavoriteArticle try to favorite an article", "userID", user.ID, "articleID", articleID)

	af := model2.ArticleFavorite{
		UserID:    user.ID,
		ArticleID: articleID,
	}
	if err := adb.db.WithContext(ctx).Create(&af).Error; err != nil {
		logger.Errorw("ArticleDB_FavoriteArticle failed to favorite an article", "userID", user.ID, "articleID", articleID, "err", err)
		return database.WrapError(err)
	}
	return nil
}

func (adb *articleDB) UnFavoriteArticle(ctx context.Context, user *userModel.User, articleID uint) error {
	logger := logging.FromContext(ctx)
	if user == nil {
		logger.Error("ArticleDB_UnFavoriteArticle no user")
		return database.WrapError(gorm.ErrRecordNotFound)
	}
	logger = logger.With("userID", user.ID, "articleID", articleID)
	logger.Debug("ArticleDB_UnFavoriteArticle try to unfavorite an article")

	result := adb.db.WithContext(ctx).Unscoped().Where("article_id = ? AND user_id = ?", articleID, user.ID).Delete(new(model2.ArticleFavorite))
	if result.Error != nil {
		logger.Errorw("ArticleDB_UnFavoriteArticle failed to unfavorite", "err", result.Error)
		return database.WrapError(result.Error)
	}
	if result.RowsAffected != 1 {
		logger.Error("ArticleDB_UnFavoriteArticle failed to unfavorite the article. zero rows affected")
		return database.WrapError(gorm.ErrRecordNotFound)
	}
	return nil
}

func (adb *articleDB) FindTags(ctx context.Context) ([]*model2.Tag, error) {
	logger := logging.FromContext(ctx)
	logger.Debug("ArticleDB_FindTags try to find tags all")

	var tags []*model2.Tag
	if err := adb.db.WithContext(ctx).Find(&tags).Error; err != nil {
		logger.Errorw("ArticleDB_FindTags failed to find tags all", "err", err)
		return nil, database.WrapError(err)
	}
	return tags, nil
}
