// Package rl provides the basic data model for a roguelike game.
package rl

import "github.com/jefflund/stones/pkg/hjkl"

// Event is a message handled by Mob.
type Event any

// CollideEvent is an Event arising from a Mob coliided with a Tile obstacle.
type CollideEvent struct {
	Obstacle *Tile
}

// BumpeEvent is an Event arising from a Mob bumping another Mob.
type BumpEvent struct {
	Bumped *Mob
}

// Component processes Event for a Mob.
type Component interface {
	Process(*Mob, Event)
}

// ComponentFunc is a function which acts as Component.
type ComponentFunc func(m *Mob, v Event)

// Process implements Component by calling the function.
func (c ComponentFunc) Process(m *Mob, v Event) {
	c(m, v)
}

// EventProcessor creates a ComponentFunc which processes a single Event type.
func EventProcessor[V Event](f func(m *Mob, v *V)) ComponentFunc {
	return ComponentFunc(func(m *Mob, v Event) {
		if v, ok := v.(*V); ok {
			f(m, v)
		}
	})
}

// Mob is a game object which occupies Tile.
type Mob struct {
	Face hjkl.Glyph
	Pos  *Tile

	Components []Component
}

// NewMob constructs a Mob with the given Glyph face.
func NewMob(face hjkl.Glyph) *Mob {
	return &Mob{Face: face}
}

// Handle has each constituent Component process the Event for the Mob.
func (m *Mob) Handle(v Event) {
	for _, c := range m.Components {
		c.Process(m, v)
	}
}

// AddComponent appends a Component to the Mob.
func (m *Mob) AddComponent(c Component) {
	m.Components = append(m.Components, c)
}

// Move attempts to move the Mob to a new Tile.
func (m *Mob) Move(delta hjkl.Vector) {
	dst := m.Pos.Adjacent[delta]
	if !dst.Pass {
		m.Handle(&CollideEvent{dst})
	} else if dst.Occupant != nil {
		m.Handle(&BumpEvent{dst.Occupant})
	} else {
		m.Pos.Occupant = nil
		dst.Occupant = m
		m.Pos = dst
	}
}

// Tile is a square in the game map.
type Tile struct {
	Face     hjkl.Glyph
	Pass     bool
	Occupant *Mob

	Offset   hjkl.Vector
	Adjacent map[hjkl.Vector]*Tile
}

// NewTile creates a new Tile with the given Vector offset.
func NewTile(offset hjkl.Vector) *Tile {
	return &Tile{
		Face:     hjkl.Ch('.'),
		Pass:     true,
		Offset:   offset,
		Adjacent: make(map[hjkl.Vector]*Tile),
	}
}

// PlaceMob places a Mob on a Tile.
func PlaceMob(m *Mob, t *Tile) {
	m.Pos = t
	t.Occupant = m
}
