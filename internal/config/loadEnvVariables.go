package config

import(
	"log"
	"github.com/joho/godotenv"
)

// LoadEnvVariables loads environment variables from a .env file
func LoadEnvVariables() {
  err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
}
