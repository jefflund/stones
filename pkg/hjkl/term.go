package hjkl

import (
	"errors"

	"github.com/nsf/termbox-go"
)

// Canvas blits Glyph.
type Canvas interface {
	Blit(Vector, Glyph)
}

// Terminal provides I/O capabilities.
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

// Blit places a Glyph into the termbox cell buffer.
func (TermboxTerminal) Blit(v Vector, g Glyph) {
	// The offsets in the Attribute conversions are because termbox shifts the
	// ANSI color codes, while hjkl prefers the standard ANSI color codes.
	termbox.SetCell(v.X, v.Y, g.Ch, termbox.Attribute(g.Fg+1), termbox.Attribute(g.Bg+1))
}

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
	input := make(chan Key)
	go func() {
		defer close(input)
		for {
			switch event := termbox.PollEvent(); event.Type {
			case termbox.EventInterrupt:
				// Interupt is called by Done.
				return
			case termbox.EventKey:
				// Unlike termbox, hjkl makes no distinctio between character
				// keys and other special keys, so union the two.
				input <- Key(event.Ch) | Key(event.Key)
			}
		}
	}()
	return input
}

// Done interupts the input goroutine and closes termbox.
func (TermboxTerminal) Done() {
	termbox.Interrupt()
	termbox.Close()
}

// Clear clears the termbox cell buffer.
func (TermboxTerminal) Clear() {
	// Use CellBuffer instead of termbox.Clear to avoid extra Flush.
	cells := termbox.CellBuffer()
	for i := 0; i < len(cells); i++ {
		cells[i] = termbox.Cell{}
	}
}

// Flush ensures the screen reflects the state of the termbox cell buffer.
func (TermboxTerminal) Flush() error {
	return termbox.Flush()
}
