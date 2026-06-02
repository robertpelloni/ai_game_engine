package engine

import (
	"github.com/robertpelloni/ai_game_engine/pkg/ecs"
	"github.com/robertpelloni/ai_game_engine/pkg/schema"
	"log"
)

func PatchRegistry(registry *ecs.Registry, s *schema.GameSchema) {
	// For simplicity, this "patcher" currently re-syncs the entire state.
	// A more advanced version would only apply deltas.

	for _, entitySpec := range s.Entities {
		e := ecs.Entity(entitySpec.ID)

		for _, comp := range entitySpec.Components {
			switch comp.Type {
			case "Position":
				data := comp.Data.(map[string]interface{})
				registry.AddPosition(e, &ecs.Position{
					X: data["x"].(float64),
					Y: data["y"].(float64),
				})
			case "Velocity":
				data := comp.Data.(map[string]interface{})
				registry.AddVelocity(e, &ecs.Velocity{
					VX: data["vx"].(float64),
					VY: data["vy"].(float64),
				})
			case "SpriteRenderer":
				data := comp.Data.(map[string]interface{})
				registry.AddSprite(e, &ecs.SpriteRenderer{
					SpriteID: data["sprite_id"].(string),
				})
			case "Collider":
				data := comp.Data.(map[string]interface{})
				registry.AddCollider(e, &ecs.Collider{
					Width:  data["width"].(float64),
					Height: data["height"].(float64),
				})
			case "AIBehavior":
				data := comp.Data.(map[string]interface{})
				registry.AddAIBehavior(e, &ecs.AIBehavior{
					BehaviorType: data["behavior_type"].(string),
				})
			default:
				log.Printf("Unknown component type: %s", comp.Type)
			}
		}
	}
}
