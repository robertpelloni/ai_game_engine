package ecs

import (
	"fmt"
	"github.com/robertpelloni/ai_game_engine/pkg/schema"
	"math"
	"strings"
)

func (r *Registry) UpdatePhysics(dt float64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.SpatialGrid.Clear()

	for i := 1; i < len(r.HasVelocity); i++ {
		if r.HasVelocity[i] {
			// Apply Gravity
			r.Velocities[i].VX += r.GravityX * dt
			r.Velocities[i].VY += r.GravityY * dt

			if r.HasPosition[i] {
				r.Positions[i].X += r.Velocities[i].VX * dt
				r.Positions[i].Y += r.Velocities[i].VY * dt
			}
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
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := 1; i < len(r.HasCollider); i++ {
		if !r.HasCollider[i] || !r.HasPosition[i] {
			continue
		}

		nearby := r.SpatialGrid.GetNearby(r.Positions[i].X, r.Positions[i].Y)

		for _, otherID := range nearby {
			j := int(otherID)
			if i >= j {
				continue
			}
			if !r.HasCollider[j] || !r.HasPosition[j] {
				continue
			}

			if r.checkAABB(i, j) {
				r.resolveCollision(i, j)
				r.handleCollision(Entity(i), Entity(j), rules)
			}
		}
	}
}

func (r *Registry) checkAABB(i, j int) bool {
	p1 := r.Positions[i]
	p2 := r.Positions[j]
	c1 := r.Colliders[i]
	c2 := r.Colliders[j]

	return p1.X < p2.X+c2.Width &&
		p1.X+c1.Width > p2.X &&
		p1.Y < p2.Y+c2.Height &&
		p1.Y+c1.Height > p2.Y
}

func (r *Registry) resolveCollision(i, j int) {
	c1 := r.Colliders[i]
	c2 := r.Colliders[j]

	// If both are static, do nothing
	if c1.Static && c2.Static {
		return
	}

	p1 := &r.Positions[i]
	p2 := &r.Positions[j]

	// Calculate overlap on both axes
	overlapX := math.Min(p1.X+c1.Width, p2.X+c2.Width) - math.Max(p1.X, p2.X)
	overlapY := math.Min(p1.Y+c1.Height, p2.Y+c2.Height) - math.Max(p1.Y, p2.Y)

	if overlapX < overlapY {
		// Resolve on X axis
		separation := overlapX
		if p1.X < p2.X {
			separation = -overlapX
		}
		if !c1.Static && !c2.Static {
			p1.X += separation / 2
			p2.X -= separation / 2
		} else if !c1.Static {
			p1.X += separation
		} else {
			p2.X -= separation
		}
		r.reflectVelocity(i, j, true)
	} else {
		// Resolve on Y axis
		separation := overlapY
		if p1.Y < p2.Y {
			separation = -overlapY
		}
		if !c1.Static && !c2.Static {
			p1.Y += separation / 2
			p2.Y -= separation / 2
		} else if !c1.Static {
			p1.Y += separation
		} else {
			p2.Y -= separation
		}
		r.reflectVelocity(i, j, false)
	}
}

func (r *Registry) reflectVelocity(i, j int, axisX bool) {
	c1 := r.Colliders[i]
	c2 := r.Colliders[j]

	if axisX {
		if r.HasVelocity[i] && !c1.Static {
			r.Velocities[i].VX *= -c1.Restitution
		}
		if r.HasVelocity[j] && !c2.Static {
			r.Velocities[j].VX *= -c2.Restitution
		}
	} else {
		if r.HasVelocity[i] && !c1.Static {
			r.Velocities[i].VY *= -c1.Restitution
		}
		if r.HasVelocity[j] && !c2.Static {
			r.Velocities[j].VY *= -c2.Restitution
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
			} else if rule.Action == "Stop" {
				if r.HasVelocity[e1] { r.Velocities[e1] = Velocity{0, 0} }
				if r.HasVelocity[e2] { r.Velocities[e2] = Velocity{0, 0} }
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

func (r *Registry) Raycast(x, y, dx, dy float64, maxDist float64) (Entity, float64) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Normalize direction
	mag := math.Sqrt(dx*dx + dy*dy)
	if mag == 0 { return 0, 0 }
	dx /= mag
	dy /= mag

	step := r.SpatialGrid.CellSize / 2
	for dist := 0.0; dist < maxDist; dist += step {
		currX := x + dx*dist
		currY := y + dy*dist
		nearby := r.SpatialGrid.GetNearby(currX, currY)
		for _, e := range nearby {
			if !r.HasCollider[e] || !r.HasPosition[e] { continue }
			p := r.Positions[e]
			c := r.Colliders[e]
			if currX >= p.X && currX <= p.X+c.Width && currY >= p.Y && currY <= p.Y+c.Height {
				return e, dist
			}
		}
	}
	return 0, 0
}
