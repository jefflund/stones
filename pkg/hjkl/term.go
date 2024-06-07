package hjkl

import (
	"errors"

	"github.com/nsf/termbox-go"
)

// Canvas is used by Game to draw Glyph on screen.
type Canvas interface {
	Blit(Vector, Glyph)
}

// Terminal provides I/O capabilities for a Game.
type Terminal interface {
	Canvas

	Init() error
	Input() chan Key
	Done()

	Clear()
	Flush() error
}

// TermboxTerminal is a Terminal implemented using termbox.
type TermboxTerminal struct{}

// Init initializes termbox and sets the input and output modes.
func (TermboxTerminal) Init() error {
	if err := termbox.Init(); err != nil {
		return err
	}
	if mode := termbox.SetInputMode(termbox.InputEsc); mode != termbox.InputEsc {
		return errors.New("could not set input mode")
	}
	if mode := termbox.SetOutputMode(termbox.Output256); mode != termbox.Output256 {
		return errors.New("could not set output mode")
	}

	return nil
}

// Input gets an input channel and starts a goroutine to fill it.
func (TermboxTerminal) Input() chan Key {
	keys := make(chan Key)
	go func() {
		defer close(keys)
		for {
			switch event := termbox.PollEvent(); event.Type {
			case termbox.EventInterrupt:
				return
			case termbox.EventKey:
				keys <- Key(event.Ch) | Key(event.Key)
			}
		}
	}()
	return keys
}

// Done interupts the event polling goroutine and closes termbox.
func (t *TermboxTerminal) Done() {
	termbox.Interrupt()
	termbox.Close()
}

// Clear clears the termbox cell buffer.
func (TermboxTerminal) Clear() {
	cells := termbox.CellBuffer()
	for i := 0; i < len(cells); i++ {
		cells[i].Ch = 0
		cells[i].Fg = 0
		cells[i].Bg = 0
	}
}

// Blit places a Glyph into the termbox cell buffer.
func (TermboxTerminal) Blit(v Vector, g Glyph) {
	// The offsets in the Fg and Bg Attribute conversions are because termbox
	// shifts the ANSI color codes up by one so zero values can be defaults,
	// while hjkl prefers to use the standard code values.
	termbox.SetCell(v.X, v.Y, g.Ch, termbox.Attribute(g.Fg+1), termbox.Attribute(g.Bg+1))
}

// Flush ensures the screen reflects the state of the termbox cell buffer.
func (TermboxTerminal) Flush() error {
	return termbox.Flush()
}
