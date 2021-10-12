package article

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	articlemodel "github.com/zacscoding/echo-gorm-realworld-app/internal/article/model"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/database"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/internal/user/model"
	types2 "github.com/zacscoding/echo-gorm-realworld-app/pkg/api/types"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/logging"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/utils/httputils"
	"net/http"
	"strconv"
)

// handleGetComments handles "GET /articles/:slug/comments" to find comments.
func (h *Handler) handleGetComments(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		currentUser = h.currentUser(c)
		slug        = c.Param("slug")
	)

	// Query article
	article, err := h.getArticleBySlug(ctx, currentUser, slug)
	if err != nil {
		return err
	}

	// Query comments
	comments, err := h.articleDB.FindCommentsByArticleID(ctx, article.ID)
	if err != nil {
		return httputils.NewInternalServerError(err)
	}

	// Check follow or not given comment's authors.
	if currentUser != nil {
		if err := h.checkFollowAuthorsFromComments(ctx, currentUser, comments...); err != nil {
			return httputils.NewInternalServerError(err)
		}
	}
	return c.JSON(http.StatusOK, types2.ToCommentsResponse(comments))
}

// handleCreateComment handles "POST /api/articles/:slug/comments" to create a new comment.
func (h *Handler) handleCreateComment(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		logger      = logging.FromContext(ctx)
		currentUser = h.currentUser(c)
		slug        = c.Param("slug")
		req         = &CreateCommentRequest{}
		comment     articlemodel.Comment
	)

	// Query article
	article, err := h.getArticleBySlug(ctx, currentUser, slug)
	if err != nil {
		return err
	}

	// Bind request
	if err := req.Bind(c, article, &comment, currentUser); err != nil {
		logger.Errorw("ArticleHandler_handleCreateCommente failed to bind creating a comment", "err", err)
		return httputils.WrapBindError(err)
	}

	// Save a comment
	if err := h.articleDB.SaveComment(ctx, &comment); err != nil {
		return httputils.NewInternalServerError(err)
	}
	return c.JSON(http.StatusOK, types2.ToCommentResponse(&comment))
}

// handleDeleteComment handles "DELETE /articles/:slug/comments/:id" to delete a comment.
func (h *Handler) handleDeleteComment(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		logger      = logging.FromContext(ctx)
		currentUser = h.currentUser(c)
		slug        = c.Param("slug")
		commentID   = c.Param("id")
	)

	// Bind request
	cid, err := strconv.ParseUint(commentID, 10, 64)
	if err != nil {
		logger.Errorw("ArticleHandler_handleCreateComment invalid comment id", "commentID", commentID, "err", err)
		return httputils.NewBindError("id", "uint")
	}

	// Query article
	article, err := h.getArticleBySlug(ctx, currentUser, slug)
	if err != nil {
		return err
	}

	// Delete a comment
	if err := h.articleDB.DeleteCommentByID(ctx, currentUser, article.ID, uint(cid)); err != nil {
		if err == database.ErrRecordNotFound {
			return httputils.NewNotFoundError(fmt.Sprintf("comment(%d) not found", cid))
		}
		return httputils.NewInternalServerError(err)
	}
	return c.JSON(http.StatusOK, types2.ToStatusResponse(types2.StatusDeleted, nil))
}

// TODO: change author in article and comment to pointer because of duplicate codes.
func (h *Handler) checkFollowAuthorsFromComments(ctx context.Context, u *userModel.User, comments ...*articlemodel.Comment) error {
	if len(comments) == 0 {
		return nil
	}
	var authors []uint
	for _, c := range comments {
		authors = append(authors, c.AuthorID)
	}

	fm, err := h.userDB.IsFollows(ctx, u.ID, authors)
	if err != nil {
		return err
	}

	for _, c := range comments {
		if follow, ok := fm[c.AuthorID]; ok && follow {
			c.Author.Following = true
		}
	}
	return nil
}
