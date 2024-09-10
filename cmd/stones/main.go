package main

import (
	"github.com/jefflund/stones/pkg/hjkl"
	"github.com/jefflund/stones/pkg/hjkl/math/rand"
	"github.com/jefflund/stones/pkg/hjkl/rl"
	"github.com/jefflund/stones/pkg/hjkl/rl/gen"
	"github.com/jefflund/stones/pkg/hjkl/tui"
)

type Game struct {
	Hero  *rl.Mob
	Tiles []*rl.Tile
}

func NewGame() *Game {
	tiles := gen.GenTileGrid(78, 22, func(o hjkl.Vector) *rl.Tile {
		t := rl.NewTile(o)
		if rand.Chance(0.1) {
			t.Face = hjkl.Glyph{Ch: '%', Fg: hjkl.ColorGreen}
			t.Pass = false
		}
		return t
	})
	open := func(t *rl.Tile) bool {
		return t.Pass && t.Occupant == nil
	}

	hero := rl.NewMob(hjkl.Ch('@'))
	rl.PlaceMob(hero, rand.Select(tiles, open))

	return &Game{hero, tiles}
}

func (g *Game) Update(ks []hjkl.Key) error {
	for _, k := range ks {
		if k == 'q' || k == hjkl.KeyEsc {
			return hjkl.Termination
		}
		if delta, ok := hjkl.VIKeyMap[k]; ok {
			if g.Hero.Pos.Adjacent[delta] != nil {
				g.Hero.Move(delta)
			}
		}
	}
	return nil
}

func (g *Game) Draw(c hjkl.Canvas) {
	tui.DrawBorder(c, hjkl.Vec(0, 0), hjkl.Vec(80, 24))
	tui.DrawTiles(c, hjkl.Vec(1, 1), g.Tiles)
}

func main() {
	if err := hjkl.Run(NewGame()); err != nil {
		panic(err)
	}
}
