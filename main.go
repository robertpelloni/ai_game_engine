package main

import (
	"fmt"
	"github.com/robertpelloni/ai_game_engine/pkg/ecs"
	"github.com/robertpelloni/ai_game_engine/pkg/engine"
	"github.com/robertpelloni/ai_game_engine/pkg/schema"
	"log"
	"os"
	"time"
)

func main() {
	fmt.Println("Starting AI Game Engine v0.0.8...")

	// Initial mock schema
	initialSchema := &schema.GameSchema{
		World: schema.WorldConfig{
			GridSpacing: 1.0,
			Gravity:     []float64{0, -9.81},
		},
		Entities: []schema.EntitySpec{
			{
				ID: 1,
				Components: []schema.ComponentData{
					{Type: "Position", Data: map[string]interface{}{"x": 0.0, "y": 0.0}},
					{Type: "Velocity", Data: map[string]interface{}{"vx": 100.0, "vy": 0.0}},
					{Type: "SpriteRenderer", Data: map[string]interface{}{"sprite_id": "player"}},
					{Type: "Collider", Data: map[string]interface{}{"width": 10.0, "height": 10.0}},
					{Type: "Health", Data: map[string]interface{}{"current": 100.0, "max": 100.0}},
					{Type: "CombatState", Data: map[string]interface{}{
						"state": "Startup", "frames_left": 2.0, "startup_frames": 2.0, "active_frames": 3.0, "recovery_frames": 2.0,
					}},
				},
			},
			{
				ID: 2,
				Components: []schema.ComponentData{
					{Type: "Position", Data: map[string]interface{}{"x": 5.0, "y": 0.0}},
					{Type: "Collider", Data: map[string]interface{}{"width": 10.0, "height": 10.0}},
					{Type: "Health", Data: map[string]interface{}{"current": 50.0, "max": 50.0}},
				},
			},
		},
		Rules: []schema.EventAction{
			{Trigger: "COLLIDES_WITH", Action: "Damage"},
		},
		StyleKeywords: []string{"Retro Raycaster", "Souls Combat"},
	}

	registry := ecs.NewRegistry()
	engine.PatchRegistry(registry, initialSchema)

	styleConfig := engine.GetStyleConfig(initialSchema.StyleKeywords)
	engine.ApplyStyle(styleConfig)

	// Create a dummy schema file for the watcher
	schemaFile := "schema.json"
	os.WriteFile(schemaFile, []byte(`{"style_keywords": ["Gritty Noir"]}`), 0644)
	defer os.Remove(schemaFile)

	// Run hot-reload watcher in background
	go engine.WatchSchema(schemaFile, func(data []byte) {
		s, err := schema.ParseSchemaBytes(data)
		if err != nil {
			log.Printf("Failed to parse reloaded schema: %v", err)
			return
		}
		engine.PatchRegistry(registry, s)
		newStyle := engine.GetStyleConfig(s.StyleKeywords)
		engine.ApplyStyle(newStyle)
	})

	// Main game loop
	dt := 0.016
	for i := 0; i < 10; i++ {
		fmt.Printf("\n--- Frame %d ---\n", i)
		registry.UpdatePhysics(dt)
		registry.UpdateCombat()
		registry.UpdateBehavior()
		registry.UpdateCollision(initialSchema.Rules)
		registry.UpdateRender()
		time.Sleep(16 * time.Millisecond)
	}

	fmt.Println("\nEngine shut down.")
}
