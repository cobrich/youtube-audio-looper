package handler

import (
	"log"
	"net/http"

	"github.com/cobrich/youtube-audio-looper/internal/entity"
	"github.com/cobrich/youtube-audio-looper/internal/service"
	"github.com/cobrich/youtube-audio-looper/internal/validator"
	"github.com/gin-gonic/gin"
)

type AudioHandler struct {
	audioService *service.AudioService
}

func NewAudioHandler(audioService *service.AudioService) *AudioHandler {
	return &AudioHandler{
		audioService: audioService,
	}
}

func (h *AudioHandler) CreateLoopedAudio(c *gin.Context) {
	var req entity.Request

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if err := validator.ValidateCreateLoopedAudioRequest(req); err != nil {
		log.Println(err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	audioData, err := h.audioService.CreateLoopedAudio(req)
	if err != nil {
		log.Println("create audio error:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create audio",
		})
		return
	}

	c.Header("Content-Disposition", `attachment; filename="looped-audio.mp3"`)

	c.Data(
		http.StatusOK,
		"audio/mpeg",
		audioData,
	)
}
