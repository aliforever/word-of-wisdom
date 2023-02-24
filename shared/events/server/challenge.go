package server

import "word-of-wisdom/shared/events"

type Challenge struct {
	Hash   string `json:"hash"`
	Nounce []byte `json:"nounce"`
}

func (h Challenge) Type() string {
	return "challenge"
}

func NewChallenge(hash string, nounce []byte) events.OutgoingI {
	return Challenge{Hash: hash, Nounce: nounce}
}
