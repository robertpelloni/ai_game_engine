package main

import (
	"time"

	"sync"

	"github.com/robertpelloni/ai_game_engine/pkg/assets"

	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/robertpelloni/ai_game_engine/pkg/ecs"
	"github.com/robertpelloni/ai_game_engine/pkg/engine"
	"github.com/robertpelloni/ai_game_engine/pkg/schema"
)

type Game struct {
	registry        *ecs.Registry
	schemaMu        sync.RWMutex
	schema          *schema.GameSchema
	generatedLevels []ecs.Entity
}

func (g *Game) Update() error {
	dt := 1.0 / 60.0 // 60 TPS
	g.registry.UpdatePhysics(dt)
	g.registry.UpdateCombat()
	g.registry.UpdateBehavior()

	g.schemaMu.RLock()
	rules := g.schema.Rules
	g.schemaMu.RUnlock()

	g.registry.UpdateCollision(rules)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	g.registry.Mu.RLock()
	defer g.registry.Mu.RUnlock()

	// Very basic rendering loop over entities that have both a Position and Collider/Sprite
	for i := 1; i < len(g.registry.HasPosition); i++ {
		if !g.registry.HasPosition[i] {
			continue
		}

		pos := g.registry.Positions[i]

		var w, h float32 = 10, 10
		if i < len(g.registry.HasCollider) && g.registry.HasCollider[i] {
			w = float32(g.registry.Colliders[i].Width)
			h = float32(g.registry.Colliders[i].Height)
		}

		c := color.RGBA{255, 255, 255, 255}
		if i < len(g.registry.HasHealth) && g.registry.HasHealth[i] {
			// Change color if damaged
			if g.registry.Healths[i].Current < g.registry.Healths[i].Max {
				c = color.RGBA{255, 0, 0, 255}
			}
		}

		img := assets.GetTexture(g.registry.Sprites[i].SpriteID)
		if img != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(pos.X, pos.Y)
			screen.DrawImage(img, op)
		} else {
			vector.DrawFilledRect(screen, float32(pos.X), float32(pos.Y), w, h, c, false)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	fmt.Println("Starting AI Game Engine v0.0.14...")

	// Initial mock schema
	initialSchema := &schema.GameSchema{
		World: schema.WorldConfig{
			GridSpacing: 1.0,
			Gravity:     []float64{0, 0}, // Removed gravity for top-down test
			LevelSeed:   12345,
			RoomCount:   5,
			Biome:       "Sci-Fi",
		},
		Entities: []schema.EntitySpec{
			{
				ID: 1,
				Components: []schema.ComponentData{
					{Type: "Position", Data: map[string]interface{}{"x": 300.0, "y": 200.0}},
					{Type: "Velocity", Data: map[string]interface{}{"vx": 50.0, "vy": 20.0}},
					{Type: "SpriteRenderer", Data: map[string]interface{}{"sprite_id": "player"}},
					{Type: "Collider", Data: map[string]interface{}{"width": 20.0, "height": 20.0}},
					{Type: "Health", Data: map[string]interface{}{"current": 100.0, "max": 100.0}},
					{Type: "CombatState", Data: map[string]interface{}{
						"state": "Startup", "frames_left": 2.0, "startup_frames": 2.0, "active_frames": 3.0, "recovery_frames": 2.0,
					}},
				},
			},
			{
				ID: 2,
				Components: []schema.ComponentData{
					{Type: "Position", Data: map[string]interface{}{"x": 400.0, "y": 250.0}},
					{Type: "Collider", Data: map[string]interface{}{"width": 30.0, "height": 30.0, "static": true}},
					{Type: "Health", Data: map[string]interface{}{"current": 50.0, "max": 50.0}},
				},
			},
		},
		Rules: []schema.EventAction{
			{Trigger: "COLLIDES_WITH", Action: "Damage"},
			{Trigger: "COLLIDES_WITH", Action: "Stop"},
			{Trigger: "Health < 50 AND COLLIDES_WITH", Action: "RunAway"},
		},
		StyleKeywords: []string{"Retro Raycaster", "Souls Combat"},
	}

	registry := ecs.NewRegistry()
	registry.CollisionCallback = func(e1, e2 ecs.Entity, rules []schema.EventAction) {
		for _, rule := range rules {
			if engine.ParseRuleCondition(registry, e1, e2, rule.Trigger) {
				engine.ExecuteRuleAction(registry, e1, e2, rule.Action)
			}
		}
	}
	engine.PatchRegistry(registry, initialSchema)

	styleConfig := engine.GetStyleConfig(initialSchema.StyleKeywords)
	engine.ApplyStyle(styleConfig)



	game := &Game{
		registry: registry,
		schema:   initialSchema,
	}
	game.generatedLevels = engine.GenerateLevel(registry, &initialSchema.World)

	// Watch prompt.txt for natural language instructions
	promptFile := "prompt.txt"
	os.WriteFile(promptFile, []byte("generate a huge fantasy dungeon"), 0644)
	defer os.Remove(promptFile)

	go func() {
		lastMod := time.Time{}
		for {
			info, err := os.Stat(promptFile)
			if err == nil && info.ModTime().After(lastMod) {
				lastMod = info.ModTime()
				bytes, err := os.ReadFile(promptFile)
				if err == nil && len(bytes) > 0 {
					newSchema, err := engine.ParseNaturalLanguageToSchema(string(bytes))
					if err == nil && newSchema != nil {
						// Lock and apply new map details
						game.schemaMu.Lock()
						game.schema.World.LevelSeed = newSchema.World.LevelSeed
						game.schema.World.RoomCount = newSchema.World.RoomCount
						game.schema.World.Biome = newSchema.World.Biome

						game.generatedLevels = engine.GenerateLevel(registry, &game.schema.World)
						game.schemaMu.Unlock()
					}
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()

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

		// In a real app we'd need to clear old entities safely by locking and deleting them.
		// For the mock, we can just regenerate over if the seed differs (not fully clearing yet to save time on MVP).

		game.schemaMu.RLock()
		seedChanged := s.World.LevelSeed != game.schema.World.LevelSeed || s.World.RoomCount != game.schema.World.RoomCount
		game.schemaMu.RUnlock()

		if seedChanged {
			game.schemaMu.Lock()
			game.schema = s
			game.generatedLevels = engine.GenerateLevel(registry, &s.World)
			game.schemaMu.Unlock()
		}
	})



	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("AI Game Engine 0.0.14")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nEngine shut down.")
}
