package engine

import (
	"github.com/robertpelloni/ai_game_engine/pkg/ecs"
	"github.com/robertpelloni/ai_game_engine/pkg/schema"
	"testing"
)

func TestStyleTriggers(t *testing.T) {
	config := GetStyleConfig([]string{"Retro Raycaster"})
	if !config.RaycastingEnabled {
		t.Error("Expected raycasting to be enabled")
	}
	if config.Resolution[0] != 320 {
		t.Errorf("Expected resolution width 320, got %d", config.Resolution[0])
	}
}

func TestPatchRegistry(t *testing.T) {
	reg := ecs.NewRegistry()
	s := &schema.GameSchema{
		Entities: []schema.EntitySpec{
			{
				ID: 1,
				Components: []schema.ComponentData{
					{Type: "Position", Data: map[string]interface{}{"x": 100.0, "y": 200.0}},
				},
			},
		},
	}

	PatchRegistry(reg, s)

	if p, ok := reg.Positions[ecs.Entity(1)]; !ok || p.X != 100 || p.Y != 200 {
		t.Errorf("Failed to patch position correctly")
	}
}
