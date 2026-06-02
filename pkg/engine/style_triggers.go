package engine

import "fmt"

type StyleConfig struct {
	Resolution      [2]int
	RaycastingEnabled bool
	PaletteClamping bool
	DeferredLighting bool
	RainEffect       bool
	FrameStateMachine bool
}

func GetStyleConfig(keywords []string) StyleConfig {
	config := StyleConfig{
		Resolution: [2]int{1920, 1080},
	}

	for _, kw := range keywords {
		switch kw {
		case "Retro Raycaster":
			config.Resolution = [2]int{320, 240}
			config.RaycastingEnabled = true
			config.PaletteClamping = true
		case "Gritty Noir":
			config.DeferredLighting = true
			config.RainEffect = true
		case "Souls Combat":
			config.FrameStateMachine = true
		}
	}

	return config
}

func ApplyStyle(config StyleConfig) {
	fmt.Printf("Applying Style Config: %+v\n", config)
	// In a real engine, this would set global state, load shaders, etc.
}
