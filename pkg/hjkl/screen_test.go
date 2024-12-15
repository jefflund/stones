package hjkl

import "testing"

type MockCanvas map[Vector]Glyph

func (c MockCanvas) Blit(v Vector, g Glyph) {
	c[v] = g
}

func (c MockCanvas) Equals(expected []string) bool {
	n := 0
	for y, row := range expected {
		for x, ch := range row {
			if ch == ' ' {
				continue
			}
			if c[Vec(x, y)] != Ch(ch) {
				return false
			}
			n++
		}
	}
	return len(c) == n
}

func TestTilesWidget(t *testing.T) {
	c := make(MockCanvas)
	NewTestTile := func(x, y int, f Glyph, m *Mob) *Tile {
		t := NewTile(Vec(x, y))
		t.Face = f
		t.Occupant = m
		return t
	}
	tiles := []*Tile{
		NewTestTile(0, 0, Ch('#'), nil),
		NewTestTile(1, 0, Ch('#'), nil),
		NewTestTile(2, 0, Ch('#'), nil),
		NewTestTile(0, 1, Ch('.'), nil),
		NewTestTile(1, 1, Ch('.'), NewMob(Ch('@'))),
		NewTestTile(2, 1, Ch('.'), NewMob(Ch('D'))),
		NewTestTile(0, 2, Ch('#'), nil),
		NewTestTile(1, 2, Ch('+'), nil),
		NewTestTile(2, 2, Ch('#'), nil),
	}
	NewTilesWidget(Vec(2, 1), Vec(3, 3), tiles).Draw(c)
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
