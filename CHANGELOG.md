# CHANGELOG
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
