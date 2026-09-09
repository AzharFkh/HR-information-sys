package config

import (
	"os"
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Println("Failed to load .env:", err)
		}
	}
}