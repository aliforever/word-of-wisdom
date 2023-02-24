package quotestorage

import (
	"errors"
	"math/rand"
	"sync"
)

type QuoteStorage struct {
	m sync.Mutex

	quotes []string
}

func NewQuoteStorage(quotes ...string) *QuoteStorage {
	return &QuoteStorage{quotes: quotes}
}

func (q *QuoteStorage) AddQuote(quote string) *QuoteStorage {
	q.m.Lock()
	defer q.m.Unlock()

	q.quotes = append(q.quotes, quote)

	return q
}

func (q *QuoteStorage) GetRandomQuote() (string, error) {
	q.m.Lock()
	defer q.m.Unlock()

	if len(q.quotes) == 0 {
		return "", errors.New("no_quotes_fed")
	}

	return q.quotes[rand.Intn(len(q.quotes))], nil
}
