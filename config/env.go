package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file
func LoadEnv() error {
	// Try to load .env file, but don't fail if it doesn't exist (in production)
	err := godotenv.Load()
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("⚠️ .env file not found, using environment variables")
			return nil
		}
		return err
	}
	log.Println("✅ .env file loaded successfully")
	return nil
}
