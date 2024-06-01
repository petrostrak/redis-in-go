package main

import (
	"context"
	"fmt"
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

type Message struct {
	data []byte
	peer *Peer
}

type Server struct {
	Config
	ln          net.Listener
	peers       map[*Peer]bool
	addPeerChan chan *Peer
	quit        chan struct{}
	msgChan     chan Message
	kv          *KV
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
		msgChan:     make(chan Message),
		kv:          NewKV(),
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
		case msg := <-s.msgChan:
			if err := s.handleMsg(msg); err != nil {
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

	if err := peer.read(); err != nil {
		slog.Error("read error", "err", err, "remoteAddr", conn.RemoteAddr().String())
	}
}

func (s *Server) handleMsg(msg Message) error {
	cmd, err := parseCommand(string(msg.data))
	if err != nil {
		return err
	}

	switch v := cmd.(type) {
	case SetCommand:
		return s.kv.Set(v.key, v.val)
	case GetCommand:
		value, found := s.kv.Get(string(v.key))
		if !found {
			return fmt.Errorf("key not found")
		}
		_, err := msg.peer.Send(value)
		if err != nil {
			slog.Error("peer send error", "err", err)
		}
	}

	return nil
}

func main() {
	server := NewServer(Config{})

	go func() {
		log.Fatal(server.Start())
	}()
	time.Sleep(time.Second)

	client := client.New("localhost:5001")

	if err := client.Set(context.Background(), "petros", "trakadas"); err != nil {
		slog.Error("error calling SET", "err", err)
	}

	time.Sleep(time.Second)
	fmt.Println(server.kv.data)
}
