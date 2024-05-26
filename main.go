package main

import (
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
	quit        chan struct{}
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}

	return &Server{
		Config:      cfg,
		peers:       make(map[*Peer]bool),
		addPeerChan: make(chan *Peer),
		quit:        make(chan struct{}),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln

	go s.loop()

	slog.Info("server running", "ListenAddr", s.ListenAddr)

	return s.listen()
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addPeerChan:
			s.peers[peer] = true
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

	slog.Info("new peer connected", "remoteAddr", conn.RemoteAddr().String())

	if err := peer.read(); err != nil {
		slog.Error("read error", "err", err, "remoteAddr", conn.RemoteAddr().String())
	}
}

func main() {
	server := NewServer(Config{})
	log.Fatal(server.Start())
}
