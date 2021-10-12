package database

import (
	"context"
	"errors"
	model2 "github.com/zacscoding/echo-gorm-realworld-app/internal/article/model"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/database"
	userModel "github.com/zacscoding/echo-gorm-realworld-app/internal/user/model"
	"github.com/zacscoding/echo-gorm-realworld-app/pkg/logging"
	"gorm.io/gorm"
	"time"
)

func (adb *articleDB) FindBySlug(ctx context.Context, user *userModel.User, slug string) (*model2.Article, error) {
	logger := logging.FromContext(ctx)
	logger.Debugw("ArticleDB_FindBySlug try to find an article", "user", user, "slug", slug)

	var (
		article model2.Article
		db      = adb.db.WithContext(ctx)
	)

	// find article with author
	if err := db.Joins("Author").First(&article, "slug = ?", slug).Error; err != nil {
		logger.Error("ArticleDB_FindBySlug failed to find an article", "slug", slug, "err", err)
		return nil, database.WrapError(err)
	}
	// load tags
	if err := db.Model(&article).Association("Tags").Find(&article.Tags); err != nil {
		logger.Error("ArticleDB_FindBySlug failed to fetch tags", "slug", slug, "err", err)
		return nil, database.WrapError(err)
	}
	// set favorite count
	if err := setFavoriteCount(db, &article); err != nil {
		logger.Error("ArticleDB_FindBySlug failed to fetch favorites count", "articleID", article.ID, "err", err)
		return nil, database.WrapError(err)
	}
	if user != nil {
		// set favorited from current user
		if err := setFavorited(db, user, &article); err != nil {
			logger.Error("ArticleDB_FindBySlug failed to fetch favorited", "user", user, "articleID", article.ID, "err", err)
			return nil, database.WrapError(err)
		}
	}
	return &article, nil
}

func (adb *articleDB) FindArticlesByQuery(ctx context.Context, user *userModel.User, query model2.ArticleQuery, offset, limit int) (*model2.Articles, error) {
	var (
		logger = logging.FromContext(ctx)
		userID = uint(0)
	)
	if user != nil && user.ID != 0 {
		userID = user.ID
	}
	logger.Debugw("ArticleDB_FindArticlesByQuery try to find recent articles", "userID", userID, "query", query, "offset", offset, "limit", limit)

	if limit <= 0 {
		return &model2.Articles{
			Articles:      make([]*model2.Article, 0),
			ArticlesCount: 0,
		}, nil
	}

	// find article ids from given query and offset, limit.
	ids, err := articleIdsByQuery(ctx, adb.db, query, offset, limit)
	if err != nil {
		logger.Error("ArticleDB_FindArticlesByQuery failed to fetch article ids", "userID", userID, "query", query, "offset", offset, "limit", limit, "err", err)
		return nil, database.WrapError(err)
	}

	// find total count from given query.
	total, err := countArticleByQuery(ctx, adb.db, query)
	if err != nil {
		logger.Error("ArticleDB_FindArticlesByQuery failed to fetch total count", "userID", userID, "query", query, "err", err)
		return nil, database.WrapError(err)
	}

	if len(ids) == 0 {
		return &model2.Articles{
			Articles:      make([]*model2.Article, 0),
			ArticlesCount: total,
		}, nil
	}

	articles, err := articlesByIds(adb.db, ids)
	if err != nil {
		logger.Error("ArticleDB_FindArticlesByQuery failed to fetch articles with author and tags.", "userID", userID, "ids", ids, "err", err)
		return nil, database.WrapError(err)
	}
	if err := fillArticlesExtraData(adb.db, user, articles); err != nil {
		logger.Error("ArticleDB_FindArticlesByQuery failed to update extra data", "userID", userID, "ids", ids, "err", err)
	}
	return &model2.Articles{
		Articles:      articles,
		ArticlesCount: total,
	}, nil
}

