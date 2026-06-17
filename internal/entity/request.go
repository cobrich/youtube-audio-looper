package entity

type Request struct {
	YoutubeURL string `json:"youtube_url" binding:"required"`
	Start      string `json:"start" binding:"required"`
	End        string `json:"end" binding:"required"`
	Duration   string `json:"duration" binding:"required"`
}
