package godot

import (
	"fmt"

	"github.com/robertpelloni/ai_game_engine/pkg/ecs"
)

func HookRegistryCreateEntity(reg *ecs.Registry) ecs.Entity {
	e := reg.CreateEntity()
	entityIDStr := fmt.Sprintf("entity_%d", e)
	Bridge.SpawnNode(entityIDStr, int(e))
	return e
}

func HookRegistryDestroyEntity(reg *ecs.Registry, e ecs.Entity) {
	entityIDStr := fmt.Sprintf("entity_%d", e)
	Bridge.DespawnNode(entityIDStr)
	reg.DestroyEntity(e)
}
