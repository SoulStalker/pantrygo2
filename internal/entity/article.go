package entity

import "time"

// ArticleStatus статус статьи
type ArticleStatus string

const (
	StatusDraft      ArticleStatus = "draft"
	StatusPublished  ArticleStatus = "published"
	StatusDeprecated ArticleStatus = "deprecated"
)

// Article статья в базе знаний
type Article struct {
	ID          uint          `json:"id" gorm:"primaryKey"`
	Title       string        `json:"title" gorm:"not null"`
	Content     string        `json:"content" gorm:"type:text"`
	Status      ArticleStatus `json:"status" gorm:"type:varchar(20);default:'draft'"`
	CreatedByID uint          `json:"created_by_id"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}
