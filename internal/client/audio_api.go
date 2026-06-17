package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cobrich/youtube-audio-looper/internal/entity"
)

type AudioAPIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAudioAPIClient(baseURL string) *AudioAPIClient {
	return &AudioAPIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

func (c *AudioAPIClient) CreateLoopedAudio(req entity.Request) ([]byte, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := c.baseURL + "/api/v1/audio/loop"

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("audio api failed: status=%d body=%s", resp.StatusCode, string(data))
	}

	return data, nil
}
