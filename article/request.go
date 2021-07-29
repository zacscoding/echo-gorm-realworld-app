package article

import (
	"github.com/labstack/echo/v4"
	articlemodel "github.com/zacscoding/echo-gorm-realworld-app/article/model"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/user/model"
	"github.com/zacscoding/echo-gorm-realworld-app/utils/httputils"
)

type PageableQuery struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

func (r *PageableQuery) Bind(ctx echo.Context) error {
	if err := httputils.BindAndValidate(ctx, r); err != nil {
		return err
	}
	if err := r.Validate(); err != nil {
		return err
	}
	return nil
}

func (r *PageableQuery) Validate() error {
	if r.Limit < 0 {
		return httputils.NewStatusUnprocessableEntity("limit must greater than or equals to 0")
	}
	if r.Offset < 0 {
		return httputils.NewStatusUnprocessableEntity("offset must greater than or equals to 0")
	}
	if r.Limit == 0 {
		r.Limit = 20
	}
	return nil
}

type ArticleQuery struct {
	*PageableQuery
	Tag       string `query:"tag"`
	Author    string `query:"author"`
	Favorited string `query:"favorited"`
}

func (r *ArticleQuery) Bind(ctx echo.Context) error {
	if err := httputils.BindAndValidate(ctx, r); err != nil {
		return err
	}
	if err := r.PageableQuery.Validate(); err != nil {
		return err
	}
	return nil
}

// CreateArticleRequest represents request body data of creating an article.
type CreateArticleRequest struct {
	Article struct {
		Title       string   `json:"title" validate:"required"`
		Description string   `json:"description" validate:"required"`
		Body        string   `json:"body" validate:"required"`
		Tags        []string `json:"tagList"`
	} `json:"article"`
}

func (r *CreateArticleRequest) Bind(ctx echo.Context, a *articlemodel.Article, u *userModel.User) error {
	if err := httputils.BindAndValidate(ctx, r); err != nil {
		return err
	}
	a.Title = r.Article.Title
	a.Description = r.Article.Description
	a.Body = r.Article.Body
	for _, tag := range r.Article.Tags {
		a.Tags = append(a.Tags, &articlemodel.Tag{
			Name: tag,
		})
	}
	a.AuthorID = u.ID
	a.Author = *u
	return nil
}

type UpdateArticleRequest struct {
	Article struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Body        string `json:"body"`
	} `json:"article"`
}

func (r *UpdateArticleRequest) Bind(ctx echo.Context, a *articlemodel.Article, u *userModel.User) error {
	if err := httputils.BindAndValidate(ctx, r); err != nil {
		return err
	}
	if r.Article.Title != "" {
		a.Title = r.Article.Title
	}
	if r.Article.Description != "" {
		a.Description = r.Article.Description
	}
	if r.Article.Body != "" {
		a.Body = r.Article.Body
	}
	return nil
}
