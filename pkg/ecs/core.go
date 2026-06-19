package ecs

import "sync"

type Entity uint32

type ComponentType string

const (
	PositionType   ComponentType = "Position"
	VelocityType   ComponentType = "Velocity"
	SpriteType     ComponentType = "SpriteRenderer"
	ColliderType   ComponentType = "Collider"
	AIBehaviorType ComponentType = "AIBehavior"
	HealthType     ComponentType = "Health"
	CombatStateType ComponentType = "CombatState"
)

type Position struct {
	X, Y float64
}

type Velocity struct {
	VX, VY float64
}

type SpriteRenderer struct {
	SpriteID string
	Prompt string
}

type Collider struct {
	Width, Height float64
	Restitution   float64 // 0 for no bounce, 1 for perfect bounce
	Static        bool    // If true, the object doesn't move on collision
	Layer         uint32  // Collision layer bitmask
	Mask          uint32  // Collision mask bitmask
	IsTrigger     bool    // If true, collision resolution is skipped
}

type AIBehavior struct {
	BehaviorType string
}

type Health struct {
	Current, Max float64
}

type CombatState struct {
	State         string // "Startup", "Active", "Recovery", "Idle"
	FramesLeft    int
	StartupFrames  int
	ActiveFrames   int
	RecoveryFrames int
}

type Registry struct {
	Mu     sync.RWMutex
	nextID uint32

	// Global Physics Properties
	GravityX float64
	GravityY float64
	Damping  float64 // Global velocity damping (0 to 1)
	Friction float64 // Surface friction (not yet fully implemented in resolution)

	// Contiguous memory for components
	Positions    []Position
	Velocities   []Velocity
	Sprites      []SpriteRenderer
	Colliders    []Collider
	AIBehaviors  []AIBehavior
	Healths      []Health
	CombatStates []CombatState

	// Presence bitsets
	HasPosition    []bool
	HasVelocity    []bool
	HasSprite      []bool
	HasCollider    []bool
	HasAIBehavior  []bool
	HasHealth      []bool
	HasCombatState []bool

	SpatialGrid *Grid
}

func NewRegistry() *Registry {
	return &Registry{
		Positions:     make([]Position, 0),
		Velocities:    make([]Velocity, 0),
		Sprites:       make([]SpriteRenderer, 0),
		Colliders:     make([]Collider, 0),
		AIBehaviors:   make([]AIBehavior, 0),
		Healths:       make([]Health, 0),
		CombatStates:  make([]CombatState, 0),
		HasPosition:   make([]bool, 0),
		HasVelocity:   make([]bool, 0),
		HasSprite:     make([]bool, 0),
		HasCollider:   make([]bool, 0),
		HasAIBehavior: make([]bool, 0),
		HasHealth:     make([]bool, 0),
		HasCombatState: make([]bool, 0),
		SpatialGrid:   NewGrid(100.0),
		Damping:       1.0, // Default: no damping
	}
}

func (r *Registry) ensureCapacity(id uint32) {
	if int(id) >= len(r.HasPosition) {
		newSize := int(id) + 1
		r.Positions = append(r.Positions, make([]Position, newSize-len(r.Positions))...)
		r.Velocities = append(r.Velocities, make([]Velocity, newSize-len(r.Velocities))...)
		r.Sprites = append(r.Sprites, make([]SpriteRenderer, newSize-len(r.Sprites))...)
		r.Colliders = append(r.Colliders, make([]Collider, newSize-len(r.Colliders))...)
		r.AIBehaviors = append(r.AIBehaviors, make([]AIBehavior, newSize-len(r.AIBehaviors))...)
		r.Healths = append(r.Healths, make([]Health, newSize-len(r.Healths))...)
		r.CombatStates = append(r.CombatStates, make([]CombatState, newSize-len(r.CombatStates))...)

		r.HasPosition = append(r.HasPosition, make([]bool, newSize-len(r.HasPosition))...)
		r.HasVelocity = append(r.HasVelocity, make([]bool, newSize-len(r.HasVelocity))...)
		r.HasSprite = append(r.HasSprite, make([]bool, newSize-len(r.HasSprite))...)
		r.HasCollider = append(r.HasCollider, make([]bool, newSize-len(r.HasCollider))...)
		r.HasAIBehavior = append(r.HasAIBehavior, make([]bool, newSize-len(r.HasAIBehavior))...)
		r.HasHealth = append(r.HasHealth, make([]bool, newSize-len(r.HasHealth))...)
		r.HasCombatState = append(r.HasCombatState, make([]bool, newSize-len(r.HasCombatState))...)
	}
}

func (r *Registry) CreateEntity() Entity {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.nextID++
	id := r.nextID
	r.ensureCapacity(id)
	return Entity(id)
}

func (r *Registry) AddPosition(e Entity, p Position) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.ensureCapacity(uint32(e))
	r.Positions[e] = p
	r.HasPosition[e] = true
}

func (r *Registry) AddVelocity(e Entity, v Velocity) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.ensureCapacity(uint32(e))
	r.Velocities[e] = v
	r.HasVelocity[e] = true
}

func (r *Registry) AddSprite(e Entity, s SpriteRenderer) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.ensureCapacity(uint32(e))
	r.Sprites[e] = s
	r.HasSprite[e] = true
}

func (r *Registry) AddCollider(e Entity, c Collider) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.ensureCapacity(uint32(e))
	r.Colliders[e] = c
	r.HasCollider[e] = true
}

func (r *Registry) AddAIBehavior(e Entity, b AIBehavior) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.ensureCapacity(uint32(e))
	r.AIBehaviors[e] = b
	r.HasAIBehavior[e] = true
}

func (r *Registry) AddHealth(e Entity, h Health) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.ensureCapacity(uint32(e))
	r.Healths[e] = h
	r.HasHealth[e] = true
}

func (r *Registry) AddCombatState(e Entity, s CombatState) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.ensureCapacity(uint32(e))
	r.CombatStates[e] = s
	r.HasCombatState[e] = true
}
