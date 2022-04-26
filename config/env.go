package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	err error
)

func Get(key string) string {
	return getEnv(".env", key)
}

func getEnv(envFile string, key string) string {
	err = godotenv.Load(envFile)

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
