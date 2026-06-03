package scene

import (
	"fmt"
	"github.com/robertpelloni/ai_game_engine/pkg/ecs"
	"github.com/robertpelloni/ai_game_engine/pkg/engine"
	"log"
)

type SceneManager struct {
	Registry     *ecs.Registry
	CurrentScene *Scene
}

func NewSceneManager(registry *ecs.Registry) *SceneManager {
	return &SceneManager{
		Registry: registry,
	}
}

func (sm *SceneManager) LoadScene(scene *Scene) {
	fmt.Printf("Loading scene: %s\n", scene.ID)

	// 1. Unload current scene (clear registry)
	sm.UnloadScene()

	// 2. Patch registry with new scene schema
	engine.PatchRegistry(sm.Registry, scene.Schema)

	sm.CurrentScene = scene
}

func (sm *SceneManager) UnloadScene() {
	if sm.CurrentScene != nil {
		fmt.Printf("Unloading scene: %s\n", sm.CurrentScene.ID)
	}

	// Reset the registry to a clean state
	// In a real implementation, we'd preserve some global entities or state.
	*sm.Registry = *ecs.NewRegistry()
}

func (sm *SceneManager) TransitionTo(scene *Scene) {
	log.Printf("Transitioning from %v to %s", sm.CurrentScene, scene.ID)
	sm.LoadScene(scene)
}
