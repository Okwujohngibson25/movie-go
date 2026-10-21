package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Port          string
	Env           string
	JWTSecret     string
	AdminEmail    string
	AdminPassword string
}

type DatabaseConfig struct {
	URL string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		App: AppConfig{
			Port:          os.Getenv("PORT"),
			Env:           os.Getenv("APP_ENV"),
			JWTSecret:     os.Getenv("JWT_SECRET"),
			AdminEmail:    os.Getenv("Admin_email"),
			AdminPassword: os.Getenv("Admin_Password"),
		},
		Database: DatabaseConfig{
			URL: os.Getenv("DATABASE_URL"),
		},
	}, nil
}
