package types

import (
	articlemodel "github.com/zacscoding/echo-gorm-realworld-app/article/model"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/user/model"
	"time"
)

// ArticleResponse represents a single article response.
type ArticleResponse struct {
	Article *Article `json:"article"`
}

// ToArticleResponse converts given a to ArticleResponse.
func ToArticleResponse(a *articlemodel.Article) *ArticleResponse {
	return &ArticleResponse{
		Article: toArticle(a),
	}
}

// ArticlesResponse represents multiple articles response.
type ArticlesResponse struct {
	Articles      []*Article `json:"articles"`
	ArticlesCount int64      `json:"articlesCount"`
}

// ToArticlesResponse converts given as to ArticlesResponse.
func ToArticlesResponse(as *articlemodel.Articles) *ArticlesResponse {
	res := new(ArticlesResponse)
	res.Articles = make([]*Article, len(as.Articles))
	for i, a := range as.Articles {
		res.Articles[i] = toArticle(a)
	}
	res.ArticlesCount = as.ArticlesCount
	return res
}

type Article struct {
	Slug           string    `json:"slug"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	Tags           []string  `json:"tagList"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int       `json:"favoritesCount"`
	Author         Author    `json:"author"`
}

type Author struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

func toAuthor(u *userModel.User) Author {
	return Author{
		Username:  u.Name,
		Bio:       u.Bio,
		Image:     u.Image,
		Following: u.Following,
	}
}

func toArticle(a *articlemodel.Article) *Article {
	tags := make([]string, len(a.Tags))
	for i := 0; i < len(a.Tags); i++ {
		tags[i] = a.Tags[i].Name
	}
	return &Article{
		Slug:           a.Slug,
		Title:          a.Title,
		Description:    a.Description,
		Body:           a.Body,
		Tags:           tags,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
		Favorited:      a.Favorited,
		FavoritesCount: a.FavoritesCount,
		Author:         toAuthor(&a.Author),
	}
}
