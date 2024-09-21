// Package tui provides utilities for building terminal user interfaces.
package tui

import (
	"github.com/jefflund/stones/pkg/hjkl"
	"github.com/jefflund/stones/pkg/hjkl/rl"
)

// Widget is something which can draw itself on a Canvas.
type Widget interface {
	Draw(hjkl.Canvas)
}

// TUI is a slice of Widget.
type TUI []Widget

// Draw has each constituent Widget draw itself.
func (t TUI) Draw(c hjkl.Canvas) {
	for _, w := range t {
		w.Draw(c)
	}
}

// Window serves as a base for various Widget which need relative drawing.
type Window struct {
	Pos  hjkl.Vector
	Size hjkl.Vector
}

// Blit performs a Blit on the canvas relative to the location of the Window.
func (w *Window) Blit(c hjkl.Canvas, v hjkl.Vector, g hjkl.Glyph) {
	if 0 <= v.X && v.X < w.Size.X && 0 <= v.Y && v.Y < w.Size.Y {
		c.Blit(w.Pos.Add(v), g)
	}
}

// BorderWidget is a Widget which draws a border at its bounds.
type BorderWidget struct {
	Window
}

// NewBorder creates a BorderWidget with the given bounds.
func NewBorder(pos, size hjkl.Vector) *BorderWidget {
	return &BorderWidget{Window{pos, size}}
}

// Draw draws a border at the BorderWidget bounds.
func (w *BorderWidget) Draw(c hjkl.Canvas) {
	maxX, maxY := w.Size.X-1, w.Size.Y-1
	for x := 1; x < maxX; x++ {
		w.Blit(c, hjkl.Vec(x, 0), hjkl.Ch('-'))
		w.Blit(c, hjkl.Vec(x, maxY), hjkl.Ch('-'))
	}
	for y := 1; y <= maxY; y++ {
		w.Blit(c, hjkl.Vec(0, y), hjkl.Ch('|'))
		w.Blit(c, hjkl.Vec(maxX, y), hjkl.Ch('|'))
	}
	w.Blit(c, hjkl.Vec(0, 0), hjkl.Ch('+'))
	w.Blit(c, hjkl.Vec(maxX, 0), hjkl.Ch('+'))
	w.Blit(c, hjkl.Vec(0, maxY), hjkl.Ch('+'))
	w.Blit(c, hjkl.Vec(maxX, maxY), hjkl.Ch('+'))
}

// TilesWidget is a Widget which draws a collection of Tile.
type TilesWidget struct {
	Window
	Tiles []*rl.Tile
}

// NewTiles creates a TileWidget with the given collection of Tile.
func NewTiles(pos, size hjkl.Vector, tiles []*rl.Tile) *TilesWidget {
	return &TilesWidget{Window{pos, size}, tiles}
}

// Draw draws the collection of Tile within the TilesWidget bounds.
func (w *TilesWidget) Draw(c hjkl.Canvas) {
	for _, t := range w.Tiles {
		f := t.Face
		if t.Occupant != nil {
			f = t.Occupant.Face
		}
		w.Blit(c, t.Offset, f)
	}
}

// LogWidget is a Widget which draws log messages.
type LogWidget struct {
	Window
	History []string
}

// NewLog creates a LogWidget with an empty history.
func NewLog(pos, size hjkl.Vector) *LogWidget {
	return &LogWidget{
		Window: Window{pos, size},
	}
}

// Update appends a message to the history.
func (w *LogWidget) Update(msg string) {
	w.History = append(w.History, msg)
	if len(w.History) > w.Size.Y {
		w.History = w.History[len(w.History)-w.Size.Y:]
	}
}

// Draw displays the log history.
func (w *LogWidget) Draw(c hjkl.Canvas) {
	for y, msg := range w.History {
		for x, ch := range msg {
			w.Blit(c, hjkl.Vec(x, y), hjkl.Ch(ch))
		}
	}
}
