package main

import (
	"log"
	"os"

	"github.com/cobrich/youtube-audio-looper/internal/handler"
	"github.com/cobrich/youtube-audio-looper/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	audioService := service.NewAudioService()
	audioHandler := handler.NewAudioHandler(audioService)

	router.Static("/web", "./web")
	router.GET("/", func(c *gin.Context) {
		c.File("./web/index.html")
	})

	router.GET("/health", handler.Health)

	api := router.Group("api/v1")
	{
		api.POST("/audio/loop", audioHandler.CreateLoopedAudio)
	}

	log.Println("server started")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
