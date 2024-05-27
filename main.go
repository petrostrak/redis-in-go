package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"time"

	"github.com/petrostrak/redis-in-go/client"
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
	msgChan     chan []byte
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
		msgChan:     make(chan []byte),
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
		case rawMsg := <-s.msgChan:
			if err := s.handleRawMsg(rawMsg); err != nil {
				slog.Info("raw msg error", "err", err)
			}
		case <-s.quit:
			return
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
	peer := NewPeer(conn, s.msgChan)
	s.addPeerChan <- peer

	slog.Info("new peer connected", "remoteAddr", conn.RemoteAddr().String())

	if err := peer.read(); err != nil {
		slog.Error("read error", "err", err, "remoteAddr", conn.RemoteAddr().String())
	}
}

func (s *Server) handleRawMsg(msg []byte) error {
	cmd, err := parseCommand(string(msg))
	if err != nil {
		return err
	}

	switch v := cmd.(type) {
	case SetCommand:
		slog.Info("SET", "key", v.key, "value", v.val)
	}

	return nil
}

func main() {
	go func() {
		server := NewServer(Config{})
		log.Fatal(server.Start())
	}()
	time.Sleep(time.Second)

	client := client.New("localhost:5001")

	if err := client.Set(context.Background(), "petros", "trakadas"); err != nil {
		slog.Error("error calling SET", "err", err)
	}

	select {}
}
