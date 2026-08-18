package downloader

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var ErrAuthenticationRequired = errors.New("youtube authentication required")

type YouTubeDownloader struct{}

func NewYouTubeDownloader() *YouTubeDownloader {
	return &YouTubeDownloader{}
}

func (d *YouTubeDownloader) DownloadAudio(ctx context.Context, youtubeURL string, outputPath string) error {
	args := []string{
		"yt-dlp",
		"--no-playlist",
		"-x",
		"--audio-format", "mp3",
		"--audio-quality", "0",
		"-o", outputPath,
		youtubeURL,
	}

	if cookiesPath := os.Getenv("YTDLP_COOKIES_PATH"); cookiesPath != "" {
		args = append([]string{"yt-dlp", "--cookies", cookiesPath}, args[1:]...)
	}

	if _, err := exec.LookPath("node"); err == nil {
		args = append([]string{"yt-dlp", "--js-runtimes", "node"}, args[1:]...)
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "Sign in to confirm") {
			return fmt.Errorf("%w: %s", ErrAuthenticationRequired, string(output))
		}

		return fmt.Errorf("yt-dlp failed: %w, output: %s", err, string(output))
	}

	return nil
}
