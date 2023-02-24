package server

import (
	"errors"
	"net"
	"word-of-wisdom/libs/server/config"
)

var (
	connectionClosed = errors.New("connection_closed")
)

type Server struct {
	cfg *config.Config

	server *net.TCPListener

	close chan bool

	onConnect func(conn *Connection)
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg:   cfg,
		close: make(chan bool),
	}
}

func (s *Server) SetOnConnect(fn func(conn *Connection)) *Server {
	s.onConnect = fn

	return s
}

func (s *Server) Start() error {
	var (
		tcpAddr *net.TCPAddr
		err     error
	)

	tcpAddr, err = net.ResolveTCPAddr("tcp", s.cfg.Address)
	if err != nil {
		return err
	}

	s.server, err = net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	return s.receiveConnections()
}

func (s *Server) Close() error {
	go s.closeListener()

	return s.server.Close()
}

func (s *Server) closeListener() {
	s.close <- true
}

func (s *Server) receiveConnections() error {
	var (
		tcpConn *net.TCPConn
		err     error
	)

	for {
		tcpConn, err = s.server.AcceptTCP()
		if err != nil {
			break
		}

		go s.addConnection(NewConnection(tcpConn))
	}

	return err
}

func (s *Server) addConnection(conn *Connection) {
	if s.onConnect != nil {
		s.onConnect(conn)
	}
}
