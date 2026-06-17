package validator

import (
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/cobrich/youtube-audio-looper/internal/entity"
	"github.com/cobrich/youtube-audio-looper/internal/utils"
)

const (
	MaxDurationSeconds   = 600
	MaxClipLengthSeconds = 60
)

func ValidateCreateLoopedAudioRequest(req entity.Request) error {
	if !isValidURL(req.YoutubeURL) {
		return errors.New("invalid youtube url")
	}

	startSeconds, err := utils.ParseTimeToSeconds(req.Start)
	if err != nil {
		log.Println(err)

		return err
	}

	endSeconds, err := utils.ParseTimeToSeconds(req.End)
	if err != nil {
		log.Println(err)

		return err
	}

	if endSeconds <= startSeconds {
		log.Println(err)

		return errors.New("end time must be older than start")
	}

	durationSeconds, err := utils.ParseTimeToSeconds(req.Duration)
	if err != nil {
		log.Println(err)

		return err
	}

	if durationSeconds <= 0 {
		log.Println(err)

		return errors.New("duration must be grater than zero")
	}

	if endSeconds-startSeconds > durationSeconds {
		log.Println(err)

		return errors.New("selected segment cannot be longer than final duration")
	}

	if durationSeconds > MaxDurationSeconds {
		return errors.New("maximum duration is 10 minutes")
	}

	if endSeconds-startSeconds > MaxClipLengthSeconds {
		return errors.New("the clip length should be less than minute")
	}

	return nil
}

func isValidURL(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		log.Println(err)

		return false
	}

	host := strings.ToLower(parsedURL.Host)

	return host == "youtube.com" ||
		host == "www.youtube.com" ||
		host == "youtu.be" ||
		host == "m.youtube.com" ||
		host == "music.youtube.com"
}
