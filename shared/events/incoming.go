package events

import "encoding/json"

type Incoming struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

func IncomingFromJSON(js []byte) (*Incoming, error) {
	var incoming *Incoming

	err := json.Unmarshal(js, &incoming)
	if err != nil {
		return nil, err
	}

	return incoming, nil
}

func (i *Incoming) DataToStruct(target interface{}) error {
	return json.Unmarshal(i.Data, &target)
}
