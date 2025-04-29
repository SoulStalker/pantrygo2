package http

import (
	"github.com/gin-gonic/gin"
	"github.com/soulstalker/pantrygo2/internal/usecase"
)

// ArticleHandler обрабатывает HTTP-запросы для работы со статьями
type ArticleHandler struct {
	articleUseCase *usecase.ArticleUseCase
}

func (h *ArticleHandler) RegisterRoutes(r *gin.Engine) {
	articles := r.Group("/api/articles")
	articles.POST("/", h.CreateArticle)
	articles.GET("/:id", h.GetArticle)
	articles.POST("/:id", h.UpdateArticle)
	articles.DELETE("/:id", h.DeleteArticle)
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	// create
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
	// get
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	// update
}

func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	// delete
}
