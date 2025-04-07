package authserver

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var configPath = "configs/.env.auth"

type envConfig struct {
	ServerAddress         string `envconfig:"SERVER_ADDRESS" default:"localhost:8080"`
	LogLevel              string `envconfig:"LOG_LEVEL" default:"info"`
	DatabaseDSN           string `envconfig:"DATABASE_DSN"`
	TokenSecret           string `envconfig:"TOKEN_SECRET"`
	DefaultRequestTimeout string `envconfig:"DEFAULT_REQUEST_TIMEOUT"`
}

func NewConfig() *envConfig {
	env := os.Getenv("GO_ENV")

	if env == "" {
		env = ".development"
	}

	godotenv.Load(configPath + env)

	var cfg envConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal("Failed to process env vars: ", err)
	}

	fmt.Println("LogLevel", cfg.LogLevel)
	return &cfg
}
