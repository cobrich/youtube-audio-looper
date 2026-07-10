FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o audio-looper ./cmd/api


FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache \
    ffmpeg \
    python3 \
    py3-pip \
    ca-certificates

RUN pip install --break-system-packages -U yt-dlp

COPY --from=builder /app/audio-looper .
COPY --from=builder /app/web ./web

EXPOSE 10000

CMD ["./audio-looper"]
