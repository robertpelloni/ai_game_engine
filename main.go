package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
	"image/color"

	"github.com/fsnotify/fsnotify"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/robertpelloni/ai_game_engine/pkg/ecs"
	"github.com/robertpelloni/ai_game_engine/pkg/engine"
	"github.com/robertpelloni/ai_game_engine/pkg/schema"
	"github.com/robertpelloni/ai_game_engine/pkg/assets"
)

type Game struct {
	registry *ecs.Registry
	schemaMu sync.RWMutex
	schema   *schema.GameSchema
}

func (g *Game) Update() error {
	// Rebuild ECS Registry dynamically if we parse natural language prompts
	select {
	case prompt := <-promptQueue:
		g.regenerateFromPrompt(prompt)
	default:
	}

	g.schemaMu.RLock()
	defer g.schemaMu.RUnlock()

	// Multi-threaded logic chunks
	g.registry.UpdatePhysics(1.0 / 60.0)
	g.registry.UpdateCombat()
	g.registry.UpdateBehavior()

	// Physics triggers rules
	g.registry.UpdateCollision(g.schema.Rules)

	return nil
}

func (g *Game) regenerateFromPrompt(prompt string) {
	log.Printf("Executing NLP schema regeneration: %s", prompt)
	s, err := engine.ParseNaturalLanguageToSchema(prompt)
	if err != nil {
		log.Printf("NLP failed: %v", err)
		return
	}

	g.schemaMu.Lock()
	defer g.schemaMu.Unlock()
	g.schema = s

	// Reset standard registry
	g.registry = ecs.NewRegistry()

	// Rebind the Collision Callback interface
	g.registry.CollisionCallback = func(e1, e2 ecs.Entity, rules []schema.EventAction) {
		for _, rule := range rules {
			if engine.ParseRuleCondition(g.registry, e1, e2, rule.Trigger) {
				engine.ExecuteRuleAction(g.registry, e1, e2, rule.Action)
			}
		}
	}

	// Rehydrate with base entities, then procedural generation, then styling overrides
	engine.PatchRegistry(g.registry, g.schema)
	engine.GenerateLevel(g.registry, &g.schema.World)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.schemaMu.RLock()
	defer g.schemaMu.RUnlock()

	g.registry.Mu.RLock()
	defer g.registry.Mu.RUnlock()

	screen.Fill(color.RGBA{20, 20, 30, 255})

	for i := 1; i < len(g.registry.HasSprite); i++ {
		if g.registry.HasSprite[i] && g.registry.HasPosition[i] {
			p := g.registry.Positions[i]
			sprite := g.registry.Sprites[i]

			img := assets.GetTexture(sprite.SpriteID)

			if img != nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(p.X, p.Y)
				screen.DrawImage(img, op)
			} else {
				// Fallback rectangle if loading
				c := g.registry.Colliders[i]
				w, h := 32.0, 32.0
				if g.registry.HasCollider[i] {
					w, h = c.Width, c.Height
				}
				ebitenutil.DrawRect(screen, p.X, p.Y, w, h, color.RGBA{255, 0, 0, 255})
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}

var promptQueue = make(chan string, 5)

func main() {
	schemaPath := filepath.Join(".", "pkg", "schema", "schema.json")
	bytes, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("Failed to read schema file: %v", err)
	}

	initialSchema, err := schema.ParseSchemaBytes(bytes)
	if err != nil {
		log.Fatalf("Failed to parse schema JSON: %v", err)
	}

	registry := ecs.NewRegistry()

	// Bind the Collision Callback interface
	registry.CollisionCallback = func(e1, e2 ecs.Entity, rules []schema.EventAction) {
		for _, rule := range rules {
			if engine.ParseRuleCondition(registry, e1, e2, rule.Trigger) {
				engine.ExecuteRuleAction(registry, e1, e2, rule.Action)
			}
		}
	}
	engine.PatchRegistry(registry, initialSchema)

	game := &Game{
		registry: registry,
		schema:   initialSchema,
	}

	go watchSchemaFile(schemaPath, game)
	go watchPromptFile("prompt.txt")

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Jules AI Generative Engine")

	log.Println("Starting game loop...")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Game execution failed: %v", err)
	}
}

func watchSchemaFile(path string, game *Game) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create file watcher: %v", err)
	}
	defer watcher.Close()

	if err := watcher.Add(path); err != nil {
		log.Fatalf("Failed to add file to watcher: %v", err)
	}

	var timer *time.Timer

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(100*time.Millisecond, func() {
					reloadSchema(path, game)
				})
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Watcher error: %v", err)
		}
	}
}

func reloadSchema(path string, game *Game) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read reloaded schema: %v", err)
		return
	}

	s, err := schema.ParseSchemaBytes(bytes)
	if err != nil {
		log.Printf("Failed to parse reloaded schema: %v", err)
		return
	}

	// Lock the global game state to safely update the schema
	game.schemaMu.Lock()
	defer game.schemaMu.Unlock()
	game.schema = s

	// Safely patch the running registry with the new components
	engine.PatchRegistry(game.registry, s)

	log.Println("Schema successfully hot-reloaded and registry patched!")
}

func watchPromptFile(path string) {
	// Create an empty prompt.txt file if it doesn't exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.WriteFile(path, []byte(""), 0644)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Failed to create file watcher: %v", err)
	}
	defer watcher.Close()

	if err := watcher.Add(path); err != nil {
		log.Fatalf("Failed to add file to watcher: %v", err)
	}

	var timer *time.Timer

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(500*time.Millisecond, func() {
					bytes, err := os.ReadFile(path)
					if err == nil && len(bytes) > 0 {
						prompt := string(bytes)
						select {
						case promptQueue <- prompt:
							// clear it out
							os.WriteFile(path, []byte(""), 0644)
						default:
							log.Println("Prompt queue is full, dropping prompt")
						}
					}
				})
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Prompt Watcher error: %v", err)
		}
	}
}
