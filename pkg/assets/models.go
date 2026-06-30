package assets

import (
	"fmt"
	"sync"
)

type ModelCache struct {
	mu     sync.RWMutex
	models map[string]interface{}
}

var GlobalModelCache = &ModelCache{
	models: make(map[string]interface{}),
}

func GenerateModelAsync(modelID string, prompt string) {
	GlobalModelCache.mu.Lock()
	defer GlobalModelCache.mu.Unlock()

	// Mocking procedural model generation based on text
	fmt.Printf("Generating 3D model for %s with prompt: %s\n", modelID, prompt)
	GlobalModelCache.models[modelID] = struct{}{}
}
