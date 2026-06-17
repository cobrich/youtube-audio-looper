package processor

import (
	"context"
	"fmt"
	"os/exec"
)

type FFmpegProcessor struct {
}

func NewFFmpegProcessor() *FFmpegProcessor {
	return &FFmpegProcessor{}
}

func (p *FFmpegProcessor) CutAudio(ctx context.Context, inputPath, outputPath, start, end string) error {
	cmd := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-y",
		"-i", inputPath,
		"-ss", start,
		"-to", end,
		"-vn",
		"-acodec", "libmp3lame",
		"-q:a", "2",
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg cut failed: %w, output: %s", err, string(output))
	}

	return nil
}

func (p *FFmpegProcessor) LoopAudio(ctx context.Context, inputPath string, outputPath string, duration string) error {
	cmd := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-y",
		"-stream_loop", "-1",
		"-i", inputPath,
		"-t", duration,
		"-vn",
		"-acodec", "libmp3lame",
		"-q:a", "2",
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg loop failed: %w, output: %s", err, string(output))
	}

	return nil
}
