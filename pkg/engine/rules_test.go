package engine

import (
	"testing"
	"github.com/robertpelloni/ai_game_engine/pkg/ecs"
)

func TestParseRuleCondition(t *testing.T) {
	reg := ecs.NewRegistry()
	e1 := reg.CreateEntity()

	// Test basic
	if !ParseRuleCondition(reg, e1, 0, "COLLIDES_WITH") {
		t.Errorf("Expected basic collision to be true")
	}

	// Test logic gate
	reg.AddHealth(e1, ecs.Health{Current: 40.0, Max: 100.0})

	if !ParseRuleCondition(reg, e1, 0, "Health < 50 AND COLLIDES_WITH") {
		t.Errorf("Expected logic parse to be true")
	}

	reg.Healths[e1].Current = 60.0
	if ParseRuleCondition(reg, e1, 0, "Health < 50 AND COLLIDES_WITH") {
		t.Errorf("Expected logic parse to be false")
	}
}

func TestExecuteRuleAction(t *testing.T) {
	reg := ecs.NewRegistry()
	e1 := reg.CreateEntity()
	reg.AddVelocity(e1, ecs.Velocity{VX: 10, VY: -10})

	ExecuteRuleAction(reg, e1, 0, "RunAway")

	if reg.Velocities[e1].VX != -20 || reg.Velocities[e1].VY != 20 {
		t.Errorf("RunAway logic did not reverse and multiply velocities")
	}
}
