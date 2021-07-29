package article

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/zacscoding/echo-gorm-realworld-app/api/types"
	articlemodel "github.com/zacscoding/echo-gorm-realworld-app/article/model"
	"github.com/zacscoding/echo-gorm-realworld-app/database"
	"github.com/zacscoding/echo-gorm-realworld-app/logging"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/user/model"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/authutils"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/httputils"
	"net/http"
)

// handleGetArticles handles "GET /api/articles?tag=&author=&favorited=&limit=&size=" to get articles.
func (h *Handler) handleGetArticles(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		logger      = logging.FromContext(ctx)
		query       = ArticleQuery{}
		currentUser = h.currentUser(c)
	)

	// Bind request
	if err := query.Bind(c); err != nil {
		logger.Errorw("ArticleHandler_handleGetArticles failed to bind query", "err", err)
		return httputils.WrapBindError(err)
	}

	// Query articles
	articles, err := h.articleDB.FindArticlesByQuery(ctx, currentUser, articlemodel.ArticleQuery{
		Tag:         query.Tag,
		Author:      query.Author,
		FavoritedBy: query.Favorited,
	}, query.PageableQuery.Offset, query.PageableQuery.Limit)
	if err != nil {
		return httputils.NewInternalServerError(err)
	}

	// Check follow or not given article's authors.
	if currentUser != nil {
		if err := h.checkFollowAuthorsArticles(ctx, currentUser, articles.Articles...); err != nil {
			return httputils.NewInternalServerError(err)
		}
	}
	return c.JSON(http.StatusOK, types.ToArticlesResponse(articles))
}

// handleGetFeeds handles "GET /api/articles/feed" to get feeds.
func (h *Handler) handleGetFeeds(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		logger      = logging.FromContext(ctx)
		query       = PageableQuery{}
		currentUser = h.currentUser(c)
	)

	// Bind request
	if err := query.Bind(c); err != nil {
		logger.Errorw("ArticleHandler_handleGetFeeds failed to bind query", "err", err)
		return httputils.WrapBindError(err)
	}

	// Find followers
	followers, err := h.userDB.FindFollowerIDs(ctx, currentUser.ID)
	if err != nil {
		return httputils.NewInternalServerError(err)
	}
	if len(followers) == 0 {
		return c.JSON(http.StatusOK, types.ToArticlesResponse(articlemodel.EmptyArticles))
	}

	// Find feeds
	feeds, err := h.articleDB.FindArticlesByAuthors(ctx, currentUser, followers, query.Offset, query.Limit)
	if err != nil {
		return httputils.NewInternalServerError(err)
	}
	for _, a := range feeds.Articles {
		a.Author.Following = true
	}
	return c.JSON(http.StatusOK, types.ToArticlesResponse(feeds))
}

// handleGetArticle handles "GET /api/articles/:slug" to get an article.
func (h *Handler) handleGetArticle(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		slug        = c.Param("slug")
		currentUser = h.currentUser(c)
	)

	// Query article
	article, err := h.getArticleBySlug(ctx, currentUser, slug)
	if err != nil {
		return err
	}

	// Check follow or not given article's authors.
	if currentUser != nil {
		if err := h.checkFollowAuthorsArticles(ctx, currentUser, article); err != nil {
			return httputils.NewInternalServerError(err)
		}
	}
	return c.JSON(http.StatusOK, types.ToArticleResponse(article))
}

// handleCreateArticle handles "POST /api/articles" to post an article.
func (h *Handler) handleCreateArticle(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		logger      = logging.FromContext(ctx)
		req         = &CreateArticleRequest{}
		currentUser = h.currentUser(c)
		a           articlemodel.Article
	)

	// Bind request
	if err := req.Bind(c, &a, currentUser); err != nil {
		logger.Errorw("ArticleHandler_handleCreateArticle failed to bind creating an article", "err", err)
		return httputils.WrapBindError(err)
	}

	// Save an article
	if err := h.articleDB.Save(ctx, &a); err != nil {
		if err == database.ErrKeyConflict {
			return httputils.NewStatusUnprocessableEntity("duplicate title")
		}
		return httputils.NewInternalServerError(err)
	}
	return c.JSON(http.StatusOK, types.ToArticleResponse(&a))
}

// handleUpdateArticle handles "PUT /api/articles/:slug" to update an article.
func (h *Handler) handleUpdateArticle(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		logger      = logging.FromContext(ctx)
		currentUser = h.currentUser(c)
		slug        = c.Param("slug")
	)

	// Query article
	a, err := h.getArticleBySlug(ctx, currentUser, slug)
	if err != nil {
		return err
	}

	// Bind request
	req := UpdateArticleRequest{}
	if err := req.Bind(c, a, currentUser); err != nil {
		logger.Errorw("ArticleHandler_handleUpdateArticle failed to bind updating an article", "err", err)
		return httputils.WrapBindError(err)
	}

	// Update article
	if err := h.articleDB.Update(ctx, currentUser, a); err != nil {
		if err == database.ErrKeyConflict {
			return httputils.NewStatusUnprocessableEntity("duplicate title")
		}
		return httputils.NewInternalServerError(err)
	}
	return c.JSON(http.StatusOK, types.ToArticleResponse(a))
}

// handleDeleteArticle handles "DELETE /api/articles/:slug" to delete an article.
func (h *Handler) handleDeleteArticle(c echo.Context) error {
	var (
		ctx         = c.Request().Context()
		currentUser = h.currentUser(c)
		slug        = c.Param("slug")
	)
	// Delete article
	if err := h.articleDB.DeleteBySlug(ctx, currentUser, slug); err != nil {
		if err == database.ErrRecordNotFound {
			return httputils.NewNotFoundError(fmt.Sprintf("article(%s) not found", slug))
		}
		return httputils.NewInternalServerError(err)
	}
	return c.String(http.StatusOK, "")
}

// getArticleBySlug returns an article if exists, otherwise wrapped http error
func (h *Handler) getArticleBySlug(ctx context.Context, currentUser *userModel.User, slug string) (*articlemodel.Article, error) {
	article, err := h.articleDB.FindBySlug(ctx, currentUser, slug)
	if err != nil {
		if err == database.ErrRecordNotFound {
			return nil, httputils.NewNotFoundError(fmt.Sprintf("article(%s) not found", slug))
		}
		return nil, httputils.NewInternalServerError(err)
	}
	return article, nil
}

func (h *Handler) checkFollowAuthorsArticles(ctx context.Context, u *userModel.User, articles ...*articlemodel.Article) error {
	var authors []uint
	for _, a := range articles {
		authors = append(authors, a.AuthorID)
	}

	fm, err := h.userDB.IsFollows(ctx, u.ID, authors)
	if err != nil {
		return err
	}

	for _, a := range articles {
		if follow, ok := fm[a.AuthorID]; ok && follow {
			a.Author.Following = true
		}
	}
	return nil
}

func (h *Handler) currentUser(c echo.Context) *userModel.User {
	uid := authutils.CurrentUser(c)
	if uid == 0 {
		return nil
	}
	return &userModel.User{
		ID: uid,
	}
}
