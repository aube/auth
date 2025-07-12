package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aube/auth/internal/application/user"
	"github.com/aube/auth/internal/infrastructure/postgres"
	"github.com/aube/auth/internal/interfaces/rest"
	"github.com/spf13/viper"
)

func main() {
	ctx := context.Background()
	viper.SetConfigFile(".env")
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

	// Инициализация репозитория и сервиса
	userRepo := postgres.NewUserRepository(dbPool)
	userService := user.NewUserService(userRepo)

	// Запуск сервера
	jwtSecret := viper.Get("JWT_SECRET").(string)
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET not found")
	}

	server := rest.NewServer(userService, jwtSecret)
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
