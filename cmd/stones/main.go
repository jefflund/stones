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
	hjkl.PlaceMob(hero, level[0])
	hjkl.PlaceMob(hjkl.NewMob(hjkl.ChFg('u', hjkl.ColorRed)), level[0])
	hjkl.PlaceMob(hjkl.NewMob(hjkl.ChFg('u', hjkl.ColorRed)), level[1])
	hjkl.PlaceMob(hjkl.NewMob(hjkl.ChFg('u', hjkl.ColorRed)), level[2])
	hjkl.PlaceMob(hjkl.NewMob(hjkl.ChFg('U', hjkl.ColorLightRed)), level[3])
	return &Game{hero, level}
}

func (g *Game) Update(ks []hjkl.Key) error {
	for _, k := range ks {
		switch k {
		case hjkl.KeyEsc, hjkl.KeyCtrlC:
			return hjkl.Termination
		default:
			if delta, ok := hjkl.VIKeyDirs[k]; ok {
				g.Hero.Handle(&hjkl.Move{Delta: delta})
			}
		}
	}
	return nil
}

func (g *Game) Draw(c hjkl.Canvas) {
	for _, t := range g.Level {
		c.Blit(t.Offset, hjkl.Get(t, &hjkl.Face{}))
	}
}

func main() {
	if err := hjkl.Run(NewGame()); err != nil {
		panic(err)
	}
}
