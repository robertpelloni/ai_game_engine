package scene

import "github.com/robertpelloni/ai_game_engine/pkg/schema"

type Scene struct {
	ID     string
	Schema *schema.GameSchema
}

func NewScene(id string, s *schema.GameSchema) *Scene {
	return &Scene{
		ID:     id,
		Schema: s,
	}
}
