package events

import "encoding/json"

type OutgoingI interface {
	Type() string
}

type Outgoing struct {
	Type string    `json:"type"`
	Data OutgoingI `json:"data"`
}

func NewOutgoing(data OutgoingI) *Outgoing {
	return &Outgoing{
		Type: data.Type(),
		Data: data,
	}
}

func (e *Outgoing) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
