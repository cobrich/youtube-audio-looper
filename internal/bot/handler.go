package bot

import (
	"fmt"
	"strings"

	"github.com/cobrich/youtube-audio-looper/internal/client"
	"github.com/cobrich/youtube-audio-looper/internal/entity"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	api *tgbotapi.BotAPI

	sessions map[int64]*Session

	audioAPIClient *client.AudioAPIClient
}

func NewHandler(api *tgbotapi.BotAPI, client *client.AudioAPIClient) *Handler {
	return &Handler{
		api:            api,
		sessions:       make(map[int64]*Session),
		audioAPIClient: client,
	}
}

func (h *Handler) HandleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	text := strings.TrimSpace(update.Message.Text)

	if text == "/start" {
		h.handleStart(chatID)
		return
	}

	h.handleText(chatID, text)
}

func (h *Handler) handleStart(chatID int64) {
	h.sessions[chatID] = &Session{
		State: StateWaitingURL,
	}

	msg := tgbotapi.NewMessage(chatID, "Ассалаумағаллейкум, ютубқа сілтеме жүктеңіз...")

	h.api.Send(msg)
}

func (h *Handler) handleText(chatID int64, text string) {
	session, ok := h.sessions[chatID]
	if !ok {
		h.handleStart(chatID)
		return
	}

	switch session.State {
	case StateWaitingURL:
		session.YoutubeURL = text
		session.State = StateWaitingStart

		msg := tgbotapi.NewMessage(chatID, "Басталу уақытын енгізіңіз, мысалы: 00:15")
		h.api.Send(msg)
	case StateWaitingStart:
		session.Start = text
		session.State = StateWaitingEnd

		msg := tgbotapi.NewMessage(chatID, "Аяқталу уақытын енгізіңіз, мысалы: 01:00")
		h.api.Send(msg)
	case StateWaitingEnd:
		session.End = text
		session.State = StateWaitingDuration

		msg := tgbotapi.NewMessage(chatID, "Жазба ұзақтығын таңдаңыз, мысалы: 03:00")
		h.api.Send(msg)
	case StateWaitingDuration:
		session.Duration = text
		session.State = StateIdle

		msg := tgbotapi.NewMessage(chatID, "Ақпараттар дұрыстығын тексеріп, жазбаңызды дайындап жіберемін. Өтініш, күтіңіз...")
		h.api.Send(msg)

		req := entity.Request{
			YoutubeURL: session.YoutubeURL,
			Start:      session.Start,
			End:        session.End,
			Duration:   session.Duration,
		}

		fmt.Println(session)
		fmt.Println(req)

		audioData, err := h.audioAPIClient.CreateLoopedAudio(req)
		if err != nil {
			h.api.Send(tgbotapi.NewMessage(chatID, "Жазба шығарылмады, ақпараттар дұрыстыңын тексеруіңізді сұраймын."))
			return
		}

		audioFile := tgbotapi.FileBytes{
			Name:  "looped-audio.mp3",
			Bytes: audioData,
		}

		audioMsg := tgbotapi.NewAudio(chatID, audioFile)
		audioMsg.Caption = "Жазба дайын"

		h.api.Send(audioMsg)

	default:
		h.handleStart(chatID)
	}
}
