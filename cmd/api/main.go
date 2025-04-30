package main

import (
	"github.com/gin-gonic/gin"
	"github.com/soulstalker/pantrygo2/internal/delivery/http"
	"github.com/soulstalker/pantrygo2/internal/repository/postgres"
	"github.com/soulstalker/pantrygo2/internal/usecase"
	"github.com/spf13/viper"
	"log"
)

func main() {
	// Загрузка конфигурации
	viper.SetConfigFile("./config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// Подключение к базе данных
	db, err := postgres.NewPostgresDB(
		viper.GetString("db.host"),
		viper.GetString("db.user"),
		viper.GetString("db.password"),
		viper.GetString("db.dbname"),
		viper.GetInt("db.port"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %w", err)
	}

	articleRepo := postgres.NewArticleRepo(db)
	tagRepo := "todo"
	mediaRepo := "todo"

	articleUseCase := usecase.NewArticleUseCase(articleRepo, tagRepo, mediaRepo)

	articleHandler := http.NewArticleHandler(articleUseCase)

	router := gin.Default()

	articleHandler.RegisterRoutes(router)

	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
