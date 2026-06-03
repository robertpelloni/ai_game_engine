package engine_test

import (
	"github.com/robertpelloni/ai_game_engine/pkg/ecs"
	"github.com/robertpelloni/ai_game_engine/pkg/schema"
	"github.com/robertpelloni/ai_game_engine/pkg/scene"
	"testing"
)

func TestIntegrationGameLoop(t *testing.T) {
	// 1. Initialize Registry and Manager
	reg := ecs.NewRegistry()
	sm := scene.NewSceneManager(reg)

	// 2. Define initial scene
	s1 := &schema.GameSchema{
		World: schema.WorldConfig{Gravity: []float64{0, -9.81}},
		Entities: []schema.EntitySpec{
			{
				ID: 1,
				Components: []schema.ComponentData{
					{Type: "Position", Data: map[string]interface{}{"x": 0.0, "y": 100.0}},
					{Type: "Velocity", Data: map[string]interface{}{"vx": 0.0, "vy": 0.0}},
					{Type: "Collider", Data: map[string]interface{}{"width": 10.0, "height": 10.0}},
				},
			},
		},
	}
	scene1 := scene.NewScene("Level1", s1)

	// 3. Load scene
	sm.LoadScene(scene1)

	// 4. Run loop
	dt := 0.016
	reg.UpdatePhysics(dt)
	reg.UpdateCollision(nil)

	if reg.Positions[1].Y >= 100.0 {
		t.Errorf("Expected gravity to move entity 1 down, got Y=%f", reg.Positions[1].Y)
	}

	// 5. Simulate Scene Transition
	s2 := &schema.GameSchema{
		Entities: []schema.EntitySpec{
			{
				ID: 2,
				Components: []schema.ComponentData{
					{Type: "Position", Data: map[string]interface{}{"x": 50.0, "y": 50.0}},
				},
			},
		},
	}
	scene2 := scene.NewScene("Level2", s2)
	sm.TransitionTo(scene2)

	if !reg.HasPosition[2] {
		t.Error("Expected entity 2 to exist after transition")
	}
	if reg.HasPosition[1] {
		t.Error("Expected entity 1 to be cleared after transition")
	}
}
