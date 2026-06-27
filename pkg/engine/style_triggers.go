package engine

import (
	"fmt"
	"github.com/robertpelloni/ai_game_engine/pkg/ecs"
)

type StyleConfig struct {
	Resolution        [2]int
	RaycastingEnabled bool
	PaletteClamping   bool
	DeferredLighting  bool
	RainEffect        bool
	FrameStateMachine bool
	GravityModifier   float64
	DampingModifier   float64
	AssetPromptSuffix string
}

func GetStyleConfig(keywords []string) StyleConfig {
	config := StyleConfig{
		Resolution:        [2]int{1920, 1080},
		GravityModifier:   1.0,
		DampingModifier:   1.0,
		AssetPromptSuffix: "",
	}

	for _, kw := range keywords {
		switch kw {
		case "Retro Raycaster":
			config.Resolution = [2]int{320, 240}
			config.RaycastingEnabled = true
			config.PaletteClamping = true
			config.AssetPromptSuffix = ", retro doom style pixel art, 32-bit"
		case "Gritty Noir":
			config.DeferredLighting = true
			config.RainEffect = true
			config.AssetPromptSuffix = ", dark gritty neo-noir detective, black and white style"
		case "Souls Combat":
			config.FrameStateMachine = true
			config.GravityModifier = 1.5 // Heavier gravity for souls-like physics
			config.DampingModifier = 0.8
		case "Cyberpunk":
			config.AssetPromptSuffix = ", neon cyberpunk high tech low life, synthwave lighting"
		}
	}

	return config
}

// ApplyStyle applies global configuration from the StyleConfig directly into the ECS Registry.
func ApplyStyle(reg *ecs.Registry, config StyleConfig) {
	fmt.Printf("Applying Style Config: %+v\n", config)

	reg.Mu.Lock()
	defer reg.Mu.Unlock()

	// Apply modifiers to ECS global state
	reg.Damping = 1.0 // Reset to default before multiplier
	reg.GravityY *= config.GravityModifier
	reg.Damping *= config.DampingModifier
}
