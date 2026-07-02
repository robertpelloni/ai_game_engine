# Project Architecture & Patterns: AI Game Engine

## Core Philosophy and Design
*   **Generative & Data-Driven:** The core directive is to build an engine that creates deterministic, playable games entirely out of natural language processing and dynamically hot-swappable JSON files.
*   **Continuous Autonomous Execution:** The environment is structured to be constantly updated by AI tools operating semi-autonomously, with strict documentation standard handling (updating `TODO.md`, `ROADMAP.md`, `CHANGELOG.md`, `VERSION.md`, etc., continuously).
*   **Language:** Written primarily in **Go** (Golang), chosen for strong concurrency handling and performance.

## The Architecture: Strict Entity-Component-System (ECS)
The engine utilizes a strict, highly optimized custom ECS built within `pkg/ecs/`.

*   **Entities:** Represented purely as generic numerical IDs (`uint32` or `int`).
*   **Components:** Rather than allocating objects or structs, components are packed together into strictly contiguous slice arrays (e.g., `registry.Positions`, `registry.Velocities`). This design choice dramatically increases CPU cache-coherency and runtime performance. A parallel set of boolean slices (`registry.HasPosition`) acts as bitmasks to verify component existence.
*   **Systems:** Located in `systems.go`. Logic operates on arrays sequentially. We have recently **parallelized** these systems:
    *   **Physics:** Applies Gravity, Damping, and integrates Velocity into Position, processing entity chunks concurrently using goroutines.
    *   **Combat State Machine:** Simulates deep, Souls-like combat states natively, also processed concurrently.
    *   **Behavior:** Processes AI behavior chunks in parallel.

## Concurrency & Thread Safety
*   To support the multi-threaded system updates and concurrent Ebitengine rendering loop, the `ecs.Registry` struct exposes a public `Mu sync.RWMutex`.
*   Systems that mutate data arrays lock the registry (`r.Mu.Lock()`), while systems that only read data (like the `Draw` loop in `main.go`) use an RLock (`r.Mu.RLock()`).
*   Parallel chunks are synchronized using `sync.WaitGroup`, preventing data races.

## Collision & Physics Module
*   **Spatial Partitioning:** To avoid O(N²) collision checks, the engine uses a custom Grid-based Spatial Partitioning algorithm (`pkg/ecs/grid.go`). Because grid insertion is not thread-safe, it is executed sequentially *after* the parallelized position updates.
*   **AABB Resolution:** Standard Axis-Aligned Bounding Box (AABB) physics are used.
*   **Triggers & Filtering:** Supports non-colliding `IsTrigger` hitboxes and uses bitmask logic (`Layer` vs `Mask`) for collision filtering.

## State Patching & Hot-Reloading
*   **JSON Schema:** The "source of truth" is defined in `pkg/schema`. A JSON file dictates world configuration, entity components, rules, and global styles.
*   **Patcher (`pkg/engine/patcher.go`):** Translates untyped JSON maps into the ECS registry.
*   **File Watching:** Dynamically watches `schema.json` via background goroutines, dynamically hot-reloading data while safely locking the ECS registry to prevent data races.

## Rendering Engine
*   **Ebitengine Integration:** Uses `github.com/hajimehoshi/ebiten/v2` for 2D visual outputs running a 60 TPS loop in `main.go`.

## Asset Generation Integration
*   The system implements an asynchronous generator API framework (`pkg/assets/generator.go`). The `SpriteRenderer` component supports a `Prompt` string, enabling the dynamic generation of actual texture sprites dynamically at runtime when entities are loaded via the JSON patcher.
*   A global `TextureCache` guarded by `sync.RWMutex` prevents duplicate generation triggers for identical overlapping requests while drawing dynamically in the Ebitengine loop.

## Dynamic Rule Parser
*   The engine evaluates complex scripted behaviors and game logic using a pseudo-AST string parser in `pkg/engine/rules.go`. To prevent circular dependencies between the high-level engine and the low-level `pkg/ecs` packages, the ECS Registry exposes a `CollisionCallback` used to inject and evaluate these rules dynamically (e.g. `Health < 50 AND COLLIDES_WITH`) directly from text parameters defined in the JSON schema.

## Submodules
*   The project tracks `third_party/godot` and `third_party/godot-cpp` as git submodules to bridge complex 3D integration in future iterations.