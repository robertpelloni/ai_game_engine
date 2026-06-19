package engine

import (
	"encoding/json"
	"log"
	"math/rand"
	"strings"

	"github.com/robertpelloni/ai_game_engine/pkg/schema"
)

// ParseNaturalLanguageToSchema simulates an LLM receiving a prompt and generating a JSON GameSchema delta.
func ParseNaturalLanguageToSchema(prompt string) (*schema.GameSchema, error) {
	log.Printf("NLP Engine: Interpreting prompt: '%s'", prompt)

	// Mock LLM Heuristics
	lowerPrompt := strings.ToLower(prompt)

	biome := "Sci-Fi"
	if strings.Contains(lowerPrompt, "castle") || strings.Contains(lowerPrompt, "fantasy") || strings.Contains(lowerPrompt, "dungeon") {
		biome = "Fantasy"
	}

	roomCount := 3
	if strings.Contains(lowerPrompt, "huge") || strings.Contains(lowerPrompt, "massive") || strings.Contains(lowerPrompt, "large") {
		roomCount = 10
	} else if strings.Contains(lowerPrompt, "small") || strings.Contains(lowerPrompt, "tiny") {
		roomCount = 1
	}

	seed := int64(rand.Intn(99999))

	// Create the "AI Generated" JSON payload
	mockJSON := map[string]interface{}{
		"world": map[string]interface{}{
			"level_seed": seed,
			"room_count": roomCount,
			"biome":      biome,
		},
	}

	bytes, err := json.Marshal(mockJSON)
	if err != nil {
		return nil, err
	}

	return schema.ParseSchemaBytes(bytes)
}
