package main

import (
	"log"

	"github.com/cobrich/youtube-audio-looper/internal/handler"
	"github.com/cobrich/youtube-audio-looper/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	audioService := service.NewAudioService()
	audioHandler := handler.NewAudioHandler(audioService)
	router.GET("/health", handler.Health)

	api := router.Group("api/v1")
	{
		api.POST("/audio/loop", audioHandler.CreateLoopedAudio)
	}

	log.Println("server started")

	if err := router.Run(":8084"); err != nil {
		log.Fatal(err)
	}
}
