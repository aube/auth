package main

import (
	"context"
	"log"
	"os"

	"github.com/aube/auth/internal/application/user"
	"github.com/aube/auth/internal/infrastructure/postgres"
	"github.com/aube/auth/internal/interfaces/rest"
)

func main() {
	ctx := context.Background()

	// Инициализация БД
	pgConfig := postgres.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  "disable",
	}

	dbPool, err := postgres.NewPool(ctx, pgConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Инициализация репозитория и сервиса
	userRepo := postgres.NewUserRepository(dbPool)
	userService := user.NewUserService(userRepo)

	// Запуск сервера
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-key" // В продакшене используйте надежный секрет
	}

	server := rest.NewServer(userService, jwtSecret)
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
