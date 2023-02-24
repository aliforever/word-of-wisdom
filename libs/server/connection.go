package server

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"github.com/bwesterb/go-pow"
	"net"
	"word-of-wisdom/shared/events"
	"word-of-wisdom/shared/events/client"
	"word-of-wisdom/shared/events/server"
	"word-of-wisdom/shared/transport"
)

type state int

const (
	StateChallenge state = iota
	StateAuthorized
)

type Connection struct {
	conn *net.TCPConn

	state state

	transport *transport.Transport

	close chan bool

	updates chan []byte

	hash   string
	nounce []byte

	onUpdate func(update []byte)
	onError  func(err error)
}

func NewConnection(conn *net.TCPConn) *Connection {
	return &Connection{
		conn:      conn,
		transport: transport.NewTransport(conn),
		state:     StateChallenge,
		close:     make(chan bool),
		updates:   make(chan []byte),
	}
}

func (c *Connection) SetOnUpdate(fn func([]byte)) *Connection {
	c.onUpdate = fn

	return c
}

func (c *Connection) SetOnError(fn func(error)) *Connection {
	c.onError = fn

	return c
}

func (c *Connection) Start() error {
	c.nounce = make([]byte, 12)

	_, err := rand.Read(c.nounce)
	if err != nil {
		return err
	}

	c.hash = pow.NewRequest(5, c.nounce)

	js, err := events.NewOutgoing(server.NewChallenge(c.hash, c.nounce)).ToJSON()
	if err != nil {
		return err
	}

	_, _ = c.transport.WriteNextBytes(js)

	return c.handleUpdates()
}

func (c *Connection) Close() error {
	return c.closeConnection()
}

func (c *Connection) Updates() <-chan []byte {
	return c.updates
}

func (c *Connection) WriteBytes(data []byte) error {
	_, err := c.transport.WriteNextBytes(data)

	return err
}

func (c *Connection) closeConnection() error {
	c.close <- true

	return c.conn.Close()
}

func (c *Connection) validateHash(incomingEvent *events.Incoming) error {
	var challenge *client.ChallengeResponse

	err := incomingEvent.DataToStruct(&challenge)
	if err != nil {
		return err
	}

	if challenge.Response == "" {
		return errors.New("empty_challenge")
	}

	result, err := pow.Check(c.hash, challenge.Response, c.nounce)
	if err != nil {
		return err
	}

	if !result {
		return errors.New("invalid_result")
	}

	return nil
}

func (c *Connection) handleUpdates() error {
	updates := make(chan []byte)

	go c.processUpdates(updates)

	return c.transport.ReceiveBytes(updates)
}

func (c *Connection) processUpdates(ch <-chan []byte) {
	for data := range ch {
		var incomingEvent *events.Incoming

		incomingEvent, err := events.IncomingFromJSON(data)
		if err != nil {
			continue
		}

		switch c.state {
		case StateChallenge:
			err := c.validateHash(incomingEvent)
			if err != nil {
				go c.reportError(err)
				return
			}

			c.notifyAuthorization()
		case StateAuthorized:
			if c.onUpdate != nil {
				c.onUpdate(data)
			}
		}
	}
}

func (c *Connection) notifyAuthorization() {
	c.state = StateAuthorized

	j, _ := json.Marshal(events.NewOutgoing(server.NewAuthorized()))

	_, _ = c.transport.WriteNextBytes(j)
}

func (c *Connection) reportError(err error) {
	if c.onError != nil {
		c.onError(err)
	}
}
