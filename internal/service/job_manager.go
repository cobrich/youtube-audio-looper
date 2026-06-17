package service

import (
	"os"
	"path/filepath"

	"github.com/cobrich/youtube-audio-looper/internal/entity"
	"github.com/google/uuid"
)

type JobManager struct {
	baseDir string
}

func NewJobManager(baseDir string) *JobManager {
	return &JobManager{
		baseDir: baseDir,
	}
}

func (m *JobManager) Create() (*entity.Job, error) {
	id := uuid.New().String()

	jobDir := filepath.Join(m.baseDir, id)

	if err := os.MkdirAll(jobDir, 0755); err != nil {
		return nil, err
	}

	return &entity.Job{
		ID:  id,
		Dir: jobDir,
	}, nil
}
