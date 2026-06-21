package godot

import (
	"log"
	"sync"
)

// GDExtensionBridge handles the mock connection between the Go ECS and Godot C++
type GDExtensionBridge struct {
	mu           sync.RWMutex
	isActive     bool
	nodeRegistry map[string]int
}

var Bridge = &GDExtensionBridge{
	nodeRegistry: make(map[string]int),
}

// Initialize prepares the Godot integration layer.
func (b *GDExtensionBridge) Initialize() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.isActive = true
	log.Println("Godot GDExtension Bridge: Initialized")
}

// SyncEntity syncs a Go ECS entity with a Godot 3D Node.
func (b *GDExtensionBridge) SyncEntity(entityID string, ecsID int, x, y float64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.isActive {
		return
	}

	// Track the node
	b.nodeRegistry[entityID] = ecsID

	// Mock CGO call to update transform
	// C.update_godot_node(C.CString(entityID), C.double(x), C.double(0), C.double(y))
}

// Shutdown safely closes the bridge.
func (b *GDExtensionBridge) Shutdown() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.isActive = false
	b.nodeRegistry = make(map[string]int)
	log.Println("Godot GDExtension Bridge: Shut Down")
}
