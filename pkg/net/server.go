package net

import (
	"fmt"
	"net"
	"sync"

	"github.com/robertpelloni/ai_game_engine/pkg/ecs"
)

type Server struct {
	mu          sync.Mutex
	clients     map[string]*net.UDPAddr
	registry    *ecs.Registry
	conn        *net.UDPConn
}

func NewServer(reg *ecs.Registry) *Server {
	return &Server{
		clients:  make(map[string]*net.UDPAddr),
		registry: reg,
	}
}

func (s *Server) Start(port int) error {
	addr := &net.UDPAddr{
		Port: port,
		IP:   net.ParseIP("0.0.0.0"),
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	s.conn = conn

	fmt.Printf("Started UDP Server on port %d\n", port)

	go s.listen()
	return nil
}

func (s *Server) listen() {
	buf := make([]byte, 1024)
	for {
		n, addr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("Error reading from UDP: %v\n", err)
			continue
		}

		s.mu.Lock()
		clientKey := addr.String()
		if _, exists := s.clients[clientKey]; !exists {
			fmt.Printf("New client connected: %s\n", clientKey)
			s.clients[clientKey] = addr
		}
		s.mu.Unlock()

		// process client data if necessary
		_ = n
	}
}

func (s *Server) BroadcastState() {
	if s.conn == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Extremely simplified serialization mock
	stateMsg := []byte("STATE_UPDATE")

	for _, addr := range s.clients {
		_, err := s.conn.WriteToUDP(stateMsg, addr)
		if err != nil {
			fmt.Printf("Error broadcasting to %s: %v\n", addr.String(), err)
		}
	}
}
