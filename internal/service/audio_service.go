package service

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/cobrich/youtube-audio-looper/internal/downloader"
	"github.com/cobrich/youtube-audio-looper/internal/entity"
	"github.com/cobrich/youtube-audio-looper/internal/processor"
)

type AudioService struct {
	downloader *downloader.YouTubeDownloader
	processor  *processor.FFmpegProcessor
	jobManager *JobManager
}

func NewAudioService() *AudioService {
	return &AudioService{
		downloader: downloader.NewYouTubeDownloader(),
		processor:  processor.NewFFmpegProcessor(),
		jobManager: NewJobManager("tmp/jobs"),
	}
}

func (s *AudioService) CreateLoopedAudio(req entity.Request) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	job, err := s.jobManager.Create()
	if err != nil {
		return nil, err
	}
	log.Printf(
		"job=%s youtube=%s start=%s end=%s duration=%s",
		job.ID,
		req.YoutubeURL,
		req.Start,
		req.End,
		req.Duration,
	)
	defer os.RemoveAll(job.Dir)
	defer func() {
		log.Printf("job=%s finished", job.ID)
	}()

	originalPath := filepath.Join(job.Dir, "original.mp3")
	clipPath := filepath.Join(job.Dir, "clip.mp3")
	resultPath := filepath.Join(job.Dir, "result.mp3")

	if err := s.downloader.DownloadAudio(ctx, req.YoutubeURL, originalPath); err != nil {
		return nil, err
	}

	if err := s.processor.CutAudio(ctx, originalPath, clipPath, req.Start, req.End); err != nil {
		return nil, err
	}

	if err := s.processor.LoopAudio(ctx, clipPath, resultPath, req.Duration); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(resultPath)
	if err != nil {
		return nil, err
	}

	return data, nil
}
