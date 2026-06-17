package utils

import (
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var TimeRegex = regexp.MustCompile(`^([0-9]{2}:)?[0-9]{2}:[0-9]{2}$`)

func ParseTimeToSeconds(value string) (int, error) {
	if !TimeRegex.MatchString(value) {
		return 0, errors.New("invalid time format")
	}

	parts := strings.Split(value, ":")

	if len(parts) == 2 {
		minutes, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Println(err)

			return 0, err
		}

		if minutes >= 60 {
			return 0, errors.New("minutes must be less than 60")
		}

		seconds, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Println(err)

			return 0, err
		}

		if seconds >= 60 {
			return 0, errors.New("seconds must be less than 60")
		}

		return minutes*60 + seconds, nil
	}

	if len(parts) == 3 {
		hours, err := strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}

		if hours >= 24 {
			return 0, errors.New("minutes must be less than 24")
		}

		minutes, err := strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}

		if minutes >= 60 {
			return 0, errors.New("minutes must be less than 60")
		}

		seconds, err := strconv.Atoi(parts[2])
		if err != nil {
			return 0, err
		}

		if seconds >= 60 {
			return 0, errors.New("seconds must be less than 60")
		}

		return hours*3600 + minutes*60 + seconds, nil
	}

	return 0, errors.New("invalid time format")
}
