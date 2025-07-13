package main

import (
	"context"
	"fmt"
	"log"
	"os"

	appFile "github.com/aube/auth/internal/application/file"
	appUser "github.com/aube/auth/internal/application/user"
	"github.com/aube/auth/internal/infrastructure/fs"
	"github.com/aube/auth/internal/infrastructure/postgres"
	"github.com/aube/auth/internal/interfaces/rest"
	"github.com/spf13/viper"
)

func main() {
	ctx := context.Background()
	viper.SetConfigFile(".env")
	viper.SetDefault("STORAGE_PATH", "./_storage")
	viper.SetDefault("API_PATH", "/api/v1")
	viper.ReadInConfig()

	// Инициализация БД
	pgConfig := postgres.Config{
		Host:     viper.Get("DB_HOST").(string),
		Port:     viper.Get("DB_PORT").(string),
		User:     viper.Get("DB_USER").(string),
		Password: viper.Get("DB_PASSWORD").(string),
		DBName:   viper.Get("DB_NAME").(string),
		SSLMode:  "disable",
	}
	fmt.Println(pgConfig)

	dbPool, err := postgres.NewPool(ctx, pgConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Инициализация хранилища файлов
	storagePath := viper.Get("STORAGE_PATH").(string)
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		log.Fatalf("Failed to create storage directory: %v", err)
	}
	fsRepo, err := fs.NewFileSystemRepository(storagePath)
	if err != nil {
		log.Fatalf("Failed to initialize file repository: %v", err)
	}

	fileService := appFile.NewFileService(fsRepo)

	// Инициализация репозитория и сервиса
	userRepo := postgres.NewUserRepository(dbPool)
	userService := appUser.NewUserService(userRepo)

	// Запуск сервера
	jwtSecret := viper.Get("JWT_SECRET").(string)
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET not found")
	}
	apiPath := viper.Get("API_PATH").(string)

	server := rest.NewServer(userService, fileService, jwtSecret, apiPath)
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
