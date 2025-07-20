package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aube/auth/internal/api/rest"
	appFile "github.com/aube/auth/internal/application/file"
	appUpload "github.com/aube/auth/internal/application/upload"
	appUser "github.com/aube/auth/internal/application/user"
	"github.com/aube/auth/internal/infrastructure/fs"
	"github.com/aube/auth/internal/infrastructure/postgres"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/spf13/viper"
)

func main() {
	ctx := context.Background()
	viper.SetConfigFile(".env")
	viper.SetDefault("STORAGE_PATH", "./_storage")
	viper.SetDefault("API_PATH", "/api/v1")
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.ReadInConfig()

	logger.Init(viper.Get("LOG_LEVEL").(string))

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
	fsRepo, err := fs.NewFileSystemRepository(storagePath)
	if err != nil {
		log.Fatalf("Failed to initialize file repository: %v", err)
	}

	fileService := appFile.NewFileService(fsRepo)

	uploadRepo := postgres.NewUploadRepository(dbPool)
	userRepo := postgres.NewUserRepository(dbPool)
	uploadService := appUpload.NewUploadService(uploadRepo)
	userService := appUser.NewUserService(userRepo)

	// Запуск сервера
	jwtSecret := viper.Get("JWT_SECRET").(string)
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET not found")
	}
	apiPath := viper.Get("API_PATH").(string)

	server := rest.NewServer(userService, fileService, uploadService, jwtSecret, apiPath)
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
