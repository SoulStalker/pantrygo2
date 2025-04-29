package entity

import "time"

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
)

type MediaFile struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Filename  string    `json:"filename" gorm:"not null"`
	Path      string    `json:"path" gorm:"not null"`
	Type      MediaType `json:"type" gorm:"type:varchar(10);not null"`
	Size      int64     `json:"size"`
	ArticleID uint      `json:"article_id" gorm:"index"`
	CreatedAt time.Time `json:"created_at"`
}
