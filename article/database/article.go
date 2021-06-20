package database

import (
	"context"
	"database/sql"
	"github.com/zacscoding/echo-gorm-realworld-app/article/model"
	"github.com/zacscoding/echo-gorm-realworld-app/database"
	"github.com/zacscoding/echo-gorm-realworld-app/logging"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/user/model"
	"gorm.io/gorm"
)

//go:generate mockery --name ArticleDB --filename article_mock.go
type ArticleDB interface {
	ArticleQueryDB
	CommentDB

	// Save saves a given article a and saves tags in article a.
	// database.ErrKeyConflict will return if duplicate emails.
	Save(ctx context.Context, a *model.Article) error

	// Update updates a given model.Article from articleID and authorID.
	// title, description, body will be updated.
	// database.ErrRecordNotFound will be returned if not exists.
	// database.ErrKeyConflict will be returned if duplicate slug
	Update(ctx context.Context, user *userModel.User, a *model.Article) error

	// DeleteBySlug deletes an article matched by user's id and slug.
	// database.ErrRecordNotFound will be returned if zero row affected.
	DeleteBySlug(ctx context.Context, user *userModel.User, slug string) error
}

type ArticleQueryDB interface {
	// FindBySlug returns a model.Article with Author, Tags, FavoritesCount if exists.
	// Favorited field will be setted if provide user.
	// database.ErrRecordNotFound will be returned if not exists.
	FindBySlug(ctx context.Context, user *userModel.User, slug string) (*model.Article, error)

	// FindArticlesByQuery returns ([]*model.Articles, total count, error) from given queries.
	// each articles contains Author, Tags, FavoritesCount and Favorited(if provide user).
	FindArticlesByQuery(ctx context.Context, user *userModel.User, query model.ArticleQuery, offset, limit int) (*model.Articles, error)

	// FindArticlesByAuthors returns ([]*model.Articles, total count, error) from given author ids.
	// each articles contains Author, Tags, FavoritesCount and Favorited.
	FindArticlesByAuthors(ctx context.Context, user *userModel.User, authors []uint, offset, limit int) (*model.Articles, error)
}

type CommentDB interface {
	// SaveComment saves a given comment c.
	// database.ErrFKConstraint will be returned if not exist article id or author id.
	SaveComment(ctx context.Context, c *model.Comment) error

	// FindCommentsByArticleID returns ([]*model.Comments, error) from given article id.
	FindCommentsByArticleID(ctx context.Context, articleID uint) ([]*model.Comment, error)

	// DeleteCommentByID deletes a comment matched by user'id and comment id.
	// database.ErrRecordNotFound will be returned if zero row affected.
	DeleteCommentByID(ctx context.Context, user *userModel.User, commentID uint) error
}

// NewArticleDB creates a new ArticleDB with given gorm.DB
func NewArticleDB(db *gorm.DB) ArticleDB {
	return &articleDB{
		db: db,
	}
}

type articleDB struct {
	db *gorm.DB
}

func (adb *articleDB) Save(ctx context.Context, a *model.Article) error {
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

func (adb *articleDB) Update(ctx context.Context, user *userModel.User, a *model.Article) error {
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

	result := adb.db.WithContext(ctx).Where("slug = ? AND author_id = ?", slug, user.ID).Delete(&model.Article{})
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
