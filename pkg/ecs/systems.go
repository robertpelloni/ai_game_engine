package ecs

import (
	"fmt"
	"github.com/robertpelloni/ai_game_engine/pkg/schema"
	"strings"
)

func (r *Registry) UpdatePhysics(dt float64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.SpatialGrid.Clear()

	for i := 1; i < len(r.HasVelocity); i++ {
		if r.HasVelocity[i] && r.HasPosition[i] {
			r.Positions[i].X += r.Velocities[i].VX * dt
			r.Positions[i].Y += r.Velocities[i].VY * dt
		}
		if r.HasPosition[i] {
			r.SpatialGrid.Insert(Entity(i), r.Positions[i].X, r.Positions[i].Y)
		}
	}
}

func (r *Registry) UpdateBehavior() {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i := 1; i < len(r.HasAIBehavior); i++ {
		if r.HasAIBehavior[i] {
			fmt.Printf("Executing behavior %s for entity %d\n", r.AIBehaviors[i].BehaviorType, i)
		}
	}
}

func (r *Registry) UpdateCombat() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := 1; i < len(r.HasCombatState); i++ {
		if !r.HasCombatState[i] {
			continue
		}

		cs := &r.CombatStates[i]
		if cs.State == "Idle" {
			continue
		}

		cs.FramesLeft--
		if cs.FramesLeft <= 0 {
			switch cs.State {
			case "Startup":
				cs.State = "Active"
				cs.FramesLeft = cs.ActiveFrames
				fmt.Printf("Entity %d: Combat Active!\n", i)
			case "Active":
				cs.State = "Recovery"
				cs.FramesLeft = cs.RecoveryFrames
				fmt.Printf("Entity %d: Combat Recovery...\n", i)
			case "Recovery":
				cs.State = "Idle"
				fmt.Printf("Entity %d: Combat Idle.\n", i)
			}
		}
	}
}

func (r *Registry) UpdateCollision(rules []schema.EventAction) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i := 1; i < len(r.HasCollider); i++ {
		if !r.HasCollider[i] || !r.HasPosition[i] {
			continue
		}

		// Broad-phase using Spatial Grid
		nearby := r.SpatialGrid.GetNearby(r.Positions[i].X, r.Positions[i].Y)

		for _, otherID := range nearby {
			j := int(otherID)
			if i >= j {
				continue
			}
			if !r.HasCollider[j] || !r.HasPosition[j] {
				continue
			}

			p1 := r.Positions[i]
			p2 := r.Positions[j]
			c1 := r.Colliders[i]
			c2 := r.Colliders[j]

			if p1.X < p2.X+c2.Width &&
				p1.X+c1.Width > p2.X &&
				p1.Y < p2.Y+c2.Height &&
				p1.Y+c1.Height > p2.Y {
				r.handleCollision(Entity(i), Entity(j), rules)
			}
		}
	}
}

func (r *Registry) handleCollision(e1, e2 Entity, rules []schema.EventAction) {
	for _, rule := range rules {
		if strings.Contains(rule.Trigger, "COLLIDES_WITH") {
			fmt.Printf("Collision detected between %d and %d. Rule Trigger: %s, Action: %s\n", e1, e2, rule.Trigger, rule.Action)
			if rule.Action == "Damage" {
				r.ApplyDamage(e1, 10)
				r.ApplyDamage(e2, 10)
			}
		}
	}
}

func (r *Registry) ApplyDamage(e Entity, amount float64) {
	if int(e) < len(r.HasHealth) && r.HasHealth[e] {
		r.Healths[e].Current -= amount
		fmt.Printf("Entity %d damaged. Health: %.2f/%.2f\n", e, r.Healths[e].Current, r.Healths[e].Max)
	}
}

func (r *Registry) UpdateRender() {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for i := 1; i < len(r.HasSprite); i++ {
		if r.HasSprite[i] && r.HasPosition[i] {
			fmt.Printf("Rendering entity %d (sprite: %s) at (%.2f, %.2f)\n", i, r.Sprites[i].SpriteID, r.Positions[i].X, r.Positions[i].Y)
		}
	}
}
