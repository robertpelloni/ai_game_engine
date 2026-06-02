# CHANGELOG
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
