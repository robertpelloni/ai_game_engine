package ecs

// EntityState holds custom key-value state for advanced rule scripting.
type EntityState struct {
	Flags   map[string]bool
	Numbers map[string]float64
	Strings map[string]string
}

func NewEntityState() EntityState {
	return EntityState{
		Flags:   make(map[string]bool),
		Numbers: make(map[string]float64),
		Strings: make(map[string]string),
	}
}

func (r *Registry) ensureCapacityState(id uint32) {
	if int(id) >= len(r.HasState) {
		newSize := int(id) + 1
		r.EntityStates = append(r.EntityStates, make([]EntityState, newSize-len(r.EntityStates))...)
		r.HasState = append(r.HasState, make([]bool, newSize-len(r.HasState))...)
	}
}

func (r *Registry) AddEntityState(e Entity, s EntityState) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.ensureCapacityState(uint32(e))
	r.EntityStates[e] = s
	r.HasState[e] = true
}

// GetEntityFlag safely reads a flag by value.
func (r *Registry) GetEntityFlag(e Entity, flag string) (bool, bool) {
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	if int(e) < len(r.HasState) && r.HasState[e] {
		val, exists := r.EntityStates[e].Flags[flag]
		return val, exists
	}
	return false, false
}

// SetEntityFlag safely writes a flag.
func (r *Registry) SetEntityFlag(e Entity, flag string, val bool) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	// Create state if missing
	r.ensureCapacityState(uint32(e))
	if !r.HasState[e] {
		r.EntityStates[e] = NewEntityState()
		r.HasState[e] = true
	}

	// Initialize map if nil (defensive)
	if r.EntityStates[e].Flags == nil {
		r.EntityStates[e].Flags = make(map[string]bool)
	}

	r.EntityStates[e].Flags[flag] = val
}

// GetHealth safely retrieves health values.
func (r *Registry) GetHealth(e Entity) (Health, bool) {
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	if int(e) < len(r.HasHealth) && r.HasHealth[e] {
		return r.Healths[e], true
	}
	return Health{}, false
}

// StopVelocity safely sets velocity to zero.
func (r *Registry) StopVelocity(e Entity) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	if int(e) < len(r.HasVelocity) && r.HasVelocity[e] {
		r.Velocities[e] = Velocity{0, 0}
	}
}

// ReverseVelocity safely reverses velocity vectors.
func (r *Registry) ReverseVelocity(e Entity) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	if int(e) < len(r.HasVelocity) && r.HasVelocity[e] {
		r.Velocities[e].VX = -r.Velocities[e].VX * 2
		r.Velocities[e].VY = -r.Velocities[e].VY * 2
	}
}
