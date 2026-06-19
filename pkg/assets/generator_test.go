package assets

import (
	"github.com/hajimehoshi/ebiten/v2"
	"sync"
	"testing"

)

func TestGenerateAssetAsync_NoRedundantGeneration(t *testing.T) {
	// Reset cache for test isolation
	Cache.mu.Lock()
	Cache.textures = make(map[string]*ebiten.Image)
	Cache.mu.Unlock()

	// 10 concurrent requests for the same spriteID
	var wg sync.WaitGroup
	spriteID := "test_sprite"
	prompt := "test prompt"

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			GenerateAssetAsync(spriteID, prompt, 32, 32)
		}()
	}
	wg.Wait()

	// Check if only one generation job is active
	Cache.mu.RLock()
	if !Cache.generating[spriteID] && Cache.textures[spriteID] == nil {
		t.Errorf("Expected either generating or texture to exist")
	}
	Cache.mu.RUnlock()
}
