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

// Add methods to the Registry for EntityState
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

func (r *Registry) GetEntityState(e Entity) *EntityState {
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	if int(e) < len(r.HasState) && r.HasState[e] {
		return &r.EntityStates[e]
	}
	return nil
}
