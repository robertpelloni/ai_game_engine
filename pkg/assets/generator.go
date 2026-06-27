package assets

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sashabaranov/go-openai"
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

		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey != "" {
			client := openai.NewClient(apiKey)
			req := openai.ImageRequest{
				Prompt:         prompt + ", pixel art style, top down 2d game asset, simple background",
				Size:           openai.CreateImageSize256x256,
				ResponseFormat: openai.CreateImageResponseFormatB64JSON,
				N:              1,
			}

			resp, err := client.CreateImage(context.Background(), req)
			if err == nil && len(resp.Data) > 0 {
				imgBytes, err := base64.StdEncoding.DecodeString(resp.Data[0].B64JSON)
				if err == nil {
					img, _, err := image.Decode(bytes.NewReader(imgBytes))
					if err == nil {
						ebitenImg := ebiten.NewImageFromImage(img)

						Cache.mu.Lock()
						Cache.textures[spriteID] = ebitenImg
						delete(Cache.generating, spriteID)
						Cache.mu.Unlock()
						log.Printf("AssetManager: Completed API generation for [%s]", spriteID)
						return
					}
				}
			} else {
				log.Printf("AssetManager: API Image Gen failed: %v, falling back to mock", err)
			}
		}

		// Fallback to mock behavior
		time.Sleep(time.Duration(1000 + rand.Intn(2000)) * time.Millisecond)
		img := ebiten.NewImage(width, height)

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

		log.Printf("AssetManager: Completed mock generation for [%s]", spriteID)
	}()
}
