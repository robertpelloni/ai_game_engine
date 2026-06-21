package engine

import (
	"testing"
	"github.com/robertpelloni/ai_game_engine/pkg/ecs"
)

func TestParseRuleConditionAdvanced(t *testing.T) {
	reg := ecs.NewRegistry()
	e1 := reg.CreateEntity()

	reg.AddHealth(e1, ecs.Health{Current: 40, Max: 100})

	state := ecs.NewEntityState()
	state.Flags["IsPoisoned"] = true
	reg.AddEntityState(e1, state)

	if !ParseRuleCondition(reg, e1, 0, "Health < 50 AND Flag IsPoisoned == true") {
		t.Error("Expected condition to be true")
	}

	if ParseRuleCondition(reg, e1, 0, "Health < 30 AND Flag IsPoisoned == true") {
		t.Error("Expected condition to be false")
	}
}

func TestExecuteRuleActionAdvanced(t *testing.T) {
	reg := ecs.NewRegistry()
	e1 := reg.CreateEntity()

	reg.AddHealth(e1, ecs.Health{Current: 50, Max: 100})

	ExecuteRuleAction(reg, e1, 0, "Heal 20; SetFlag IsBuffed true")

	if reg.Healths[e1].Current != 70 {
		t.Errorf("Expected health to be 70, got %f", reg.Healths[e1].Current)
	}

	state := reg.GetEntityState(e1)
	if state == nil || !state.Flags["IsBuffed"] {
		t.Error("Expected IsBuffed flag to be true")
	}
}
