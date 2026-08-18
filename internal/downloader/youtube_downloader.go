package downloader

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

var (
	ErrAuthenticationRequired = errors.New("youtube authentication required")
	ErrCookieFileUnavailable  = errors.New("yt-dlp cookie file unavailable")
)

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
		info, err := os.Stat(cookiesPath)
		if err != nil {
			return fmt.Errorf("%w: %s: %v", ErrCookieFileUnavailable, cookiesPath, err)
		}

		log.Printf("yt-dlp cookies configured path=%s size=%d", cookiesPath, info.Size())
		args = append([]string{"yt-dlp", "--cookies", cookiesPath}, args[1:]...)
	} else {
		log.Println("yt-dlp cookies not configured")
	}

	if userAgent := os.Getenv("YTDLP_USER_AGENT"); userAgent != "" {
		log.Println("yt-dlp user agent configured")
		args = append([]string{"yt-dlp", "--user-agent", userAgent}, args[1:]...)
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
