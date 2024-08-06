package habilis

import "github.com/jefflund/stones/pkg/hjkl"

var Ground = []hjkl.Glyph{
	{Ch: '.', Fg: hjkl.ColorGreen},
	{Ch: '.', Fg: hjkl.ColorLightGreen},
	{Ch: '.', Fg: hjkl.ColorLightGreen},
}

var Tree = []hjkl.Glyph{
	{Ch: '%', Fg: hjkl.ColorGreen},
	{Ch: '%', Fg: hjkl.ColorLightGreen},
	{Ch: '%', Fg: hjkl.ColorYellow},
	{Ch: '%', Fg: hjkl.ColorLightYellow},
}

func GenTile(o hjkl.Vector) *hjkl.Tile {
	t := hjkl.NewTile(o)
	if hjkl.RandChance(.1) {
		t.Pass = false
		t.Face = hjkl.RandChoice(Tree)
	} else {
		t.Face = hjkl.RandChoice(Ground)
	}
	return t
}

func ModFence(t *hjkl.Tile) {
	t.Pass = false
	t.Face = hjkl.Ch('#')
}

func GenTileGrid(cols, rows int) []*hjkl.Tile {
	return hjkl.GenTileGrid(cols, rows, GenTile)
}
