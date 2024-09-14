package tui

import (
	"testing"

	"github.com/jefflund/stones/pkg/hjkl"
	"github.com/jefflund/stones/pkg/hjkl/rl"
)

type MockCanvas map[hjkl.Vector]hjkl.Glyph

func (c MockCanvas) Blit(v hjkl.Vector, g hjkl.Glyph) {
	c[v] = g
}

func (c MockCanvas) Equals(expected []string) bool {
	n := 0
	for y, row := range expected {
		for x, ch := range row {
			if ch == ' ' {
				continue
			}
			if c[hjkl.Vec(x, y)] != hjkl.Ch(ch) {
				return false
			}
			n++
		}
	}
	return len(c) == n
}

func TestBorderWidget(t *testing.T) {
	c := make(MockCanvas)
	NewBorder(hjkl.Vec(1, 2), hjkl.Vec(5, 3)).Draw(c)
	expected := []string{
		"           ",
		"           ",
		" +---+     ",
		" |   |     ",
		" +---+     ",
		"           ",
	}
	if !c.Equals(expected) {
		t.Error("DrawBorder produced incorrect buffer")
	}
}

func TestDrawTiles(t *testing.T) {
	c := make(MockCanvas)
	asdf := func(x, y int, f, o rune) *rl.Tile {
		t := rl.NewTile(hjkl.Vec(x, y))
		t.Face = hjkl.Ch(f)
		if o != 0 {
			t.Occupant = rl.NewMob(hjkl.Ch(o))
		}
		return t
	}
	tiles := []*rl.Tile{
		asdf(0, 0, '#', 0),
		asdf(1, 0, '#', 0),
		asdf(2, 0, '#', 0),
		asdf(0, 1, '.', 0),
		asdf(1, 1, '.', '@'),
		asdf(2, 1, '.', 'D'),
		asdf(0, 2, '#', 0),
		asdf(1, 2, '+', 0),
		asdf(2, 2, '#', 0),
	}
	NewTiles(hjkl.Vec(2, 1), hjkl.Vec(3, 3), tiles).Draw(c)
	expected := []string{
		"        ",
		"  ###   ",
		"  .@D   ",
		"  #+#   ",
		"        ",
	}
	if !c.Equals(expected) {
		t.Error("DrawTilesproduced incorrect buffer", c)
	}
}
