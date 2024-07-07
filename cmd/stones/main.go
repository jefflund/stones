package main

import (
	"github.com/jefflund/stones/pkg/habilis"
	"github.com/jefflund/stones/pkg/hjkl"
)

var KeyMap = map[hjkl.Key]hjkl.Vector{
	'h': hjkl.Vec(-1, 0),
	'j': hjkl.Vec(0, 1),
	'k': hjkl.Vec(0, -1),
	'l': hjkl.Vec(1, 0),
}

type Game struct {
	Hero  *hjkl.Mob[habilis.Skin]
	Tiles []*hjkl.Tile[habilis.Skin]
}

func NewGame() *Game {
	tiles := hjkl.GenTileGrid(80, 24, hjkl.NewTile[habilis.Skin])
	hjkl.GenFence(tiles, func(t *hjkl.Tile[habilis.Skin]) {
		t.Face = hjkl.Ch('#')
		t.Pass = false
	})

	hero := habilis.NewSkinMob(
		"Grog",
		hjkl.Ch('@'),
		habilis.NewCircle("Core", habilis.StoneCore, 3),
		habilis.NewCircle("Rogok", habilis.StoneDmg, 1),
		habilis.NewCircle("Warrior", habilis.StoneMelee, 1),
		habilis.NewCircle("Tough", habilis.StoneArm, 1),
	)
	hjkl.PlaceMob(hero, tiles[1000])

	prey := habilis.NewSkinMob(
		"Mammoth",
		hjkl.Ch('M'),
		habilis.NewCircle("Core", habilis.StoneCore, 9),
		habilis.NewCircle("Tusks", habilis.StoneMelee, 1),
		habilis.NewCircle("Tough", habilis.StoneArm, 1),
	)
	hjkl.PlaceMob(prey, tiles[500])

	return &Game{hero, tiles}
}

func (g *Game) Update(ks []hjkl.Key) error {
	for _, k := range ks {
		if k == 'q' || k == hjkl.KeyEsc {
			return hjkl.Termination
		}
		if delta, ok := KeyMap[k]; ok {
			g.Hero.Move(delta)
		}
	}
	return nil
}

func (g *Game) Draw(c hjkl.Canvas) {
	for _, t := range g.Tiles {
		if t.Occupant != nil {
			c.Blit(t.Offset, t.Occupant.Face)
		} else {
			c.Blit(t.Offset, t.Face)
		}
	}
}

func main() {
	if err := hjkl.Run(NewGame()); err != nil {
		panic(err)
	}
}
