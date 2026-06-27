package engine

import (
	"strings"
	"strconv"
	"github.com/robertpelloni/ai_game_engine/pkg/ecs"
)

// ParseRuleCondition evaluates a pseudo-AST string against two colliding entities
func ParseRuleCondition(reg *ecs.Registry, e1, e2 ecs.Entity, condition string) bool {
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
				h, exists := reg.GetHealth(e1)
				if exists && h.Current >= val {
					return false
				}
			}
		}
		return true
	}

	if strings.HasPrefix(condition, "Flag ") {
		parts := strings.Split(condition, " ")
		if len(parts) >= 4 && parts[2] == "==" {
			flagName := parts[1]
			expected, err := strconv.ParseBool(parts[3])
			if err == nil {
				val, exists := reg.GetEntityFlag(e1, flagName)
				if exists && val == expected {
					return true
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
	actions := strings.Split(action, ";")

	for _, act := range actions {
		act = strings.TrimSpace(act)
		if strings.HasPrefix(act, "SetFlag ") {
			parts := strings.Split(act, " ")
			if len(parts) == 3 {
				flagName := parts[1]
				val, err := strconv.ParseBool(parts[2])
				if err == nil {
					reg.SetEntityFlag(e1, flagName, val)
				}
			}
			continue
		}

		if strings.HasPrefix(act, "Heal ") {
			parts := strings.Split(act, " ")
			if len(parts) == 2 {
				val, err := strconv.ParseFloat(parts[1], 64)
				if err == nil {
					reg.ApplyDamage(e1, -val)
				}
			}
			continue
		}

		switch act {
		case "Damage":
			reg.ApplyDamage(e1, 10)
			reg.ApplyDamage(e2, 10)
		case "Stop":
			reg.StopVelocity(e1)
			reg.StopVelocity(e2)
		case "RunAway":
			reg.ReverseVelocity(e1)
		}
	}
}
