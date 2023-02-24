package client

import (
	"encoding/json"
	"fmt"
	"github.com/bwesterb/go-pow"
	"net"
	"word-of-wisdom/shared/events"
	"word-of-wisdom/shared/events/client"
	"word-of-wisdom/shared/events/server"
	"word-of-wisdom/shared/transport"
)

type Client struct {
	address string

	transport *transport.Transport

	conn *net.TCPConn

	onUpdate func([]byte)

	onAuth func()
}

func NewClient(address string) *Client {
	client := &Client{
		address: address,
	}

	return client
}

func (c *Client) SetOnUpdate(fn func([]byte)) *Client {
	c.onUpdate = fn

	return c
}

func (c *Client) SetOnAuth(fn func()) *Client {
	c.onAuth = fn

	return c
}

func (c *Client) WriteBytes(data []byte) error {
	_, err := c.transport.WriteNextBytes(data)

	return err
}

func (c *Client) Connect() error {
	addr, err := net.ResolveTCPAddr("tcp", c.address)
	if err != nil {
		return err
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}

	c.conn = conn
	c.transport = transport.NewTransport(c.conn)

	err = c.handshake()
	if err != nil {
		return err
	}

	return c.handleUpdates()
}

func (c *Client) handshake() error {
	event, err := c.getNextIncomingEvent()
	if err != nil {
		return err
	}

	if event.Type != "challenge" {
		return fmt.Errorf("invalid_event_type_%s_received", event.Type)
	}

	var challenge *server.Challenge

	err = event.DataToStruct(&challenge)
	if err != nil {
		return err
	}

	result, err := pow.Fulfil(challenge.Hash, challenge.Nounce)
	if err != nil {
		return err
	}

	js, err := events.NewOutgoing(client.NewChallengeResponse(result)).ToJSON()
	if err != nil {
		return err
	}

	_ = c.WriteBytes(js)

	event, err = c.getNextIncomingEvent()
	if err != nil {
		return err
	}

	if event.Type != "authorized" {
		return fmt.Errorf("invalid_event_type_%s_received", event.Type)
	}

	if c.onAuth != nil {
		c.onAuth()
	}

	return nil
}

func (c *Client) getNextIncomingEvent() (*events.Incoming, error) {
	data, err := c.transport.ReadNextBytes()
	if err != nil {
		return nil, err
	}

	var event *events.Incoming

	err = json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (c *Client) handleUpdates() error {
	ch := make(chan []byte)

	go c.processUpdates(ch)

	return c.transport.ReceiveBytes(ch)
}

func (c *Client) processUpdates(ch <-chan []byte) {
	for data := range ch {
		if c.onUpdate != nil {
			c.onUpdate(data)
		}
	}
}
