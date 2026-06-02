package engine

import (
	"github.com/robertpelloni/ai_game_engine/pkg/ecs"
	"testing"
)

func Benchmark1000Entities(b *testing.B) {
	reg := ecs.NewRegistry()
	for i := 1; i <= 1000; i++ {
		e := reg.CreateEntity()
		reg.AddPosition(e, ecs.Position{X: float64(i), Y: float64(i)})
		reg.AddVelocity(e, ecs.Velocity{VX: 1.0, VY: 1.0})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reg.UpdatePhysics(0.016)
	}
}
