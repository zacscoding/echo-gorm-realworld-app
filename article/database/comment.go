package database

import (
	"context"
	"errors"
	"github.com/zacscoding/echo-gorm-realworld-app/article/model"
	"github.com/zacscoding/echo-gorm-realworld-app/database"
	"github.com/zacscoding/echo-gorm-realworld-app/logging"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/user/model"
	"gorm.io/gorm"
)

func (adb *articleDB) SaveComment(ctx context.Context, c *model.Comment) error {
	logger := logging.FromContext(ctx)
	logger.Debugw("CommentDB_SaveComment try to save a comment", "c", c)

	if c.ArticleID == 0 || c.AuthorID == 0 {
		logger.Errorw("CommentDB_SaveComment failed to save a comment because empty article id or author id",
			"articleID", c.ArticleID, "authorID", c.AuthorID)
		return errors.New("require article id and author id")
	}

	if err := adb.db.WithContext(ctx).Create(c).Error; err != nil {
		logger.Errorw("CommentDB_SaveComment failed to save a comment", "c", c, "err", err)
		return database.WrapError(err)
	}
	return nil
}

func (adb *articleDB) FindCommentsByArticleID(ctx context.Context, articleID uint) ([]*model.Comment, error) {
	logger := logging.FromContext(ctx)
	logger.Debugw("CommentDB_FindCommentsByArticleID try to find comments by article id", "articleID", articleID)

	var comments []*model.Comment
	if err := adb.db.WithContext(ctx).Model(new(model.Comment)).
		Joins("Author").
		Where("article_id = ?", articleID).
		Order("created_at DESC").
		Find(&comments).Error; err != nil {
		logger.Errorw("CommentDB_FindCommentsByArticleID failed to find comments", "articleID", articleID, "err", err)
		return nil, database.WrapError(err)
	}
	return comments, nil
}

func (adb *articleDB) DeleteCommentByID(ctx context.Context, user *userModel.User, commentID uint) error {
	logger := logging.FromContext(ctx)
	if user == nil {
		logger.Error("CommentDB_DeleteCommentByID no user provided")
		return database.WrapError(gorm.ErrRecordNotFound)
	}
	logger.Debugw("CommentDB_DeleteCommentByID try to delete a comment", "userID", user.ID, "commentID", commentID)

	result := adb.db.WithContext(ctx).Where("comment_id = ? AND author_id = ?", commentID, user.ID).Delete(&model.Comment{})
	if result.Error != nil {
		logger.Errorw("CommentDB_DeleteCommentByID failed to delete", "err", result.Error)
		return database.WrapError(result.Error)
	}
	if result.RowsAffected != 1 {
		logger.Error("CommentDB_DeleteCommentByID failed to delete the article. zero rows affected")
		return database.WrapError(gorm.ErrRecordNotFound)
	}
	return nil
}
