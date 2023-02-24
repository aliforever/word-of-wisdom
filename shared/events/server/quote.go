package server

import "word-of-wisdom/shared/events"

type Quote struct {
	Text string `json:"text"`
}

func (q Quote) Type() string {
	return "quote"
}

func NewQuote(text string) events.OutgoingI {
	return Quote{Text: text}
}
