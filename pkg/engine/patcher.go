package engine

import (
	"github.com/robertpelloni/ai_game_engine/pkg/ecs"
	"github.com/robertpelloni/ai_game_engine/pkg/schema"
	"log"
)

func getFloat(data map[string]interface{}, key string) float64 {
	if val, ok := data[key]; ok {
		if f, ok := val.(float64); ok {
			return f
		}
	}
	return 0
}

func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

func PatchRegistry(registry *ecs.Registry, s *schema.GameSchema) {
	// Sync World Gravity
	if len(s.World.Gravity) >= 2 {
		registry.GravityX = s.World.Gravity[0]
		registry.GravityY = s.World.Gravity[1]
	}

	for _, entitySpec := range s.Entities {
		e := ecs.Entity(entitySpec.ID)

		for _, comp := range entitySpec.Components {
			data, ok := comp.Data.(map[string]interface{})
			if !ok {
				log.Printf("Invalid component data for entity %d type %s", entitySpec.ID, comp.Type)
				continue
			}

			switch comp.Type {
			case "Position":
				registry.AddPosition(e, ecs.Position{
					X: getFloat(data, "x"),
					Y: getFloat(data, "y"),
				})
			case "Velocity":
				registry.AddVelocity(e, ecs.Velocity{
					VX: getFloat(data, "vx"),
					VY: getFloat(data, "vy"),
				})
			case "SpriteRenderer":
				registry.AddSprite(e, ecs.SpriteRenderer{
					SpriteID: getString(data, "sprite_id"),
				})
			case "Collider":
				registry.AddCollider(e, ecs.Collider{
					Width:       getFloat(data, "width"),
					Height:      getFloat(data, "height"),
					Restitution: getFloat(data, "restitution"),
					Static:      data["static"] == true,
				})
			case "AIBehavior":
				registry.AddAIBehavior(e, ecs.AIBehavior{
					BehaviorType: getString(data, "behavior_type"),
				})
			case "Health":
				registry.AddHealth(e, ecs.Health{
					Current: getFloat(data, "current"),
					Max:     getFloat(data, "max"),
				})
			case "CombatState":
				registry.AddCombatState(e, ecs.CombatState{
					State:          getString(data, "state"),
					FramesLeft:     int(getFloat(data, "frames_left")),
					StartupFrames:  int(getFloat(data, "startup_frames")),
					ActiveFrames:   int(getFloat(data, "active_frames")),
					RecoveryFrames: int(getFloat(data, "recovery_frames")),
				})
			default:
				log.Printf("Unknown component type: %s", comp.Type)
			}
		}
	}
}
