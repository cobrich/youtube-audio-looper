package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	AudioAPIURL   string
}

func MustLoad() *Config {
	godotenv.Load()
	cfg := &Config{
		TelegramToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		AudioAPIURL:   os.Getenv("AUDIO_API_URL"),
	}

	if cfg.TelegramToken == "" {
		log.Fatalf("TELEGRAM_BOT_TOKEN is required")
	}

	if cfg.AudioAPIURL == "" {
		cfg.AudioAPIURL = "http://localhost:8084"
	}

	return cfg
}
