package godot

/*
#include <stdlib.h>
#include <stdio.h>

// Mock C function representing a Godot C++ method
void godot_update_node_transform(char* id, double x, double y) {
    // In a real GDExtension, this would pass the memory to Godot's StringName and Vector3 classes.
    // printf("CGO: Syncing Godot Node %s to Position(%f, %f)\n", id, x, y);
}

// Mock C function for spawning nodes
void godot_spawn_node(char* id) {
}

// Mock C function for despawning nodes
void godot_despawn_node(char* id) {
}
*/
import "C"
import (
	"log"
	"sync"
	"unsafe"
)

// GDExtensionBridge handles the connection between the Go ECS and Godot C++
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
	log.Println("Godot GDExtension Bridge: Initialized with CGO")
}

// SyncEntity syncs a Go ECS entity with a Godot 3D Node across the C boundary.
func (b *GDExtensionBridge) SyncEntity(entityID string, ecsID int, x, y float64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.isActive {
		return
	}

	// Track the node
	b.nodeRegistry[entityID] = ecsID

	// Convert Go string to C string
	cEntityID := C.CString(entityID)
	defer C.free(unsafe.Pointer(cEntityID))

	// Call the mock C layer
	C.godot_update_node_transform(cEntityID, C.double(x), C.double(y))
}

func (b *GDExtensionBridge) SpawnNode(entityID string, ecsID int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.isActive {
		return
	}
	b.nodeRegistry[entityID] = ecsID
	cEntityID := C.CString(entityID)
	defer C.free(unsafe.Pointer(cEntityID))
	C.godot_spawn_node(cEntityID)
}

func (b *GDExtensionBridge) DespawnNode(entityID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.isActive {
		return
	}
	delete(b.nodeRegistry, entityID)
	cEntityID := C.CString(entityID)
	defer C.free(unsafe.Pointer(cEntityID))
	C.godot_despawn_node(cEntityID)
}

// Shutdown safely closes the bridge.
func (b *GDExtensionBridge) Shutdown() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.isActive = false
	b.nodeRegistry = make(map[string]int)
	log.Println("Godot GDExtension Bridge: Shut Down")
}
