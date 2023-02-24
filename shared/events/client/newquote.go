package client

import "word-of-wisdom/shared/events"

type NewQuote struct {
}

func (n NewQuote) Type() string {
	return "new_quote"
}

func NewNewQuote() events.OutgoingI {
	return NewQuote{}
}
