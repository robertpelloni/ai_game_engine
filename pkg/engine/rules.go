package engine

import (
	"strings"
	"strconv"
	"github.com/robertpelloni/ai_game_engine/pkg/ecs"
)

// ParseRuleCondition evaluates a pseudo-AST string against two colliding entities
func ParseRuleCondition(reg *ecs.Registry, e1, e2 ecs.Entity, condition string) bool {
	// Nested conditions evaluation (very basic AND support)
	if strings.Contains(condition, " AND ") {
		parts := strings.Split(condition, " AND ")
		for _, part := range parts {
			if !evaluateSingleCondition(reg, e1, e2, strings.TrimSpace(part)) {
				return false
			}
		}
		return true
	}

	return evaluateSingleCondition(reg, e1, e2, condition)
}

func evaluateSingleCondition(reg *ecs.Registry, e1, e2 ecs.Entity, condition string) bool {
	if strings.HasPrefix(condition, "Health <") {
		parts := strings.Split(condition, " ")
		if len(parts) >= 3 && parts[1] == "<" {
			val, err := strconv.ParseFloat(parts[2], 64)
			if err == nil {
				if int(e1) < len(reg.HasHealth) && reg.HasHealth[e1] {
					if reg.Healths[e1].Current >= val {
						return false
					}
				}
			}
		}
		return true
	}

	if strings.HasPrefix(condition, "Flag ") {
		parts := strings.Split(condition, " ")
		if len(parts) >= 4 && parts[2] == "==" {
			flagName := parts[1] // e.g. Flag IsPoisoned == true
			expected, err := strconv.ParseBool(parts[3])
			if err == nil {
				state := reg.GetEntityState(e1)
				if state != nil {
					val, exists := state.Flags[flagName]
					if exists && val == expected {
						return true
					}
				}
			}
		}
		return false
	}

	if condition == "" || condition == "COLLIDES_WITH" {
		return true
	}

	return false
}

// ExecuteRuleAction executes an arbitrary action string on entities
func ExecuteRuleAction(reg *ecs.Registry, e1, e2 ecs.Entity, action string) {
	// Support multiple actions delimited by ';'
	actions := strings.Split(action, ";")

	for _, act := range actions {
		act = strings.TrimSpace(act)
		if strings.HasPrefix(act, "SetFlag ") {
			parts := strings.Split(act, " ")
			if len(parts) == 3 {
				flagName := parts[1]
				val, err := strconv.ParseBool(parts[2])
				if err == nil {
					state := reg.GetEntityState(e1)
					if state == nil {
						newState := ecs.NewEntityState()
						newState.Flags[flagName] = val
						reg.AddEntityState(e1, newState)
					} else {
						reg.Mu.Lock()
						state.Flags[flagName] = val
						reg.Mu.Unlock()
					}
				}
			}
			continue
		}

		if strings.HasPrefix(act, "Heal ") {
			parts := strings.Split(act, " ")
			if len(parts) == 2 {
				val, err := strconv.ParseFloat(parts[1], 64)
				if err == nil {
					reg.ApplyDamage(e1, -val) // Negative damage heals
				}
			}
			continue
		}

		switch act {
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
				reg.Velocities[e1].VX = -reg.Velocities[e1].VX * 2
				reg.Velocities[e1].VY = -reg.Velocities[e1].VY * 2
			}
		}
	}
}
