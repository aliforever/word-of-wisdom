package main

import (
	"fmt"
	"os"
	"word-of-wisdom/libs/server"
	"word-of-wisdom/libs/server/config"
	"word-of-wisdom/libs/server/quotestorage"
	"word-of-wisdom/shared/events"
	serverEvent "word-of-wisdom/shared/events/server"
)

func main() {
	os.Setenv("ADDRESS", ":7000")

	quotes := []string{
		"Knowing yourself is the beginning of all wisdom. Aristotle",
		"Courage is the first of human virtues because it makes all others possible. Aristotle",
		"The best and most beautiful things in the world cannot be seen or even touched – they must be felt with the heart. Aristotle",
		"To know yourself, you must spend time with yourself, you must not be afraid to be alone. Aristotle",
		"The habits we form from childhood make no small difference, but rather they make all the difference. Aristotle",
		"Excellence is an art won by training and habituation. Aristotle",
		"I count him braver who overcomes his desires than him who conquers his enemies, for the hardest victory is over self. Aristotle",
		"No great mind has ever existed without a touch of madness. Aristotle",
		"Through discipline comes freedom. Aristotle",
		"Well begun is half done. Aristotle",
	}

	cfg, err := config.LoadFromEnv()
	if err != nil {
		panic(err)
	}

	quoteStorage := quotestorage.NewQuoteStorage(quotes...)

	srv := server.NewServer(cfg).
		SetOnConnect(func(conn *server.Connection) {
			err := conn.SetOnUpdate(func(bytes []byte) {
				event, err := events.IncomingFromJSON(bytes)
				if err != nil {
					fmt.Println("invalid event received", err)
					return
				}

				if event.Type != "new_quote" {
					fmt.Println("unknown event received", event.Type)
					return
				}

				quote, err := quoteStorage.GetRandomQuote()
				if err != nil {
					fmt.Println("cant get quote")
					return
				}

				js, _ := events.NewOutgoing(serverEvent.NewQuote(quote)).ToJSON()

				_ = conn.WriteBytes(js)
			}).SetOnError(func(err error) {
				fmt.Println("error from connection:", err)
			}).Start()
			if err != nil {
				fmt.Println("connection closed:", err)
			}
		})

	panic(srv.Start())
}
