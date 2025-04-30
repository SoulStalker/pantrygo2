package usecase

import (
	"context"
	"errors"
	"github.com/soulstalker/pantrygo2/internal/entity"
	"github.com/soulstalker/pantrygo2/internal/repository"
	"gorm.io/gorm"
	"time"
)

// ArticleUseCase содержит бизнес-логику для работы со статьями
type ArticleUseCase struct {
	articleRepo repository.ArticleRepository
	tagRepo     repository.TagRepository
	mediaRepo   repository.MediaRepository
}

func NewArticleUseCase(
	articleRepo repository.ArticleRepository,
	tagRepo repository.TagRepository,
	mediaRepo repository.MediaRepository,
) *ArticleUseCase {
	return &ArticleUseCase{
		articleRepo: articleRepo,
		tagRepo:     tagRepo,
		mediaRepo:   mediaRepo,
	}
}

// Create создает новую статью
func (uc *ArticleUseCase) Create(ctx context.Context, article *entity.Article, tagNames []string) error {
	now := time.Now()
	article.CreatedAt = now
	article.UpdatedAt = now

	// Если статус не указан, устанавливаем черновик
	if article.Status == "" {
		article.Status = entity.StatusDraft
	}

	// Создаем статью
	if err := uc.articleRepo.Create(ctx, article); err != nil {
		return err
	}

	//Добавляем теги, если указаны
	for _, tagName := range tagNames {
		tag, err := uc.tagRepo.GetByName(ctx, tagName)
		if err != nil {
			// Если тег не найден, создаем новый
			if errors.Is(err, gorm.ErrRecordNotFound) {
				tag = &entity.Tag{
					Name:      tagName,
					CreatedAt: now,
				}
				if err := uc.tagRepo.Create(ctx, tag); err != nil {
					return err
				}
			} else {
				return err
			}
		}
		// Связываем тег со статьей
		if err := uc.tagRepo.AddTagToArticle(ctx, article.ID, tag.ID); err != nil {
			return err
		}
	}
	return nil
}

// GetByID получает статью по ID
func (uc *ArticleUseCase) GetByID(ctx context.Context, id uint) (*entity.Article, []*entity.Tag, []*entity.MediaFile, error) {
	// Получаем статью
	article, err := uc.articleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}

	// Получаем теги статьи
	tags, err := uc.tagRepo.GetTagsByArticleID(ctx, id)
	if err != nil {
		return article, nil, nil, err
	}

	// Получаем медиафайлы статьи
	media, err := uc.mediaRepo.GetByArticleID(ctx, id)
	if err != nil {
		return article, tags, nil, err
	}

	return article, tags, media, nil
}

// Update обновляет статью
func (uc *ArticleUseCase) Update(ctx context.Context, article *entity.Article, tagNames []string) error {
	// Проверяем существование статьи
	existingArticle, err := uc.articleRepo.GetByID(ctx, article.ID)
	if err != nil {
		return err
	}

	// Обновляем время изменения
	article.UpdatedAt = time.Now()
	article.CreatedAt = existingArticle.CreatedAt // Время создания не меняем

	// Получаем теги
	currentTags, err := uc.tagRepo.GetTagsByArticleID(ctx, article.ID)
	if err != nil {
		return err
	}

	// Карта текущих тегов
	currentTagMap := make(map[string]*entity.Tag)
	for _, tag := range currentTags {
		currentTagMap[tag.Name] = tag
	}

	// Карта новых тегов
	newTagMap := make(map[string]bool)
	for _, tagName := range tagNames {
		newTagMap[tagName] = true
	}

	// Удаляем теги, которых нет в новом списке
	for _, tag := range currentTags {
		if _, exists := newTagMap[tag.Name]; !exists {
			if err := uc.tagRepo.RemoveTagFromArticle(ctx, article.ID, tag.ID); err != nil {
				return err
			}

		}
	}

	// Добавляем новые теги
	for tagName := range newTagMap {
		if _, exists := currentTagMap[tagName]; !exists {
			tag, err := uc.tagRepo.GetByName(ctx, tagName)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					tag = &entity.Tag{
						Name:      tagName,
						CreatedAt: time.Now(),
					}
					if err := uc.tagRepo.Create(ctx, tag); err != nil {
						return err
					}
				} else {
					return err
				}
			}
			// Связываем тег со статьей
			if err := uc.tagRepo.AddTagToArticle(ctx, article.ID, tag.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

// ArchiveArticle переносит статью в архив
func (uc *ArticleUseCase) ArchiveArticle(ctx context.Context, id uint) error {
	article, err := uc.articleRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	article.Status = entity.StatusDeprecated
	article.UpdatedAt = time.Now()

	return uc.articleRepo.Update(ctx, article)
}

// Search ищет статьи по тексту и тегам
func (uc *ArticleUseCase) Search(ctx context.Context, query string, tags []string, status entity.ArticleStatus) ([]*entity.Article, error) {
	return uc.articleRepo.Search(ctx, query, tags, status)
}

//func (uc *ArticleUseCase) Delete(id uint) error { return uc.articleRepo.Delete(id) }

// List выводит список статей с лимтом и офсетом
func (uc *ArticleUseCase) List(ctx context.Context, limit, offset int) ([]*entity.Article, error) {
	return uc.articleRepo.List(ctx, limit, offset)
}

//
