package rpg

import (
	"github.com/jefflund/stones/pkg/hjkl"
	"github.com/jefflund/stones/pkg/hjkl/rand"
)

func ForestTile(o hjkl.Vector) *hjkl.Tile {
	t := hjkl.NewTile(o)
	if rand.Chance(0.1) {
		t.Face = rand.Choice([]hjkl.Glyph{
			hjkl.ChFg('%', hjkl.ColorGreen),
			hjkl.ChFg('%', hjkl.ColorGreen),
			hjkl.ChFg('%', hjkl.ColorLightGreen),
			hjkl.ChFg('%', hjkl.ColorLightYellow),
		})
		t.Pass = false
	} else {
		t.Face = rand.Choice([]hjkl.Glyph{
			hjkl.ChFg('.', hjkl.ColorGreen),
			hjkl.ChFg('.', hjkl.ColorGreen),
			hjkl.ChFg('.', hjkl.ColorGreen),
			hjkl.ChFg('.', hjkl.ColorLightGreen),
			hjkl.ChFg('.', hjkl.ColorLightGreen),
			hjkl.ChFg('.', hjkl.ColorLightYellow),
			hjkl.ChFg('.', hjkl.ColorLightWhite),
		})
	}
	return t
}

func ForestFence(t *hjkl.Tile) {
	t.Face = hjkl.Ch('#')
	t.Pass = false
}
