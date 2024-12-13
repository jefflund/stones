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

// Run runs a game. Each tick Run calls both Update and Draw. Update updates
// the game state using the keypresses since the last tick. If Update returns
// nil, Run will continue execution by calling Draw. If Update returns
// Termination, Run will terminate without error. All other non-nil errors
// result in Run terminating with an error.
func Run(g Game) error {
	// Setup I/O via Terminal.
	term := TermboxTerminal{}
	if err := term.Init(); err != nil {
		return err
	}
	input := term.Input()
	defer term.Done()

	// Setup a ticker to trigger Update and Draw.
	ticker := time.NewTicker(time.Second / time.Duration(20))
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
				g.Draw(term) // Game only gets Canvas part of Terminal.
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
