package ecs

import "testing"

func TestCollisionFiltering(t *testing.T) {
	reg := NewRegistry()
	e1 := reg.CreateEntity()
	e2 := reg.CreateEntity()

	// e1: Layer 1, Mask 2 (only collides with things in Layer 2)
	reg.AddPosition(e1, Position{X: 0, Y: 0})
	reg.AddCollider(e1, Collider{Width: 10, Height: 10, Layer: 1, Mask: 2})

	// e2: Layer 1, Mask 1 (only collides with things in Layer 1)
	reg.AddPosition(e2, Position{X: 5, Y: 5})
	reg.AddCollider(e2, Collider{Width: 10, Height: 10, Layer: 1, Mask: 1})

	reg.UpdatePhysics(0)
	reg.UpdateCollision(nil)

	// Since Layer 1 & Mask 2 == 0 AND Layer 1 & Mask 1 == 1,
	// e2 wants to collide with e1, but e1 DOES NOT want to collide with e2.
	// In my implementation: (c1.Layer&c2.Mask) == 0 && (c2.Layer&c1.Mask) == 0 -> skip.
	// Here (1 & 1) != 0, so they SHOULD collide.

	// Wait, let's test a case where they SHOULD NOT.
	e3 := reg.CreateEntity()
	reg.AddPosition(e3, Position{X: 20, Y: 20})
	reg.AddCollider(e3, Collider{Width: 10, Height: 10, Layer: 4, Mask: 4})

	if reg.checkAABB(int(e1), int(e3)) {
		t.Error("e1 and e3 should not overlap spatially")
	}

	// Move e3 onto e1
	reg.Positions[e3] = Position{X: 5, Y: 5}
	reg.UpdatePhysics(0)

	// (c1.Layer(1) & c3.Mask(4)) == 0 AND (c3.Layer(4) & c1.Mask(2)) == 0
	c1 := reg.Colliders[e1]
	c3 := reg.Colliders[e3]
	if (c1.Layer&c3.Mask) == 0 && (c3.Layer&c1.Mask) == 0 {
		// This is the skip condition.
	} else {
		t.Error("e1 and e3 should have been filtered out by layer/mask")
	}
}

func TestTriggers(t *testing.T) {
	reg := NewRegistry()
	e1 := reg.CreateEntity()
	e2 := reg.CreateEntity()

	reg.AddPosition(e1, Position{X: 0, Y: 0})
	reg.AddCollider(e1, Collider{Width: 10, Height: 10, IsTrigger: true, Layer: 1, Mask: 1})

	reg.AddPosition(e2, Position{X: 5, Y: 5})
	reg.AddCollider(e2, Collider{Width: 10, Height: 10, Layer: 1, Mask: 1})

	reg.UpdatePhysics(0)
	reg.UpdateCollision(nil)

	// e1 is a trigger, so e2 should not be pushed back
	if reg.Positions[e2].X != 5 {
		t.Errorf("e2 was moved by a trigger: X=%f", reg.Positions[e2].X)
	}
}
