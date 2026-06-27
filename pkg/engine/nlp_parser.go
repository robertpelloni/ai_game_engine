package engine

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/robertpelloni/ai_game_engine/pkg/schema"
	"github.com/sashabaranov/go-openai"
)

// ParseNaturalLanguageToSchema uses an LLM receiving a prompt and generating a JSON GameSchema delta.
func ParseNaturalLanguageToSchema(prompt string) (*schema.GameSchema, error) {
	log.Printf("NLP Engine: Interpreting prompt: '%s'", prompt)

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Println("NLP Engine: OPENAI_API_KEY not set, falling back to mock parser")
		return mockParseNaturalLanguageToSchema(prompt)
	}

	client := openai.NewClient(apiKey)

	systemPrompt := `You are an AI game engine level designer. You parse a human prompt into a JSON schema that dictates the level seed, biome, and room count.
The schema structure is as follows:
{
  "world": {
    "level_seed": <int>,
    "room_count": <int>,
    "biome": <string>
  }
}
Biomes can be like "Sci-Fi", "Fantasy", "Dungeon", etc. Room count is usually an integer based on words like "huge" (10), "small" (1), or a direct number.
Respond ONLY with raw, valid JSON.`

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.2,
		},
	)

	if err != nil {
		log.Printf("NLP Engine: OpenAI API Error: %v. Falling back to mock.", err)
		return mockParseNaturalLanguageToSchema(prompt)
	}

	content := resp.Choices[0].Message.Content
	log.Printf("NLP Engine: Received JSON: %s", content)

	return schema.ParseSchemaBytes([]byte(content))
}

func mockParseNaturalLanguageToSchema(prompt string) (*schema.GameSchema, error) {
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
