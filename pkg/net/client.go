package net

import (
	"fmt"
	"net"
	"encoding/json"

	"github.com/robertpelloni/ai_game_engine/pkg/ecs"
)

type Client struct {
	conn     *net.UDPConn
	registry *ecs.Registry
}

func NewClient(reg *ecs.Registry) *Client {
	return &Client{
		registry: reg,
	}
}

func (c *Client) Connect(serverAddr string) error {
	addr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return err
	}
	c.conn = conn

	fmt.Printf("Connected to server at %s\n", serverAddr)

	// Send initial greeting to register with server
	_, err = c.conn.Write([]byte("HELLO"))
	if err != nil {
		return err
	}

	go c.listen()
	return nil
}

func (c *Client) listen() {
	buf := make([]byte, 8192) // larger buffer for json
	for {
		n, _, err := c.conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("Error reading from server: %v\n", err)
			continue
		}

		var state NetworkState
		err = json.Unmarshal(buf[:n], &state)
		if err == nil {
			c.registry.Mu.Lock()
			for _, ed := range state.Entities {
				e := ecs.Entity(ed.ID)
				if int(e) >= len(c.registry.HasPosition) || !c.registry.HasPosition[e] {
					// Extremely basic creation stub for remote entities
					c.registry.AddPosition(e, ecs.Position{X: ed.X, Y: ed.Y})
					if ed.SpriteID != "" {
						c.registry.AddSprite(e, ecs.SpriteRenderer{SpriteID: ed.SpriteID})
					}
				} else {
					c.registry.Positions[e].X = ed.X
					c.registry.Positions[e].Y = ed.Y
				}
			}
			c.registry.Mu.Unlock()
		}
	}
}

func (c *Client) SendAction(actionData []byte) error {
	if c.conn == nil {
		return fmt.Errorf("client not connected")
	}
	_, err := c.conn.Write(actionData)
	return err
}
