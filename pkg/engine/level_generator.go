package engine

import (
	"github.com/robertpelloni/ai_game_engine/pkg/assets"

	"fmt"
	"math/rand"

	"github.com/robertpelloni/ai_game_engine/pkg/ecs"
	"github.com/robertpelloni/ai_game_engine/pkg/schema"
)

// GenerateLevel generates procedural entities into the registry based on world config.
// It tracks generated entities by adding a specific tag/flag component if possible,
// or simply tracks them to allow clearing them on hot-reload.
func GenerateLevel(reg *ecs.Registry, config *schema.WorldConfig) []ecs.Entity {
	if config.LevelSeed == 0 || config.RoomCount == 0 {
		return nil
	}

	fmt.Printf("Generating %d rooms for biome '%s' using seed %d...\n", config.RoomCount, config.Biome, config.LevelSeed)

	rng := rand.New(rand.NewSource(config.LevelSeed))
	generatedEntities := make([]ecs.Entity, 0)

	// Determine asset prompts based on Biome
	wallPrompt := "A basic grey wall block"
	floorPrompt := "A simple dirt floor texture"
	if config.Biome == "Sci-Fi" {
		wallPrompt = "A futuristic metallic bulkhead with glowing neon lights"
		floorPrompt = "A grated steel floor panel"
	} else if config.Biome == "Fantasy" {
		wallPrompt = "A mossy cobblestone dungeon wall"
		floorPrompt = "A cracked stone pathway"
	}

	// Generate a simple scatter of rooms/walls
	for i := 0; i < config.RoomCount; i++ {
		// Room center
		cx := rng.Float64() * 800
		cy := rng.Float64() * 600

		// Add a floor prop (no collider)
		floor := reg.CreateEntity()
		reg.AddPosition(floor, ecs.Position{X: cx, Y: cy})
		spriteID := fmt.Sprintf("floor_%d", i)
		reg.AddSprite(floor, ecs.SpriteRenderer{
			SpriteID: spriteID,
			Prompt:   floorPrompt,
		})
		assets.GenerateAssetAsync(spriteID, floorPrompt, 32, 32)
		generatedEntities = append(generatedEntities, floor)

		// Surround with some random walls
		for j := 0; j < 4; j++ {
			wx := cx + (rng.Float64()*100 - 50)
			wy := cy + (rng.Float64()*100 - 50)

			wall := reg.CreateEntity()
			reg.AddPosition(wall, ecs.Position{X: wx, Y: wy})
			reg.AddCollider(wall, ecs.Collider{Width: 32, Height: 32, Static: true, Layer: 1, Mask: 1})
			spriteID2 := fmt.Sprintf("wall_%d_%d", i, j)
			reg.AddSprite(wall, ecs.SpriteRenderer{
				SpriteID: spriteID2,
				Prompt:   wallPrompt,
			})
			assets.GenerateAssetAsync(spriteID2, wallPrompt, 32, 32)
			generatedEntities = append(generatedEntities, wall)
		}
	}

	return generatedEntities
}
