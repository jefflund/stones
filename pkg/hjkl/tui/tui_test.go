package tui

import (
	"reflect"
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
		t.Error("BorderWidget.Draw produced incorrect buffer")
	}
}

func TestTilesWidget(t *testing.T) {
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
		t.Error("TilesWidget.Draw produced incorrect buffer", c)
	}
}

func TestLogWidget_Update(t *testing.T) {
	cases := []struct {
		update   string
		expected []string
	}{
		{"abc", []string{"abc"}},
		{"abc", []string{"abc", "abc"}},
		{"asdf", []string{"abc", "abc", "asdf"}},
		{"hjkl", []string{"abc", "asdf", "hjkl"}},
		{"abc", []string{"asdf", "hjkl", "abc"}},
	}
	w := NewLog(hjkl.Vec(0, 0), hjkl.Vec(3, 3))
	for _, c := range cases {
		w.Update(c.update)
		if !reflect.DeepEqual(w.History, c.expected) {
			t.Fatal("LogWidget.Update produced incorrect history")
		}
	}
}

func TestLogWidget_Draw(t *testing.T) {
	c := MockCanvas{}
	w := &LogWidget{
		Window: Window{hjkl.Vec(2, 1), hjkl.Vec(3, 3)},
		History: []string{
			"rl",
			"abc",
			"hjkl",
		},
	}
	w.Draw(c)
	expected := []string{
		"        ",
		"  rl    ",
		"  abc   ",
		"  hjk   ",
		"        ",
	}
	if !c.Equals(expected) {
		t.Error("LogWidget.Draw produced incorrect buffer", c)
	}
}
