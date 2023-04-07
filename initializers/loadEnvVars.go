package initializers

import (
    "log"
    "github.com/joho/godotenv"
    "path/filepath"
)

func LoadEnvVars(path string) {
    err := godotenv.Load(filepath.Join(path, ".env"))
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}
