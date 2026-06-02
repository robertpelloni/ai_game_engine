package ecs

import (
	"fmt"
	"github.com/robertpelloni/ai_game_engine/pkg/schema"
	"strings"
)

func (r *Registry) UpdatePhysics(dt float64) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for e, v := range r.Velocities {
		if p, ok := r.Positions[e]; ok {
			p.X += v.VX * dt
			p.Y += v.VY * dt
		}
	}
}

func (r *Registry) UpdateBehavior() {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for e, b := range r.AIBehaviors {
		fmt.Printf("Executing behavior %s for entity %d\n", b.BehaviorType, e)
	}
}

func (r *Registry) UpdateCollision(rules []schema.EventAction) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entities := make([]Entity, 0, len(r.Colliders))
	for e := range r.Colliders {
		entities = append(entities, e)
	}

	for i := 0; i < len(entities); i++ {
		for j := i + 1; j < len(entities); j++ {
			e1 := entities[i]
			e2 := entities[j]

			p1, ok1 := r.Positions[e1]
			p2, ok2 := r.Positions[e2]
			c1 := r.Colliders[e1]
			c2 := r.Colliders[e2]

			if ok1 && ok2 {
				if p1.X < p2.X+c2.Width &&
					p1.X+c1.Width > p2.X &&
					p1.Y < p2.Y+c2.Height &&
					p1.Y+c1.Height > p2.Y {
					r.handleCollision(e1, e2, rules)
				}
			}
		}
	}
}

func (r *Registry) handleCollision(e1, e2 Entity, rules []schema.EventAction) {
	for _, rule := range rules {
		if strings.Contains(rule.Trigger, "COLLIDES_WITH") {
			// Basic rule parsing logic
			fmt.Printf("Collision detected between %d and %d. Rule Trigger: %s, Action: %s\n", e1, e2, rule.Trigger, rule.Action)
			// Implementation of CALL SystemAction(Damage, 10) etc would go here
		}
	}
}

func (r *Registry) UpdateRender() {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for e, s := range r.Sprites {
		if p, ok := r.Positions[e]; ok {
			fmt.Printf("Rendering entity %d (sprite: %s) at (%.2f, %.2f)\n", e, s.SpriteID, p.X, p.Y)
		}
	}
}
