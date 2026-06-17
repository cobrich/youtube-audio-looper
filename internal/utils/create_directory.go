package utils

import (
	"log"
	"os"
)

func CreateDirectories() {
	dirs := []string{
		"tmp/downloads",
		"tmp/clips",
		"tmp/results",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatal(err)
		}
	}
}

func DeleteDirectories() {
	dirs := []string{
		"tmp/downloads",
		"tmp/clips",
		"tmp/results",
	}

	for _, dir := range dirs {
		defer os.RemoveAll(dir)
	}
}
