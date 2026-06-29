# CHANGELOG
## [0.3.0] - 2024-06-03
- **EXECUTIVE SYNC**: Completed dual-direction merge, bringing all feature branches (Phase 13, Godot Bridge, Submodule fixes) into `main`.
- Incremented version heavily to reflect structural consolidation.


## [0.2.3] - 2024-06-03
- **BUGFIX**: Fixed Style physics compounding bug. The `Damping` physics variable is now idempotently reset to 1.0 before being multiplied by the style modifiers, ensuring safe hot-reloading.
- Cleaned up duplicate style application calls in `main.go` and `pkg/engine/patcher.go` to guarantee order-of-operations.


## [0.2.2] - 2024-06-03
- **BUGFIX**: Fixed several compilation errors in `main.go` caused by the recent intelligent merge that re-introduced unused imports and mismatched pointers to `schema.WorldConfig`.


## [0.2.1] - 2024-06-03
- **BUGFIX**: Fixed a build break in `main.go` caused by the recent ECS style trigger refactor (passed correct registry pointer to `ApplyStyle`).


## [0.2.0] - 2024-06-03
- Implemented Phase 3: Style-as-Technology Trigger Engine.
- Styles defined in JSON (`style_keywords`) now dynamically mutate the global ECS simulation (e.g. 'Souls Combat' increases global gravity and damping).
- Styles now dynamically inject suffix prompts into the asynchronous DALL-E asset pipeline to enforce art-direction uniformity (e.g. 'Cyberpunk' appends synthwave/neon tags to all room generation prompts).


## [0.1.0] - 2024-06-03
- **MAJOR SYNC**: Merged Godot CGO Bridge stubs, Advanced Rule Engine AST features, and concurrent ECS system fixes back into the mainline architecture.
- Project structure officially encompasses C++ bindings, OpenAI hooks, Ebitengine pipelines, and dynamic ECS memory.


## [0.0.19] - 2024-06-03
- **BUGFIX**: Resolved guaranteed collision deadlock in `pkg/ecs/systems.go`.
- De-coupled the physics lock from the collision callback invocation loop, preventing recursive lock acquisitions when executing dynamic rules like `ApplyDamage` or `SetFlag`.


## [0.0.19] - 2024-06-03
- **BUGFIX**: Resolved guaranteed collision deadlock in `pkg/ecs/systems.go`.
- De-coupled the physics lock from the collision callback invocation loop, preventing recursive lock acquisitions when executing dynamic rules like `ApplyDamage` or `SetFlag`.


## [0.0.18] - 2024-06-03
- Implemented Phase 14: Godot GDExtension CGO Integration.
- Converted `pkg/godot/bridge.go` to use actual CGO bindings (`import "C"`).
- Added C-level stubs to test memory pointer conversion (`C.CString`) and invocation across the FFI boundary without panicking the Go runtime.


## [0.0.17] - 2024-06-03
- Implemented Phase 13: Advanced Rule Engine and Scripting.
- Added `EntityState` to ECS to store arbitrary flags and numbers dynamically via script.
- Upgraded the pseudo-AST parser to handle nested AND conditions.
- Added `SetFlag` and `Heal` script commands.


## [0.0.16] - 2024-06-03
- Implemented Godot GDExtension bridge framework in `pkg/godot/bridge.go`.
- Provides the thread-safe mock integration layer needed to sync the deterministic Go ECS state with 3D Godot nodes.


## [0.0.15] - 2024-06-03
- Integrated actual OpenAI GPT-3.5 API for Natural Language Parsing (`pkg/engine/nlp_parser.go`).
- Integrated actual OpenAI DALL-E API for async asset generation (`pkg/assets/generator.go`).
- Maintained fallback mock logic to ensure stability when `OPENAI_API_KEY` is absent.


## [0.0.14] - 2024-06-02
- Implemented Natural Language to JSON API mock (`pkg/engine/nlp_parser.go`).
- Engine can now watch text files (`prompt.txt`) and dynamically hot-swap/regenerate running ECS entities using string heuristic commands (e.g., "generate a massive sci-fi dungeon").

