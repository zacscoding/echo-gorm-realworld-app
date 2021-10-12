package article

import (
	"github.com/labstack/echo/v4"
	articleDB "github.com/zacscoding/echo-gorm-realworld-app/internal/article/database"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/config"
	"github.com/zacscoding/echo-gorm-realworld-app/internal/serverenv"
	userDB "github.com/zacscoding/echo-gorm-realworld-app/internal/user/database"
)

type Handler struct {
	cfg       *config.Config
	articleDB articleDB.ArticleDB
	userDB    userDB.UserDB
}

// NewHandler returns a new Handle from given serverenv.ServerEnv and config.Config.
func NewHandler(env *serverenv.ServerEnv, cfg *config.Config) (*Handler, error) {
	return &Handler{
		cfg:       cfg,
		articleDB: articleDB.NewArticleDB(cfg, env.GetDB()),
		userDB:    userDB.NewUserDB(cfg, env.GetDB()),
	}, nil
}

// Route configures route given "/api" echo.Group to "/api/users/**, /api/profile/**" paths.
func (h *Handler) Route(e *echo.Group, authMiddleware echo.MiddlewareFunc) {
	// articles
	articleGroup := e.Group("/articles")
	articleGroup.Use(authMiddleware)
	articleGroup.GET("", h.handleGetArticles)
	articleGroup.GET("/feed", h.handleGetFeeds)
	articleGroup.GET("/:slug", h.handleGetArticle)
	articleGroup.POST("", h.handleCreateArticle)
	articleGroup.PUT("/:slug", h.handleUpdateArticle)
	articleGroup.DELETE("/:slug", h.handleDeleteArticle)
	articleGroup.POST("/:slug/favorite", h.handleFavorite)
	articleGroup.DELETE("/:slug/favorite", h.handleUnFavorite)

	// comments
	commentGroup := e.Group("/articles/:slug/comments")
	commentGroup.Use(authMiddleware)
	commentGroup.GET("", h.handleGetComments)
	commentGroup.POST("", h.handleCreateComment)
	commentGroup.DELETE("/:id", h.handleDeleteComment)

	// tags
	e.GET("/tags", h.handleGetTags)
}
