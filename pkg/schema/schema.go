package schema

type WorldConfig struct {
	GridSpacing float64   `json:"grid_spacing"`
	Gravity     []float64 `json:"gravity"`
	Viewport    []int     `json:"viewport"`
	GlobalShader string   `json:"global_shader"`
	LevelSeed   int64     `json:"level_seed"`
	RoomCount   int       `json:"room_count"`
	Biome       string    `json:"biome"`
}

type ComponentData struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type EntitySpec struct {
	ID         uint32          `json:"id"`
	Components []ComponentData `json:"components"`
}

type EventAction struct {
	Trigger string `json:"trigger"`
	Action  string `json:"action"`
}

type GameSchema struct {
	World        WorldConfig   `json:"world"`
	Entities     []EntitySpec  `json:"entities"`
	Rules        []EventAction `json:"rules"`
	StyleKeywords []string     `json:"style_keywords"`
}
