package client

import "word-of-wisdom/shared/events"

type ChallengeResponse struct {
	Response string `json:"response"`
}

func (h ChallengeResponse) Type() string {
	return "challenge_response"
}

func NewChallengeResponse(response string) events.OutgoingI {
	return ChallengeResponse{Response: response}
}
