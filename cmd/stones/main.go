package main

import (
	"github.com/jefflund/stones/pkg/hjkl"
	"github.com/jefflund/stones/pkg/hjkl/math/rand"
)

type Game struct {
	Hero  *hjkl.Mob
	Tiles []*hjkl.Tile
}

func NewGame() *Game {
	tiles := hjkl.GenTileGrid(60, 22, func(o hjkl.Vector) *hjkl.Tile {
		t := hjkl.NewTile(o)
		if rand.Chance(0.1) {
			t.Face = hjkl.Glyph{Ch: '%', Fg: hjkl.ColorGreen}
			t.Pass = false
		}
		return t
	})
	open := func(t *hjkl.Tile) bool {
		return t.Pass && t.Occupant == nil
	}
	hjkl.GenFence(tiles, func(t *hjkl.Tile) {
		t.Face = hjkl.Ch('#')
		t.Pass = false
	})

	hero := hjkl.NewMob(hjkl.Ch('@'))
	hjkl.PlaceMob(hero, rand.Select(tiles, open))

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
	for _, t := range g.Tiles {
		c.Blit(t.Offset, t.Face)
		if t.Occupant != nil {
			c.Blit(t.Offset, t.Occupant.Face)
		}
	}
}

func main() {
	if err := hjkl.Run(NewGame()); err != nil {
		panic(err)
	}
}
