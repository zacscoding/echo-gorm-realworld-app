package model

import (
	"github.com/gosimple/slug"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/internal/user/model"
	"gorm.io/gorm"
	"time"
)

const (
	TableNameArticle         = "articles"
	TableNameArticleFavorite = "article_favorites"
	TableNameTag             = "tags"
	TableNameArticleTag      = "article_tags"
	TableNameComment         = "comments"
)

var EmptyArticles = &Articles{Articles: make([]*Article, 0), ArticlesCount: 0}

// Articles represents article list with total size.
type Articles struct {
	Articles      []*Article `json:"articles"`
	ArticlesCount int64      `json:"articlesCount"`
}

// Article represents database model for articles.
type Article struct {
	ID          uint              `gorm:"column:article_id"`
	Slug        string            `gorm:"column:slug"`
	Title       string            `gorm:"column:title"`
	Description string            `gorm:"column:description"`
	Body        string            `gorm:"column:body"`
	Author      userModel.User    `json:"-"`
	AuthorID    uint              `column:"author_id"`
	Tags        []*Tag            `gorm:"many2many:article_tags;association_autocreate:false"`
	Favorites   []ArticleFavorite `gorm:"many2many:article_favorites;"`
	Comment     []Comment         `gorm:"ForeignKey:ArticleID"`

	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`

	Favorited      bool `gorm:"-"`
	FavoritesCount int  `gorm:"-"`
}

func (a *Article) TableName() string {
	return TableNameArticle
}

func (a *Article) BeforeCreate(_ *gorm.DB) error {
	a.Slug = slug.Make(a.Title)
	return nil
}

func (a *Article) BeforeUpdate(_ *gorm.DB) error {
	a.Slug = slug.Make(a.Title)
	return nil
}

// ArticleFavorite represents relation articles and favoraties.
type ArticleFavorite struct {
	User      userModel.User
	UserID    uint
	Article   Article
	ArticleID uint
}

func (af ArticleFavorite) TableName() string {
	return TableNameArticleFavorite
}

// Tag represents database model for tags.
type Tag struct {
	ID   uint   `gorm:"column:tag_id"`
	Name string `gorm:"column:name"`
	//Articles  []Article `gorm:"many2many:article_tags;"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (t Tag) TableName() string {
	return TableNameTag
}

// ArticleTag represents relation articles and tags.
type ArticleTag struct {
	Article   Article
	ArticleID uint
	Tag       Tag
	TagID     uint
}

func (at ArticleTag) TableName() string {
	return TableNameArticleTag
}

// Comment represents database model for comments.
type Comment struct {
	ID        uint   `gorm:"column:comment_id"`
	Body      string `gorm:"column:body"`
	ArticleID uint   `gorm:"column:article_id"`
	Author    userModel.User
	AuthorID  uint           `gorm:"column:author_id"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}

func (c Comment) TableName() string {
	return TableNameComment
}
