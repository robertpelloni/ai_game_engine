package engine

import (
	"strings"
	"strconv"
	"github.com/robertpelloni/ai_game_engine/pkg/ecs"
)

// ParseRuleCondition evaluates a pseudo-AST string against two colliding entities
func ParseRuleCondition(reg *ecs.Registry, e1, e2 ecs.Entity, condition string) bool {
	if strings.Contains(condition, "Health <") {
		parts := strings.Split(condition, " ")
		for i, part := range parts {
			if part == "<" && i > 0 && i+1 < len(parts) {
				val, err := strconv.ParseFloat(parts[i+1], 64)
				if err == nil {
					if int(e1) < len(reg.HasHealth) && reg.HasHealth[e1] {
						if reg.Healths[e1].Current >= val {
							return false
						}
					}
				}
			}
		}
	}

	if condition == "" || strings.Contains(condition, "COLLIDES_WITH") {
		return true
	}

	return false
}

// ExecuteRuleAction executes an arbitrary action string on entities
func ExecuteRuleAction(reg *ecs.Registry, e1, e2 ecs.Entity, action string) {
	switch action {
	case "Damage":
		reg.ApplyDamage(e1, 10)
		reg.ApplyDamage(e2, 10)
	case "Stop":
		if int(e1) < len(reg.HasVelocity) && reg.HasVelocity[e1] {
			reg.Velocities[e1] = ecs.Velocity{0, 0}
		}
		if int(e2) < len(reg.HasVelocity) && reg.HasVelocity[e2] {
			reg.Velocities[e2] = ecs.Velocity{0, 0}
		}
	case "RunAway":
		if int(e1) < len(reg.HasVelocity) && reg.HasVelocity[e1] {
			// Reverse velocity and boost speed to "run away"
			reg.Velocities[e1].VX = -reg.Velocities[e1].VX * 2
			reg.Velocities[e1].VY = -reg.Velocities[e1].VY * 2
		}
	}
}
