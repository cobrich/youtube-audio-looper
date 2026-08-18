package downloader

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
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

func (d *YouTubeDownloader) DownloadAudioSegment(ctx context.Context, youtubeURL, outputPath, start, end string) error {
	downloadURL := normalizeYouTubeURL(youtubeURL)
	section := fmt.Sprintf("*%s-%s", normalizeTimestamp(start), normalizeTimestamp(end))

	args := []string{
		"yt-dlp",
		"--no-playlist",
		"--download-sections", section,
		"-f", "bestaudio/best",
		"-x",
		"--audio-format", "mp3",
		"--audio-quality", "0",
		"-o", outputPath,
		downloadURL,
	}

	if cookiesPath := os.Getenv("YTDLP_COOKIES_PATH"); cookiesPath != "" {
		tempCookiesPath, err := copyCookiesToTempFile(cookiesPath)
		if err != nil {
			return err
		}
		defer os.Remove(tempCookiesPath)

		log.Printf("yt-dlp cookies configured path=%s", cookiesPath)
		args = append([]string{"yt-dlp", "--cookies", tempCookiesPath}, args[1:]...)
	} else {
		log.Println("yt-dlp cookies not configured")
	}

	if userAgent := os.Getenv("YTDLP_USER_AGENT"); userAgent != "" {
		log.Println("yt-dlp user agent configured")
		args = append([]string{"yt-dlp", "--user-agent", userAgent}, args[1:]...)
	}

	if _, err := exec.LookPath("node"); err == nil {
		args = append([]string{"yt-dlp", "--remote-components", "ejs:github", "--js-runtimes", "node"}, args[1:]...)
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

func (d *YouTubeDownloader) DownloadAudio(ctx context.Context, youtubeURL string, outputPath string) error {
	return d.DownloadAudioSegment(ctx, youtubeURL, outputPath, "00:00", "inf")
}

func copyCookiesToTempFile(cookiesPath string) (string, error) {
	data, err := os.ReadFile(cookiesPath)
	if err != nil {
		return "", fmt.Errorf("%w: %s: %v", ErrCookieFileUnavailable, cookiesPath, err)
	}

	tempFile, err := os.CreateTemp("", "yt-dlp-cookies-*.txt")
	if err != nil {
		return "", fmt.Errorf("create temporary cookie file: %w", err)
	}
	defer tempFile.Close()

	if _, err := tempFile.Write(data); err != nil {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("write temporary cookie file: %w", err)
	}

	return tempFile.Name(), nil
}

func normalizeYouTubeURL(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	if strings.EqualFold(parsedURL.Host, "music.youtube.com") {
		parsedURL.Host = "www.youtube.com"
	}

	return parsedURL.String()
}

func normalizeTimestamp(value string) string {
	if value == "inf" || strings.Count(value, ":") == 2 {
		return value
	}

	return "00:" + value
}
