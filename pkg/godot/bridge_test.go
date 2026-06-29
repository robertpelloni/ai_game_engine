package godot

import "testing"

func TestBridgeInitialization(t *testing.T) {
	Bridge.Initialize()

	Bridge.mu.RLock()
	if !Bridge.isActive {
		t.Error("Bridge should be active after initialization")
	}
	Bridge.mu.RUnlock()

	Bridge.SpawnNode("player", 1)

	Bridge.mu.RLock()
	if Bridge.nodeRegistry["player"] != 1 {
		t.Error("Entity was not spawned to the registry correctly")
	}
	Bridge.mu.RUnlock()

	Bridge.SyncEntity("player", 1, 10.5, 20.0)

	Bridge.DespawnNode("player")
	Bridge.mu.RLock()
	if _, ok := Bridge.nodeRegistry["player"]; ok {
		t.Error("Entity was not despawned from the registry correctly")
	}
	Bridge.mu.RUnlock()

	Bridge.Shutdown()

	Bridge.mu.RLock()
	if Bridge.isActive {
		t.Error("Bridge should be inactive after shutdown")
	}
	Bridge.mu.RUnlock()
}
