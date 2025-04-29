package http

import (
	"github.com/gin-gonic/gin"
	"github.com/soulstalker/pantrygo2/internal/entity"
	"github.com/soulstalker/pantrygo2/internal/usecase"
	"net/http"
)

// ArticleHandler обрабатывает HTTP-запросы для работы со статьями
type ArticleHandler struct {
	articleUseCase *usecase.ArticleUseCase
}

// NewArticleHandler создает новый ArticleHandler
func NewArticleHandler(articleUseCase *usecase.ArticleUseCase) *ArticleHandler {
	return &ArticleHandler{
		articleUseCase: articleUseCase,
	}
}

// RegisterRoutes регистрирует маршруты API
func (h *ArticleHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/api/group")
	{
		group.POST("/", h.CreateArticle)
		group.GET("/:id", h.GetArticle)
		group.PUT("/:id", h.UpdateArticle)
		group.DELETE("/:id", h.ArchiveArticle)
		group.GET("/", h.ListArticles)
		group.GET("/search", h.SearchArticles)
	}
}

// Структура для запроса на создание статьи
type articleRequest struct {
	Title   string   `json:"title" binding:"required"`
	Content string   `json:"content" binding:"required"`
	Status  string   `json:"status"`
	Tags    []string `json:"tags"`
}

// CreateArticle обрабатывает создание статьи
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req articleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article := &entity.Article{
		Title:       req.Title,
		Content:     req.Content,
		Status:      entity.ArticleStatus(req.Status),
		CreatedByID: 1, // Мок, потом будет тот кто залогинен
	}

	if err := h.articleUseCase.Create(c.Request.Context(), article, req.Tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": article.ID})
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
	// get
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	// update
}

func (h *ArticleHandler) ArchiveArticle(c *gin.Context) {
	// delete
}

func (h *ArticleHandler) ListArticles(c *gin.Context) {
	// delete

}
func (h *ArticleHandler) SearchArticles(c *gin.Context) {
	// delete
}
