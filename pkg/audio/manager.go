package audio

import "fmt"

type Manager struct {
	isActive bool
}

func NewManager() *Manager {
	return &Manager{isActive: true}
}

func (m *Manager) PlaySound(soundID string) {
	if m.isActive {
		fmt.Printf("Playing sound: %s\n", soundID)
	}
}
