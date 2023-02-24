package server

import "word-of-wisdom/shared/events"

type Authorized struct{}

func (a Authorized) Type() string {
	return "authorized"
}

func NewAuthorized() events.OutgoingI {
	return Authorized{}
}
