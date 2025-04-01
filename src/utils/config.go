package utils

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"os"
)

func GetEnvConfig(key, fallback string) string {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
