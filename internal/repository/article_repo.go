package repository

import (
	"context"
	"github.com/soulstalker/pantrygo2/internal/entity"
)

// ArticleRepository определяет методы для работы со статьями
type ArticleRepository interface {
	Create(ctx context.Context, article *entity.Article) error
	GetByID(ctx context.Context, id uint) (*entity.Article, error)
	Update(ctx context.Context, article *entity.Article) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]*entity.Article, error)
	Search(ctx context.Context, query string, tags []string, status entity.ArticleStatus) ([]*entity.Article, error)
	ArchiveArticle(ctx context.Context, id uint) error
}

// TagRepository определяет методы для работы с тегами
type TagRepository interface {
	Create(ctx context.Context, tag *entity.Tag) error
	GetByID(ctx context.Context, id uint) (*entity.Tag, error)
	GetByName(ctx context.Context, name string) (*entity.Tag, error)
	List(ctx context.Context) ([]*entity.Tag, error)
	AddTagToArticle(ctx context.Context, articleID, tagID uint) error
	RemoveTagFromArticle(ctx context.Context, articleID, tagID uint) error
	GetTagsByArticleID(ctx context.Context, articleID uint) ([]*entity.Tag, error)
}

// MediaRepository определяет методы для работы с медиафайлами
type MediaRepository interface {
	Create(ctx context.Context, media *entity.MediaFile) error
	GetByID(ctx context.Context, id uint) (*entity.MediaFile, error)
	Delete(ctx context.Context, id uint) error
	GetByArticleID(ctx context.Context, articleID uint) ([]*entity.MediaFile, error)
}
