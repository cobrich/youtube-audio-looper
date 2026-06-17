package downloader

import (
	"context"
	"fmt"
	"os/exec"
)

type YouTubeDownloader struct{}

func NewYouTubeDownloader() *YouTubeDownloader {
	return &YouTubeDownloader{}
}

func (d *YouTubeDownloader) DownloadAudio(ctx context.Context, youtubeURL string, outputPath string) error {
	cmd := exec.CommandContext(
		ctx,
		"yt-dlp",
		"-x",
		"--audio-format", "mp3",
		"--audio-quality", "0",
		"-o", outputPath,
		youtubeURL,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("yt-dlp failed: %w, output: %s", err, string(output))
	}

	return nil
}
