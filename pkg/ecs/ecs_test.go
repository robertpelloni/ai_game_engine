package ecs

import (
	"github.com/robertpelloni/ai_game_engine/pkg/schema"
	"testing"
)

func TestRegistry(t *testing.T) {
	reg := NewRegistry()
	e := reg.CreateEntity()

	if e != 1 {
		t.Errorf("Expected entity ID 1, got %d", e)
	}

	reg.AddPosition(e, Position{X: 10, Y: 20})
	reg.AddVelocity(e, Velocity{VX: 1, VY: 1})

	reg.UpdatePhysics(1.0)

	p := reg.Positions[e]
	if p.X != 11 || p.Y != 21 {
		t.Errorf("Expected position (11, 21), got (%f, %f)", p.X, p.Y)
	}
}

func TestCollisionAndDamage(t *testing.T) {
	reg := NewRegistry()
	e1 := reg.CreateEntity()
	e2 := reg.CreateEntity()

	reg.AddPosition(e1, Position{X: 0, Y: 0})
	reg.AddCollider(e1, Collider{Width: 10, Height: 10, Layer: 1, Mask: 1})
	reg.AddHealth(e1, Health{Current: 100, Max: 100})

	reg.AddPosition(e2, Position{X: 5, Y: 5})
	reg.AddCollider(e2, Collider{Width: 10, Height: 10, Layer: 1, Mask: 1})
	reg.AddHealth(e2, Health{Current: 100, Max: 100})

	// Spatial grid needs to be populated
	reg.UpdatePhysics(0)

	rules := []schema.EventAction{
		{Trigger: "COLLIDES_WITH", Action: "Damage"},
	}

	reg.UpdateCollision(rules)

	if reg.Healths[e1].Current != 90 {
		t.Errorf("Expected health 90, got %f", reg.Healths[e1].Current)
	}
}

func TestCombatStateMachine(t *testing.T) {
	reg := NewRegistry()
	e := reg.CreateEntity()
	reg.AddCombatState(e, CombatState{
		State:          "Startup",
		FramesLeft:     2,
		StartupFrames:  2,
		ActiveFrames:   3,
		RecoveryFrames: 2,
	})

	reg.UpdateCombat() // FramesLeft 1
	if reg.CombatStates[e].State != "Startup" {
		t.Errorf("Expected Startup state, got %s", reg.CombatStates[e].State)
	}

	reg.UpdateCombat() // FramesLeft 0 -> Active
	if reg.CombatStates[e].State != "Active" {
		t.Errorf("Expected Active state, got %s", reg.CombatStates[e].State)
	}
}
