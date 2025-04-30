package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/soulstalker/pantrygo2/internal/entity"
	"github.com/soulstalker/pantrygo2/internal/repository"
	"gorm.io/gorm"
)

// ArticleRepo реализация ArticleRepository для PostgreSQL
type ArticleRepo struct {
	db *gorm.DB
}

// NewArticleRepo создает новый ArticleRepo
func NewArticleRepo(db *gorm.DB) repository.ArticleRepository {
	return &ArticleRepo{
		db: db,
	}
}

// Create создает новую статью
func (r *ArticleRepo) Create(ctx context.Context, article *entity.Article) error {
	return r.db.WithContext(ctx).Create(article).Error
}

// GetByID возвращает статью по ID
func (r *ArticleRepo) GetByID(ctx context.Context, id uint) (*entity.Article, error) {
	var article entity.Article
	err := r.db.WithContext(ctx).First(&article, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("article not found: %w", err)
		}
		return nil, err
	}
	return &article, nil
}

// Update обновляет статью
func (r *ArticleRepo) Update(ctx context.Context, article *entity.Article) error {
	return r.db.WithContext(ctx).Save(article).Error
}

// Delete удаляет статью
func (r *ArticleRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Article{}, id).Error
}

// ArchiveArticle архивирует статью
func (r *ArticleRepo) ArchiveArticle(ctc context.Context, id uint) error {
	return nil
	// todo доделай
}

// List возвращает список статей с пагинацией
func (r *ArticleRepo) List(ctx context.Context, limit, offset int) ([]*entity.Article, error) {
	var articles []*entity.Article
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&articles).Error
	return articles, err
}

// Search ищет статьи по тексту и тегам
func (r *ArticleRepo) Search(ctx context.Context, query string, tags []string, status entity.ArticleStatus) ([]*entity.Article, error) {
	db := r.db.WithContext(ctx)

	// Базовый запрос
	tx := db.Model(&entity.Article{})

	// Фильтрация по статусу, если указан
	if status != "" {
		tx = tx.Where("status = ?", status)
	}

	// Поиск по тексту, если запрос не пустой
	if query != "" {
		tx = tx.Where("title ILIKE ? OR content ILIKE ?", "%"+query+"%", "%"+query+"%")
	}

	// Фильтрация по тегам, если указаны
	if len(tags) > 0 {
		tx = tx.Joins("JOIN article_tags ON article_tags.article_id = articles.id")
		tx = tx.Joins("JOIN tags ON tags.id = article_tags.tag_id")
		tx = tx.Where("tags.name IN ?", tags)
		tx = tx.Group("articles.id")
		tx = tx.Having("COUNT(DISTINCT tags.id) = ?", len(tags))
	}

	var articles []*entity.Article
	err := tx.Find(&articles).Error
	return articles, err
}
