package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type R2Config struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	PublicEndpoint  string
}

type Config struct {
	Port string
	DB   DBConfig
	R2   R2Config
	URL  Url
}

type Url struct {
	BaseUrl string
}

var AppConfig *Config

func Load() *Config {
	// Load .env file if exists
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, relying on environment variables")
	}

	AppConfig = &Config{
		Port: getEnv("PORT", "8080"),
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "shortlink"),
		},
		R2: R2Config{
			AccountID:       getEnv("R2_ACCOUNT_ID", ""),
			AccessKeyID:     getEnv("R2_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("R2_SECRET_ACCESS_KEY", ""),
			BucketName:      getEnv("R2_BUCKET_NAME", ""),
			PublicEndpoint:  getEnv("R2_PUBLIC_ENDPOINT", ""),
		},
		URL: Url{
			BaseUrl: getEnv("BASE_URL", "127.0.0.1"),
		},
	}

	return AppConfig
}

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
