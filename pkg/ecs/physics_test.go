package ecs

import (
	"math"
	"testing"
)

func TestLinearMotion(t *testing.T) {
	reg := NewRegistry()
	e := reg.CreateEntity()
	reg.AddPosition(e, Position{X: 0, Y: 0})
	reg.AddVelocity(e, Velocity{VX: 10, VY: 5})

	dt := 1.0
	reg.UpdatePhysics(dt)

	p := reg.Positions[e]
	if p.X != 10 || p.Y != 5 {
		t.Errorf("Expected (10, 5), got (%f, %f)", p.X, p.Y)
	}
}

func TestGravityApplication(t *testing.T) {
	reg := NewRegistry()
	reg.GravityY = -9.81

	e := reg.CreateEntity()
	reg.AddPosition(e, Position{X: 0, Y: 0})
	reg.AddVelocity(e, Velocity{VX: 0, VY: 0})

	dt := 1.0
	// Frame 1: Velocity updated to -9.81, Position updated to -9.81
	reg.UpdatePhysics(dt)

	v := reg.Velocities[e]
	p := reg.Positions[e]

	if v.VY != -9.81 {
		t.Errorf("Expected Velocity VY -9.81, got %f", v.VY)
	}
	if p.Y != -9.81 {
		t.Errorf("Expected Position Y -9.81, got %f", p.Y)
	}

	// Frame 2: Velocity updated to -19.62, Position updated to -9.81 + (-19.62) = -29.43
	reg.UpdatePhysics(dt)

	p = reg.Positions[e]
	expectedY := -9.81 + (-19.62)
	if math.Abs(p.Y-expectedY) > 0.0001 {
		t.Errorf("Expected Position Y %f, got %f", expectedY, p.Y)
	}
}

func TestNoVelocityNoMovement(t *testing.T) {
	reg := NewRegistry()
	reg.GravityY = -9.81

	e := reg.CreateEntity()
	reg.AddPosition(e, Position{X: 100, Y: 100})
	// Entity has Position but NO Velocity component

	reg.UpdatePhysics(1.0)

	p := reg.Positions[e]
	if p.X != 100 || p.Y != 100 {
		t.Errorf("Entity moved without velocity component: got (%f, %f)", p.X, p.Y)
	}
}
