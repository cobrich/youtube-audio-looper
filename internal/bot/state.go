package bot

type State string

const (
	StateIdle            State = "idle"
	StateWaitingURL      State = "waiting_url"
	StateWaitingStart    State = "waiting_start"
	StateWaitingEnd      State = "waiting_end"
	StateWaitingDuration State = "waiting_duration"
)
