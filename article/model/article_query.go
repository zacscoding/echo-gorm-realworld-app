package model

// ArticleQuery is used for quering recent articles.
type ArticleQuery struct {
	Tag         string
	Author      string
	FavoritedBy string
}
