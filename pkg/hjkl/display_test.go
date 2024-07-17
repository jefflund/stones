package hjkl

import (
	"reflect"
	"testing"
)

type MockCanvas struct {
	Buffer map[Vector]Glyph
}

func (c *MockCanvas) Blit(v Vector, g Glyph) {
	if c.Buffer == nil {
		c.Buffer = make(map[Vector]Glyph)
	}
	c.Buffer[v] = g
}

func TestWithWindow(t *testing.T) {
	c := &MockCanvas{}
	WithWindow(c, Vec(3, 4), Vec(5, 6), func(c Canvas) {
		c.Blit(Vec(-1, 0), Ch('#'))
		c.Blit(Vec(0, -1), Ch('#'))
		c.Blit(Vec(5, 3), Ch('#'))
		c.Blit(Vec(2, 7), Ch('#'))

		c.Blit(Vec(0, 0), Ch('#'))
		c.Blit(Vec(1, 2), Ch('@'))
		c.Blit(Vec(4, 5), Ch('M'))
	})
	expected := map[Vector]Glyph{
		Vec(3, 4): Ch('#'),
		Vec(4, 6): Ch('@'),
		Vec(7, 9): Ch('M'),
	}
	if !reflect.DeepEqual(c.Buffer, expected) {
		t.Error("WithWidow gave incorrect Canvas state")
	}
}

func TestDisplayTiles(t *testing.T) {
	c := &MockCanvas{}
	m := NewMob[any](Ch('@'), nil)
	tiles := []*Tile[any]{
		{Face: Ch('#'), Offset: Vec(0, 1)},
		{Face: Ch('.'), Offset: Vec(2, 3)},
		{Face: Ch('.'), Offset: Vec(4, 5), Occupant: m},
	}
	DisplayTiles(c, tiles)
	expected := map[Vector]Glyph{
		Vec(0, 1): Ch('#'),
		Vec(2, 3): Ch('.'),
		Vec(4, 5): Ch('@'),
	}
	if !reflect.DeepEqual(c.Buffer, expected) {
		t.Error("DisplayTiles gave incorrect Canvas state")
	}
}

func TestDisplayBorder(t *testing.T) {
	c := &MockCanvas{}
	DisplayBorder(c, 4, 3)
	expected := map[Vector]Glyph{
		Vec(0, 0): Ch('+'),
		Vec(1, 0): Ch('-'),
		Vec(2, 0): Ch('-'),
		Vec(3, 0): Ch('+'),
		Vec(0, 1): Ch('|'),
		Vec(3, 1): Ch('|'),
		Vec(0, 2): Ch('+'),
		Vec(1, 2): Ch('-'),
		Vec(2, 2): Ch('-'),
		Vec(3, 2): Ch('+'),
	}
	if !reflect.DeepEqual(c.Buffer, expected) {
		t.Error("DisplayBorder gave incorrect Canvas state")
	}
}
