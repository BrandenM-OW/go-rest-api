package initializers

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnvVars(path string) {
	err := godotenv.Load(filepath.Join(path, "meta/.env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
