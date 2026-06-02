# Project Architecture: Generative AI Game Engine Runtime

## System Rules
- This project utilizes a strict Entity Component System (ECS) architecture.
- The state of the world, entities, and behaviors are entirely data-driven, defined by a unified JSON schema.
- The engine must read this JSON file dynamically at runtime and translate it into memory-allocated game objects.
- All code must support hot-swapping/hot-reloading JSON state deltas without crashing or losing core simulation loops.
