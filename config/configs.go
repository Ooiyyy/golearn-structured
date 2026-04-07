// digunakan untuk menyatukan beberapa config
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App AppConfig
	DB  DatabaseConfig
}

func LoadEnv() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		App: AppConfig{
			AppPort:   os.Getenv("APP_PORT"),
			JWTSecret: os.Getenv("JWT_SECRET"),
		},
		DB: DatabaseConfig{
			DBHost: os.Getenv("DB_HOST"),
			DBPort: os.Getenv("DB_PORT"),
			DBUser: os.Getenv("DB_USER"),
			DBPass: os.Getenv("DB_PASS"),
			DBName: os.Getenv("DB_NAME"),
		},
	}
}
