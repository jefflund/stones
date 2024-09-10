// Package tui provides utilities for building terminal user interfaces.
package tui

import (
	"github.com/jefflund/stones/pkg/hjkl"
	"github.com/jefflund/stones/pkg/hjkl/rl"
)

// DrawBorder draws the outline of a rectangle.
func DrawBorder(c hjkl.Canvas, pos, size hjkl.Vector) {
	bound := pos.Add(size).Sub(hjkl.Vec(1, 1))
	for x := pos.X + 1; x <= bound.X-1; x++ {
		c.Blit(hjkl.Vec(x, pos.Y), hjkl.Ch('-'))
		c.Blit(hjkl.Vec(x, bound.Y), hjkl.Ch('-'))
	}
	for y := pos.Y + 1; y <= bound.Y-1; y++ {
		c.Blit(hjkl.Vec(pos.X, y), hjkl.Ch('|'))
		c.Blit(hjkl.Vec(bound.X, y), hjkl.Ch('|'))
	}
	c.Blit(pos, hjkl.Ch('+'))
	c.Blit(hjkl.Vec(bound.X, pos.Y), hjkl.Ch('+'))
	c.Blit(hjkl.Vec(pos.X, bound.Y), hjkl.Ch('+'))
	c.Blit(bound, hjkl.Ch('+'))
}

// DrawTiles draws a set of Tile at a given Vector offset.
func DrawTiles(c hjkl.Canvas, offset hjkl.Vector, tiles []*rl.Tile) {
	for _, t := range tiles {
		f := t.Face
		if t.Occupant != nil {
			f = t.Occupant.Face
		}
		c.Blit(t.Offset.Add(offset), f)
	}
}