func (adb *articleDB) FindArticlesByAuthors(ctx context.Context, user *userModel.User, authors []uint, offset, limit int) (*model2.Articles, error) {
	var (
		logger = logging.FromContext(ctx)
		userID = uint(0)
	)
	if user != nil {
		userID = user.ID
	}
	logger.Debugw("ArticleDB_FindArticlesByAuthors try to find feed articles", "userID", userID, "authors", authors, "offset", offset, "limit", limit)

	if len(authors) == 0 || limit <= 0 {
		return &model2.Articles{
			Articles:      make([]*model2.Article, 0),
			ArticlesCount: 0,
		}, nil
	}

	// find article ids from given author ids and offset, limit.
	ids, err := articleIdsByAuthors(adb.db, authors, offset, limit)
	if err != nil {
		logger.Error("ArticleDB_FindArticlesByAuthors failed to fetch article ids", "authors", authors, "offset", offset, "limit", limit, "err", err)
		return nil, database.WrapError(err)
	}

	// find total count from given authors.
	total, err := countArticleByAuthors(adb.db, authors)
	if err != nil {
		logger.Error("ArticleDB_FindArticlesByAuthors failed to fetch total count", "authors", authors, "err", err)
		return nil, database.WrapError(err)
	}

	if len(ids) == 0 {
		return &model2.Articles{
			Articles:      make([]*model2.Article, 0),
			ArticlesCount: total,
		}, nil
	}

	articles, err := articlesByIds(adb.db, ids)
	if err != nil {
		logger.Error("ArticleDB_FindArticlesByAuthors failed to fetch articles with author and tags.", "userID", userID, "ids", ids, "err", err)
		return nil, database.WrapError(err)
	}
	if err := fillArticlesExtraData(adb.db, user, articles); err != nil {
		logger.Error("ArticleDB_FindArticlesByAuthors failed to update extra data", "userID", userID, "ids", ids, "err", err)
	}
	return &model2.Articles{
		Articles:      articles,
		ArticlesCount: total,
	}, nil
}

func fillArticlesExtraData(db *gorm.DB, user *userModel.User, articles []*model2.Article) error {
	// set favorites count
	if err := setFavoriteCountBulk(db, articles); err != nil {
		return err
	}
	// set is favorited from given user to articles.
	if user != nil {
		if err := setFavoritedBulk(db, user, articles); err != nil {
			return err
		}
	}
	return nil
}

// setFavoriteCount sets given article's FavoritesCount field.
func setFavoriteCount(db *gorm.DB, article *model2.Article) error {
	var favoritesCount int64
	if err := db.Model(new(model2.ArticleFavorite)).
		Where("article_id = ?", article.ID).
		Count(&favoritesCount).Error; err != nil {
		return err
	}
	article.FavoritesCount = int(favoritesCount)
	return nil
}

// setFavoriteCountBulk sets FavoritesCount field on each article.
func setFavoriteCountBulk(db *gorm.DB, articles []*model2.Article) error {
	m := make(map[uint]*model2.Article)
	ids := make([]uint, len(articles))
	for i, a := range articles {
		ids[i] = a.ID
		m[a.ID] = a
	}

	type FavoriteCount struct {
		ArticleID      uint `gorm:"column:article_id"`
		FavoritesCount int  `gorm:"column:favorites_count"`
	}

	var counts []*FavoriteCount
	if err := db.Table("(?) as g", db.Model(new(model2.ArticleFavorite)).Where("article_id IN (?)", ids)).
		Group("g.article_id").
		Select("g.article_id, count(g.user_id) as favorites_count").
		Find(&counts).Error; err != nil {
		return err
	}

	for _, c := range counts {
		if a, ok := m[c.ArticleID]; ok {
			a.FavoritesCount = c.FavoritesCount
		}
	}
	return nil
}

// setFavorited sets given article's Favorited field
func setFavorited(db *gorm.DB, user *userModel.User, article *model2.Article) error {
	if user == nil {
		return errors.New("required user")
	}
	var favorited int64
	if err := db.Model(new(model2.ArticleFavorite)).
		Where("article_id = ? AND user_id = ?", article.ID, user.ID).
		Count(&favorited).Error; err != nil {
		return err
	}
	article.Favorited = favorited == 1
	return nil
}

// setFavoritedBulk sets Favorited field from given user on each article.
func setFavoritedBulk(db *gorm.DB, user *userModel.User, articles []*model2.Article) error {
	if user == nil {
		return nil
	}
	m := make(map[uint]*model2.Article)
	ids := make([]uint, len(articles))
	for i, a := range articles {
		ids[i] = a.ID
		m[a.ID] = a
	}

	var favoritedArticleIds []uint
	if err := db.Model(new(model2.ArticleFavorite)).
		Where("article_id IN (?) AND user_id = ?", ids, user.ID).
		Select("article_id").
		Find(&favoritedArticleIds).Error; err != nil {
		return err
	}

	for _, id := range favoritedArticleIds {
		if a, ok := m[id]; ok {
			a.Favorited = true
		}
	}
	return nil
}

