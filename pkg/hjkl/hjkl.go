// Package hjkl is a terminal game engine.
package hjkl

import (
	"errors"
	"time"
)

// Game contains the functionality needed to run a game.
type Game interface {
	Update([]Key) error
	Draw(Canvas)
}

// Termination is a special error indicating normal game termination.
var Termination = errors.New("Termination")

// RunConfig stores options for Run.
type RunConfig struct {
	Terminal Terminal
	TPS      int
}

// DefaultRunConfig creates a RunConfig with default settings.
func DefaultRunConfig() *RunConfig {
	return &RunConfig{
		Terminal: TermboxTerminal{},
		TPS:      20,
	}
}

// RunOption is a function which mutates a RunConfig for Run.
type RunOption func(*RunConfig)

// WithTerm gets a RunOption which sets the Terminal of a RunConfig.
func WithTerm(t Terminal) RunOption {
	return func(r *RunConfig) {
		r.Terminal = t
	}
}

// WithTPS gets a RunOption which sets the TPS of a RunConfig.
func WithTPS(tps int) RunOption {
	return func(r *RunConfig) {
		r.TPS = tps
	}
}

// Run runs a game. Each tick Run calls both Update and Draw. Update updates
// the game state using the keypresses since the last tick. If Update returns
// nil, Run will continue execution by calling Draw. If Update returns
// Termination, Run will terminate without error. All other non-nil errors
// result in Run terminating with an error.
func Run(g Game, opts ...RunOption) error {
	// Setup config.
	config := DefaultRunConfig()
	for _, opt := range opts {
		opt(config)
	}

	// Setup term.
	term := config.Terminal
	if err := term.Init(); err != nil {
		return err
	}
	input := term.Input()
	defer term.Done()

	// Setup a ticker to trigger Update and Draw.
	ticker := time.NewTicker(time.Second / time.Duration(config.TPS))
	defer ticker.Stop()

	// Slice to accumulate keypresses each tick.
	var keys []Key

	// Run the actual game loop.
	for {
		select {
		case <-ticker.C:
			// Each tick, run Update and then either Draw or terminate.
			switch err := g.Update(keys); err {
			case nil:
				// Reset keys (in-place to avoid allocs) now they've been used.
				keys = keys[:0]
				term.Clear()
				g.Draw(term)
				if err := term.Flush(); err != nil {
					return err
				}
			case Termination:
				// Termination indicates normal termination, so return nil.
				return nil
			default:
				// All other errors are actual errors, so return them.
				return err
			}
		case key := <-input:
			// We could append to keys in the polling goroutine, but doing
			// inside the select coordinates access to keys with the ticker.
			keys = append(keys, key)
		}
	}
}

// Mob is a game object which occupies Tile.
type Mob[T any] struct {
	Face Glyph
	Pos  *Tile[T]

	Data T

	OnCollide func(*Mob[T], *Tile[T])
	OnBump    func(*Mob[T], *Mob[T])
}

// NewMob constructs a Mob with the given Glyph face.
func NewMob[T any](face Glyph, data T) *Mob[T] {
	return &Mob[T]{
		Face: face,
		Data: data,
	}
}

// Move attempts to move the Mob to a new Tile.
func (m *Mob[T]) Move(delta Vector) {
	dst := m.Pos.Adjacent[delta]
	if !dst.Pass {
		m.OnCollide(m, dst)
	} else if dst.Occupant != nil {
		m.OnBump(m, dst.Occupant)
	} else {
		m.Pos.Occupant = nil
		dst.Occupant = m
		m.Pos = dst
	}
}

// Tile is a square in the game map.
type Tile[T any] struct {
	Face     Glyph
	Pass     bool
	Occupant *Mob[T]

	Offset   Vector
	Adjacent map[Vector]*Tile[T]
}

// NewTile creates a new Tile with the given Vector offset.
func NewTile[T any](offset Vector) *Tile[T] {
	return &Tile[T]{
		Face:     Ch('.'),
		Pass:     true,
		Offset:   offset,
		Adjacent: make(map[Vector]*Tile[T]),
	}
}

// PlaceMob places a Mob on a Tile.
func PlaceMob[T any](m *Mob[T], t *Tile[T]) {
	m.Pos = t
	t.Occupant = m
}
