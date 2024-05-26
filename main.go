package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
)

const (
	defaultListenAddr = ":5001"
)

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	ln          net.Listener
	peers       map[*Peer]bool
	addPeerChan chan *Peer
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}

	return &Server{
		Config:      cfg,
		peers:       make(map[*Peer]bool),
		addPeerChan: make(chan *Peer),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln

	go s.loop()

	return s.listen()
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addPeerChan:
			s.peers[peer] = true
		default:
			fmt.Println("default case")
		}
	}
}

func (s *Server) listen() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept error", "err", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn)
	s.addPeerChan <- peer

	go peer.read()
}

func main() {
	server := NewServer(Config{})
	log.Fatal(server.Start())
}