// articleIdsByQuery returns article ids from given query and offset, limit.
func articleIdsByQuery(ctx context.Context, db *gorm.DB, query model2.ArticleQuery, offset, limit int) ([]uint, error) {
	db = buildArticleQuery(ctx, db, query)
	rows, err := db.Select("DISTINCT a.article_id, a.created_at").
		Where("a.deleted_at IS NULL").
		Order("a.created_at DESC").
		Offset(offset).
		Limit(limit).
		Rows()
	if err != nil {
		return nil, err
	}
	var ids []uint
	for rows.Next() {
		var (
			id        uint
			createdAt time.Time
		)
		if err := rows.Scan(&id, &createdAt); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// articleIdsByAuthors returns article ids from given authors and offset, limit.
func articleIdsByAuthors(db *gorm.DB, authors []uint, offset, limit int) ([]uint, error) {
	var ids []uint
	if err := db.Model(new(model2.Article)).
		Select("article_id").
		Where("author_id IN (?) AND deleted_at IS NULL", authors).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&ids).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func countArticleByQuery(ctx context.Context, db *gorm.DB, query model2.ArticleQuery) (int64, error) {
	db = buildArticleQuery(ctx, db, query)
	var count int64
	return count, db.Distinct("a.article_id").Count(&count).Error
}

func countArticleByAuthors(db *gorm.DB, authors []uint) (int64, error) {
	var count int64
	return count, db.Model(new(model2.Article)).Where("deleted_at IS NULL AND author_id IN (?)", authors).Count(&count).Error
}

func buildArticleQuery(ctx context.Context, db *gorm.DB, query model2.ArticleQuery) *gorm.DB {
	db = db.WithContext(ctx).Table("articles a").
		Joins("LEFT JOIN article_tags at ON at.article_id = a.article_id").
		Joins("LEFT JOIN tags t ON t.tag_id = at.tag_id").
		Joins("LEFT JOIN article_favorites af ON af.article_id = a.article_id").
		Joins("LEFT JOIN users u ON u.user_id = a.author_id").
		Joins("LEFT JOIN users uf ON uf.user_id = af.user_id").
		Where("a.deleted_at IS NULL")
	if query.Tag != "" {
		db = db.Where("t.name = ?", query.Tag)
	}
	if query.Author != "" {
		db = db.Where("u.name = ?", query.Author)
	}
	if query.FavoritedBy != "" {
		db = db.Where("uf.name = ?", query.FavoritedBy)
	}
	return db
}

// articlesByIds find articles with author and tags from given article ids.
func articlesByIds(db *gorm.DB, ids []uint) ([]*model2.Article, error) {
	if len(ids) == 0 {
		return []*model2.Article{}, nil
	}

	var articles []*model2.Article
	// load articles with author (eager loading)
	if err := db.Model(new(model2.Article)).
		Joins("Author").
		Where("articles.article_id IN (?)", ids).
		Order("articles.created_at DESC").
		Find(&articles).Error; err != nil {
		return nil, err
	}
	if len(articles) == 0 {
		return articles, nil
	}

	// load tags
	m := make(map[uint]*model2.Article, len(articles))
	for _, a := range articles {
		m[a.ID] = a
	}
	type ArticleTag struct {
		model2.Tag
		ArticleId uint `gorm:"article_id"`
	}
	batchSize := 50 // will use config value.
	for i := 0; i < len(articles); i += batchSize {
		var at []ArticleTag
		last := i + batchSize
		if last > len(articles) {
			last = len(articles)
		}
		if err := db.Table("tags").
			Joins("LEFT JOIN article_tags ON article_tags.tag_id = tags.tag_id").
			Where("article_tags.article_id IN (?)", ids[i:last]).
			Select("tags.*, article_tags.article_id article_id").
			Find(&at).Error; err != nil {
			return nil, err
		}
		for _, tag := range at {
			a := m[tag.ArticleId]
			a.Tags = append(a.Tags, &model2.Tag{
				ID:        tag.ID,
				Name:      tag.Name,
				CreatedAt: tag.CreatedAt,
			})
		}
	}
	return articles, nil
}
