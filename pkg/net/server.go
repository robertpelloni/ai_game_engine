package net

import (
	"fmt"
	"net"
	"sync"
	"encoding/json"

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

type NetworkState struct {
	Entities []EntityData `json:"entities"`
}

type EntityData struct {
	ID int `json:"id"`
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
	SpriteID string `json:"sprite_id"`
}

func (s *Server) BroadcastState() {
	if s.conn == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	state := NetworkState{
		Entities: make([]EntityData, 0),
	}

	s.registry.Mu.RLock()
	for i := 1; i < len(s.registry.HasPosition); i++ {
		if s.registry.HasPosition[i] {
			data := EntityData{
				ID: i,
				X:  s.registry.Positions[i].X,
				Y:  s.registry.Positions[i].Y,
			}
			if int(i) < len(s.registry.HasSprite) && s.registry.HasSprite[i] {
				data.SpriteID = s.registry.Sprites[i].SpriteID
			}
			state.Entities = append(state.Entities, data)
		}
	}
	s.registry.Mu.RUnlock()

	stateBytes, err := json.Marshal(state)
	if err != nil {
		fmt.Printf("Error marshaling state: %v\n", err)
		return
	}

	for _, addr := range s.clients {
		_, err := s.conn.WriteToUDP(stateBytes, addr)
		if err != nil {
			fmt.Printf("Error broadcasting to %s: %v\n", addr.String(), err)
		}
	}
}
