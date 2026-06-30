# Godot Integration Design

The Godot Integration relies on `pkg/godot/bridge.go` to provide a CGO layer between the deterministic Go ECS simulation and the 3D visual frontend of Godot.

## Synchronization
The synchronization is managed by `GDExtensionBridge`. The Go ECS owns the simulation (physics, rules, positions).
`HookRegistryCreateEntity` and `HookRegistryDestroyEntity` broadcast spawn/despawn events over the C++ boundary.
Godot simply acts as a dumb renderer and physics visualizer.

## Thread Safety
The CGO bridge heavily relies on Go's `sync.RWMutex` to prevent concurrent modification panics between the Go game loop (running at 60 TPS) and the Godot rendering thread (running asynchronously).
