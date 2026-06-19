package assets

import (
	"image/color"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type TextureCache struct {
	mu         sync.RWMutex
	textures   map[string]*ebiten.Image
	generating map[string]bool
}

var Cache = &TextureCache{
	textures:   make(map[string]*ebiten.Image),
	generating: make(map[string]bool),
}

func GetTexture(spriteID string) *ebiten.Image {
	Cache.mu.RLock()
	defer Cache.mu.RUnlock()
	return Cache.textures[spriteID]
}

func GenerateAssetAsync(spriteID string, prompt string, width, height int) {
	// If it already exists or is empty, skip
	if prompt == "" {
		return
	}

	Cache.mu.Lock()
	if _, exists := Cache.textures[spriteID]; exists {
		Cache.mu.Unlock()
		return
	}
	if Cache.generating[spriteID] {
		Cache.mu.Unlock()
		return
	}
	Cache.generating[spriteID] = true
	Cache.mu.Unlock()

	go func() {
		log.Printf("AssetManager: Starting generation for [%s] based on prompt: '%s'", spriteID, prompt)

		// Simulate API Latency (e.g., Stable Diffusion/DALL-E)
		time.Sleep(time.Duration(1000 + rand.Intn(2000)) * time.Millisecond)

		// Create a mock image based on the prompt hash/length to give it a unique pseudo-random color
		img := ebiten.NewImage(width, height)

		// Basic hash of string for color
		hash := 0
		for _, c := range prompt {
			hash += int(c)
		}

		r := uint8((hash * 13) % 255)
		g := uint8((hash * 27) % 255)
		b := uint8((hash * 7) % 255)

		img.Fill(color.RGBA{r, g, b, 255})

		Cache.mu.Lock()
		Cache.textures[spriteID] = img
		delete(Cache.generating, spriteID)
		Cache.mu.Unlock()

		log.Printf("AssetManager: Completed generation for [%s]", spriteID)
	}()
}
