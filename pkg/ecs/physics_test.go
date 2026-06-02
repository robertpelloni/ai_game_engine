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

func TestCollisionResolution(t *testing.T) {
	reg := NewRegistry()
	e1 := reg.CreateEntity()
	e2 := reg.CreateEntity()

	reg.AddPosition(e1, Position{X: 0, Y: 0})
	reg.AddCollider(e1, Collider{Width: 10, Height: 10, Restitution: 0.5})
	reg.AddVelocity(e1, Velocity{VX: 10, VY: 0})

	reg.AddPosition(e2, Position{X: 8, Y: 0})
	reg.AddCollider(e2, Collider{Width: 10, Height: 10, Static: true})

	reg.UpdatePhysics(0)
	reg.UpdateCollision(nil)

	if reg.Positions[e1].X >= 0 {
		t.Errorf("Expected e1 to be pushed back, got X=%f", reg.Positions[e1].X)
	}
	if reg.Velocities[e1].VX >= 0 {
		t.Errorf("Expected e1 velocity to be reflected, got VX=%f", reg.Velocities[e1].VX)
	}
}

func TestRaycasting(t *testing.T) {
	reg := NewRegistry()
	e := reg.CreateEntity()
	reg.AddPosition(e, Position{X: 100, Y: 100})
	reg.AddCollider(e, Collider{Width: 10, Height: 10})

	reg.UpdatePhysics(0)

	hit, dist := reg.Raycast(0, 105, 1, 0, 200)
	if hit != e {
		t.Errorf("Expected to hit entity %d, got %d", e, hit)
	}
	if dist < 90 || dist > 110 {
		t.Errorf("Expected distance ~100, got %f", dist)
	}
}

func TestDiagonalCollisionResolution(t *testing.T) {
	reg := NewRegistry()
	e1 := reg.CreateEntity()
	e2 := reg.CreateEntity()

	// Entities approaching diagonally
	reg.AddPosition(e1, Position{X: 0, Y: 0})
	reg.AddCollider(e1, Collider{Width: 10, Height: 10, Restitution: 0.0})
	reg.AddVelocity(e1, Velocity{VX: 10, VY: 10})

	reg.AddPosition(e2, Position{X: 8, Y: 8})
	reg.AddCollider(e2, Collider{Width: 10, Height: 10, Static: true})

	reg.UpdatePhysics(0)
	reg.UpdateCollision(nil)

	// One axis should be resolved
	p1 := reg.Positions[e1]
	if p1.X >= 0 && p1.Y >= 0 {
		t.Errorf("Expected diagonal resolution to push e1 back, got (%f, %f)", p1.X, p1.Y)
	}
}
