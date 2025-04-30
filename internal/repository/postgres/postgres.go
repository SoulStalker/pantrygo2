package postgres

import (
	"fmt"
	"github.com/soulstalker/pantrygo2/internal/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(host, user, password, dbname string, port int) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	err = db.AutoMigrate(&entity.Article{}, &entity.Tag{}, &entity.ArticleTag{}, &entity.MediaFile{})
	if err != nil {
		return nil, fmt.Errorf("failed ti migrate database: %w", err)
	}

	return db, nil
}
