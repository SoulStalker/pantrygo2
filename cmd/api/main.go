package main

import (
	"github.com/soulstalker/pantrygo2/internal/repository/postgres"
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

}
