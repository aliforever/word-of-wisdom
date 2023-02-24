package main

import (
	"encoding/json"
	"fmt"
	"time"
	"word-of-wisdom/libs/client"
	"word-of-wisdom/shared/events"
	clientEvents "word-of-wisdom/shared/events/client"
	"word-of-wisdom/shared/events/server"
)

func main() {
	client := client.NewClient("127.0.0.1:7000")

	client.SetOnAuth(func() {
		js, _ := events.NewOutgoing(clientEvents.NewNewQuote()).ToJSON()

		_ = client.WriteBytes(js)
	})

	client.SetOnUpdate(func(data []byte) {
		var incomingEvent *events.Incoming

		err := json.Unmarshal(data, &incomingEvent)
		if err != nil {
			return
		}

		if incomingEvent.Type != "quote" {
			fmt.Println("received unknown event from server", incomingEvent.Type)
			return
		}

		var quote *server.Quote

		err = incomingEvent.DataToStruct(&quote)
		if err != nil {
			fmt.Println("invalid json received", err)
			return
		}

		fmt.Println("received quote:", quote.Text)
	})

	err := client.Connect()
	if err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 3)
	panic(err)
}
