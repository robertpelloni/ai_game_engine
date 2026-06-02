package ecs

import "sync"

type Entity uint32

type ComponentType string

const (
	PositionType ComponentType = "Position"
	VelocityType ComponentType = "Velocity"
	SpriteType   ComponentType = "SpriteRenderer"
	ColliderType ComponentType = "Collider"
	AIBehaviorType ComponentType = "AIBehavior"
)

type Position struct {
	X, Y float64
}

type Velocity struct {
	VX, VY float64
}

type SpriteRenderer struct {
	SpriteID string
}

type Collider struct {
	Width, Height float64
}

type AIBehavior struct {
	BehaviorType string
}

type Registry struct {
	mu sync.RWMutex
	nextID uint32
	entities map[Entity]bool

	Positions  map[Entity]*Position
	Velocities map[Entity]*Velocity
	Sprites    map[Entity]*SpriteRenderer
	Colliders  map[Entity]*Collider
	AIBehaviors map[Entity]*AIBehavior
}

func NewRegistry() *Registry {
	return &Registry{
		entities:    make(map[Entity]bool),
		Positions:   make(map[Entity]*Position),
		Velocities:  make(map[Entity]*Velocity),
		Sprites:     make(map[Entity]*SpriteRenderer),
		Colliders:   make(map[Entity]*Collider),
		AIBehaviors: make(map[Entity]*AIBehavior),
	}
}

func (r *Registry) CreateEntity() Entity {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nextID++
	id := Entity(r.nextID)
	r.entities[id] = true
	return id
}

func (r *Registry) AddPosition(e Entity, p *Position) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Positions[e] = p
}

func (r *Registry) AddVelocity(e Entity, v *Velocity) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Velocities[e] = v
}

func (r *Registry) AddSprite(e Entity, s *SpriteRenderer) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Sprites[e] = s
}

func (r *Registry) AddCollider(e Entity, c *Collider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Colliders[e] = c
}

func (r *Registry) AddAIBehavior(e Entity, b *AIBehavior) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.AIBehaviors[e] = b
}
