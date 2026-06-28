package net

import (
	"fmt"
	"net"

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
	buf := make([]byte, 1024)
	for {
		n, _, err := c.conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("Error reading from server: %v\n", err)
			continue
		}

		// process state updates
		msg := string(buf[:n])
		if msg == "STATE_UPDATE" {
			// Stub for parsing and updating local registry
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