## [0.0.13] - 2024-06-02
- Implemented AI-driven procedural level generator `pkg/engine/level_generator.go`.
- Added support for `LevelSeed`, `RoomCount`, and `Biome` configuration in `schema.go`.
- Wired generator to the hot-reload pipeline; the level map will now regenerate dynamically when JSON seed variables are modified.
- Integrated generative asset streaming: room generation automatically tags sprites with text prompts based on Biome selection, calling the asynchronous API mock.

## [0.0.12] - 2024-06-02
- Implemented complex rule parsing engine (`pkg/engine/rules.go`).
- Added a basic pseudo-AST to evaluate entity conditionals defined entirely via string parameters in JSON schemas (e.g., `Health < 50 AND COLLIDES_WITH`).
- Updated `UpdateCollision` systems to invoke a `CollisionCallback`, avoiding circular dependency while evaluating high-level rules dynamically.
- Included robust testing suite for rule string evaluation.

## [0.0.11] - 2024-06-02
- Implemented Asynchronous Asset Generation API mock (`pkg/assets`).
- Updated `SpriteRenderer` component to include a `Prompt` string for generative assets.
- Connected JSON Patcher to trigger async generation for new prompt-based sprites.
- Created thread-safe `TextureCache` for Ebitengine to dynamically fetch textures post-generation.

## [0.0.10] - 2024-06-02
- Implemented multi-threaded system updates using goroutines and sync.WaitGroup.
- Parallelized UpdatePhysics (excluding grid insert), UpdateCombat, and UpdateBehavior.

## [0.0.9] - 2024-06-02
- Integrated Ebitengine for 2D visual output.
- Replaced dummy console loop with a functional 60 TPS Ebitengine game loop.
- Added basic entity rendering using filled rectangles in Ebitengine.
- Fixed mutex lock issues in ECS to allow safe concurrent rendering.

## [0.0.8] - 2024-06-02
- Implemented Collision Filtering via `Layer` and `Mask` bitmasks.
- Added `IsTrigger` support to skip collision resolution while keeping event triggers.
- Implemented global velocity `Damping` in the physics system.
- Enhanced JSON Patcher to support advanced physics properties.

## [0.0.7] - 2024-06-02
- Implemented full-loop Integration Test Suite for the game engine.
- Enhanced Physics test suite with diagonal collision resolution coverage.
- Validated Scene Management and Registry clearing during transitions.

## [0.0.6] - 2024-06-02
- Improved Spatial Grid with neighbor-cell checking for robust broad-phase collisions.
- Implemented AABB Collision Resolution (separation and velocity reflection).
- Added `Static` and `Restitution` properties to `Collider` component.
- Implemented high-performance Raycasting system using the spatial grid.
- Expanded Rule Engine with "Stop" action.

## [0.0.5] - 2024-06-02
- Enhanced Physics layer with global Gravity support.
- Added comprehensive unit tests for the Physics simulation.
- Implemented Scene Management (Loading, Unloading, Transitions).
- Expanded Registry and Patcher to support gravity sync.

## [0.0.4] - 2024-06-02
- Integrated Godot Engine, Ebitengine, and Godot-cpp as submodules for future rendering and GDExtension support.
- Organized third-party dependencies under `third_party/`.

## [0.0.3] - 2024-06-02
- Implemented Grid-based Spatial Partitioning for optimized broad-phase collision detection.
- Added Souls-like Combat State Machine (Startup, Active, Recovery frames).
- Optimized Registry and Systems to use the spatial grid.
- Expanded patcher to support CombatState components.

## [0.0.2] - 2024-06-02
- Refactored ECS to use contiguous memory (slices) for improved performance.
- Implemented Health component and Damage action in the rule engine.
- Added mock Asset Generation interface.
- Updated systems to support new slice-based architecture.

## [0.0.1] - 2024-06-02
- Initial project structure.
- Basic ECS and JSON schema definition.

## [0.4.0] - 2026-06-29
### Added
- Created `pkg/net` package containing base `Server` and `Client` structs for UDP networking synchronization (Phase 15).
### Changed
- Added Godot CGO bridge entity spawn/despawn hooks and updated `Registry.DestroyEntity` behavior to facilitate better integration between the Go ECS and Godot C++ nodes.
- Expanded `pkg/net` Server and Client to serialize and synchronize ECS Position and Sprite state payloads using JSON over UDP.
