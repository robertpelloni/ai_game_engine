package engine

import (
	"testing"
	"github.com/robertpelloni/ai_game_engine/pkg/ecs"
)

func TestApplyStyle(t *testing.T) {
	reg := ecs.NewRegistry()
	reg.GravityY = 10.0
	reg.Damping = 1.0

	config := GetStyleConfig([]string{"Souls Combat"})
	ApplyStyle(reg, config)

	if reg.GravityY != 15.0 {
		t.Errorf("Expected GravityY to be modified to 15.0 by Souls Combat style, got %f", reg.GravityY)
	}

	if reg.Damping != 0.8 {
		t.Errorf("Expected Damping to be modified to 0.8 by Souls Combat style, got %f", reg.Damping)
	}
}

func TestStylePromptSuffix(t *testing.T) {
	config := GetStyleConfig([]string{"Cyberpunk"})
	if config.AssetPromptSuffix == "" {
		t.Error("Expected Cyberpunk style to provide an asset prompt suffix")
	}
}
