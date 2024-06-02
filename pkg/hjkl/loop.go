// Package hjkl is a terminal game engine.
package hjkl

import (
	"errors"
	"time"

	"github.com/nsf/termbox-go"
)

// Key represents a single keypress.
type Key rune

// Key constants which normally require escapes.
const (
	KeyEsc   Key = 0x1B
	KeyEnter Key = 0x0D
	KeyCtrlC Key = 0x03
)

// Vector is a two-dimension int vector.
type Vector struct {
	X, Y int
}

// Vec is shorthand for Vector{X: x, Y: y}.
func Vec(x, y int) Vector {
	return Vector{X: x, Y: y}
}

// Color describes the color of a Glyph as a uint8 ANSI 256 color code.
type Color uint8

// Color constants for use with Glyph.
const (
	ColorBlack Color = iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
	ColorLightBlack
	ColorLightRed
	ColorLightGreen
	ColorLightYellow
	ColorLightBlue
	ColorLightMagenta
	ColorLightCyan
	ColorLightWhite
)

// Glyph represents a single onscreen character.
type Glyph struct {
	Ch rune
	Fg Color
	Bg Color
}

// Ch is shorthand for Glyph{Ch: ch, Fg: ColorWhite, Bg: ColorBlack}.
func Ch(ch rune) Glyph {
	return Glyph{Ch: ch, Fg: ColorWhite, Bg: ColorBlack}
}

// Canvas is used by Game to draw Glyph on screen.
type Canvas interface {
	Blit(Vector, Glyph)
}

// Game contains the functionality needed to run a game.
type Game interface {
	Update([]Key) error
	Draw(Canvas)
}

// Termination is a special error indicating normal game termination.
var Termination = errors.New("Termination")

func Run(g Game) error {
	// Setup termbox for use.
	if err := termbox.Init(); err != nil {
		return err
	}
	if mode := termbox.SetInputMode(termbox.InputEsc); mode != termbox.InputEsc {
		return errors.New("could not set input mode")
	}
	if mode := termbox.SetOutputMode(termbox.Output256); mode != termbox.Output256 {
		return errors.New("could not set output mode")
	}
	defer termbox.Close()

	// Setup a goroutine to constantly poll for events.
	events := make(chan termbox.Event)
	go func() {
		defer close(events)
		for {
			event := termbox.PollEvent()
			if event.Type == termbox.EventInterrupt {
				return
			}
			events <- event
		}
	}()
	defer termbox.Interrupt() // Causes event polling to end gracefully.

	// Setup a ticker to trigger Update and Draw.
	ticker := time.NewTicker(time.Second / 20)
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
				// Use CellBuffer rather than Clear to avoid extra Flush.
				cells := termbox.CellBuffer()
				for i := 0; i < len(cells); i++ {
					cells[i].Ch = 0
					cells[i].Fg = 0
					cells[i].Bg = 0
				}
				g.Draw(termboxCanvas{})
				if err := termbox.Flush(); err != nil {
					return err
				}
			case Termination:
				// Termination indicates normal termination, so return nil.
				return nil
			default:
				// All other errors are actual errors, so return them.
				return err
			}
		case event := <-events:
			// We could append to keys from inside the polling goroutine, but
			// if we recieve from inside a select, we can coordinate access to
			// the keys slice with the ticker.
			if event.Type == termbox.EventKey {
				keys = append(keys, Key(event.Ch)|Key(event.Key))
			}
		}
	}
}

// termboxCanvas implements Canvas using termbox.SetCell.
type termboxCanvas struct{}

// Blit place a Glyph into the screen buffer.
func (termboxCanvas) Blit(v Vector, g Glyph) {
	termbox.SetCell(v.X, v.Y, g.Ch, termbox.Attribute(g.Fg+1), termbox.Attribute(g.Bg+1))
}
