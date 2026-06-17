# YouTube Audio Looper

A Go-based service and Telegram bot that extracts audio segments from YouTube videos, loops them to a user-defined duration, and returns the final MP3 file.

The project consists of two components:

* REST API for audio processing
* Telegram Bot for user interaction

---

## Why I Built This

Sometimes a song contains only a small section that I want to listen to repeatedly, while the rest of the track is less interesting.

Instead of manually downloading and editing audio files, I decided to build a service that automatically:

* downloads audio from YouTube
* extracts the desired segment
* loops it to a specified duration
* delivers the final MP3 directly through Telegram

This project was built to solve a real personal problem while exploring audio processing, external tools integration, Docker, and Telegram bot development with Go.

---

## Screenshot

> Add your Telegram bot screenshots here.

![Telegram Bot Demo](docs/images/demo.png)

---

## Demo

```text
User: https://youtube.com/watch?v=...

Bot: Send start timestamp.

User: 00:42

Bot: Send end timestamp.

User: 00:47

Bot: Send final duration.

User: 01:00

Bot: Processing...

Bot: 🎵 looped-audio.mp3
```

---

## Features

* Download audio from YouTube videos
* Extract any audio segment
* Loop the selected segment to a custom duration
* Generate MP3 output
* Telegram Bot integration
* Docker support
* Automatic temporary file cleanup
* Request validation
* Health check endpoint
* Job-based processing workflow

---

## Architecture

```text
Telegram User
      │
      ▼
Telegram Bot
      │ HTTP
      ▼
Audio API
      │
      ├── yt-dlp
      ├── ffmpeg (cut)
      ├── ffmpeg (loop)
      └── cleanup
      │
      ▼
Generated MP3
      │
      ▼
Telegram User
```

---

## Project Structure

```text
cmd/
├── api/
└── bot/

internal/
├── bot/
├── client/
├── config/
├── downloader/
├── entity/
├── handler/
├── processor/
├── service/
├── utils/
└── validator/
```

---

## API Usage

### Endpoint

```http
POST /api/v1/audio/loop
```

### Request Body

```json
{
  "youtube_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
  "start": "00:42",
  "end": "00:47",
  "duration": "01:00"
}
```

### Response

Returns a generated MP3 file.

---

## Telegram Bot Usage

1. Start the bot using `/start`
2. Send a YouTube URL
3. Send the start timestamp
4. Send the end timestamp
5. Send the desired final duration
6. Receive the generated MP3 file

---

## Running Locally

### Requirements

* Go 1.26+
* ffmpeg
* yt-dlp

### Run API

```bash
go run ./cmd/api
```

### Run Telegram Bot

```bash
export TELEGRAM_BOT_TOKEN=your_token
export AUDIO_API_URL=http://localhost:8084

go run ./cmd/bot
```

---

## Running with Docker

### Build and Start

```bash
docker compose up --build
```

### Health Check

```bash
curl http://localhost:8084/health
```

### Example Request

```bash
curl -X POST http://localhost:8084/api/v1/audio/loop \
  -H "Content-Type: application/json" \
  -d '{
    "youtube_url": "https://www.youtube.com/watch?v=jNQXAC9IVRw",
    "start": "00:00",
    "end": "00:05",
    "duration": "00:30"
  }' \
  --output result.mp3
```

---

## Processing Workflow

```text
Request
    │
    ▼
Create Job Directory
    │
    ▼
Download Audio
    │
    ▼
Extract Segment
    │
    ▼
Loop Segment
    │
    ▼
Generate MP3
    │
    ▼
Return File
    │
    ▼
Cleanup Temporary Files
```

---

## Tech Stack

### Backend

* Go 1.26
* Gin

### Audio Processing

* ffmpeg
* yt-dlp

### Telegram

* Telegram Bot API

### Infrastructure

* Docker
* Docker Compose

---

## Future Improvements

* Railway deployment
* VPS deployment
* Rate limiting
* User quotas
* Audio format selection
* Web interface
* Usage analytics
* User history
* Caching

---

## Author

Bekzat Tursun

GitHub: https://github.com/cobrich

LinkedIn: https://linkedin.com/in/tursunbekzat
