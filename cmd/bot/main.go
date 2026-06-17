package main

import (
	"log"

	"github.com/cobrich/youtube-audio-looper/internal/bot"
	"github.com/cobrich/youtube-audio-looper/internal/client"
	"github.com/cobrich/youtube-audio-looper/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	cfg := config.MustLoad()

	api, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	audioAPIClient := client.NewAudioAPIClient(cfg.AudioAPIURL)
	handler := bot.NewHandler(api, audioAPIClient)

	_ = handler

	log.Printf("authorized as %s", api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := api.GetUpdatesChan(u)

	for update := range updates {
		handler.HandleUpdate(update)
	}
}
