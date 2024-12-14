package main

import (
	"math/rand"

	"github.com/jefflund/stones/pkg/hjkl"
)

type Game struct {
	Hero  *hjkl.Mob
	Level []*hjkl.Tile
}

func NewGame() *Game {
	hero := hjkl.NewMob(hjkl.Ch('@'))
	level := hjkl.GenTileGrid(80, 24, func(o hjkl.Vector) *hjkl.Tile {
		t := hjkl.NewTile(o)
		if rand.Float64() < 0.1 {
			t.Face = hjkl.ChFg('%', hjkl.ColorGreen)
			t.Pass = false
		}
		return t
	})
	hjkl.GenFence(level, func(t *hjkl.Tile) {
		t.Face = hjkl.Ch('#')
		t.Pass = false
	})
	hero.Pos = level[0]
	hero.Pos.Occupant = hero
	return &Game{hero, level}
}

func (g *Game) Update(ks []hjkl.Key) error {
	for _, k := range ks {
		switch k {
		case hjkl.KeyEsc, hjkl.KeyCtrlC:
			return hjkl.Termination
		default:
			if delta, ok := hjkl.VIKeyDirs[k]; ok {
				if dst := g.Hero.Pos.Adjacent[delta]; ok && dst.Pass {
					g.Hero.Pos.Occupant = nil
					dst.Occupant = g.Hero
					g.Hero.Pos = dst
				}
			}
		}
	}
	return nil
}

func (g *Game) Draw(c hjkl.Canvas) {
	for _, g := range g.Level {
		c.Blit(g.Offset, g.Face)
	}
	c.Blit(g.Hero.Pos.Offset, g.Hero.Face)
}

func main() {
	if err := hjkl.Run(NewGame()); err != nil {
		panic(err)
	}
}
