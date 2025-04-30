package http

import (
	"github.com/gin-gonic/gin"
	"github.com/soulstalker/pantrygo2/internal/entity"
	"github.com/soulstalker/pantrygo2/internal/usecase"
	"net/http"
	"strconv"
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

// GetArticle обрабатывает получение статьи по ID
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	article, tags, media, err := h.articleUseCase.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Преобразуем теги в список строк для ответа
	tagNames := make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}

	c.JSON(http.StatusOK, gin.H{
		"article": article,
		"tags":    tags,
		"media":   media,
	})
}

// UpdateArticle обрабатывает обновление статьи
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req articleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article := &entity.Article{
		ID:      uint(id),
		Title:   req.Title,
		Content: req.Content,
		Status:  entity.ArticleStatus(req.Status),
	}

	if err := h.articleUseCase.Update(c.Request.Context(), article, req.Tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// ArchiveArticle обрабатывает архивирование статьи
func (h *ArticleHandler) ArchiveArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.articleUseCase.ArchiveArticle(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "archived"})
}

// ListArticles обрабатывает получение списка статей
func (h *ArticleHandler) ListArticles(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	articles, err := h.articleUseCase.List(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": articles})
}

// SearchArticles обрабатывает поиск статей
func (h *ArticleHandler) SearchArticles(c *gin.Context) {
	query := c.Query("q")
	tags := c.QueryArray("tag")
	status := entity.ArticleStatus(c.Query("status"))

	articles, err := h.articleUseCase.Search(c.Request.Context(), query, tags, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": articles})
}
