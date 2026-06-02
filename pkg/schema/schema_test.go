package schema

import (
	"testing"
)

func TestParseSchemaBytes(t *testing.T) {
	jsonData := `{
		"world": {
			"grid_spacing": 1.0,
			"gravity": [0.0, -9.81],
			"viewport": [800, 600],
			"global_shader": "basic"
		},
		"entities": [
			{
				"id": 1,
				"components": [
					{"type": "Position", "data": {"x": 10.0, "y": 20.0}}
				]
			}
		],
		"rules": [],
		"style_keywords": ["Retro Raycaster"]
	}`

	schema, err := ParseSchemaBytes([]byte(jsonData))
	if err != nil {
		t.Fatalf("Failed to parse schema: %v", err)
	}

	if schema.World.GridSpacing != 1.0 {
		t.Errorf("Expected GridSpacing 1.0, got %f", schema.World.GridSpacing)
	}

	if len(schema.Entities) != 1 {
		t.Errorf("Expected 1 entity, got %d", len(schema.Entities))
	}

	if schema.Entities[0].ID != 1 {
		t.Errorf("Expected entity ID 1, got %d", schema.Entities[0].ID)
	}
}
