package assets

import (
	"fmt"
	"log"
)

type AssetType string

const (
	SpriteAsset AssetType = "Sprite"
	AudioAsset  AssetType = "Audio"
)

type AssetRequest struct {
	Type        AssetType
	Description string
}

func RequestAsset(req AssetRequest) string {
	log.Printf("Requesting dynamic asset: %s (%s)", req.Description, req.Type)
	// In a real implementation, this would call an external API (e.g., Stable Diffusion)
	// and return a URL or local path to the generated asset.
	return fmt.Sprintf("generated_%s_path", req.Description)
}
