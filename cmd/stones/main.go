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
	'n': hjkl.Vec(1, 1),
	'b': hjkl.Vec(-1, 1),
	'u': hjkl.Vec(1, -1),
	'y': hjkl.Vec(-1, -1),
}

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

type Game struct {
	Hero  *hjkl.Mob[habilis.Skin]
	Tiles []*hjkl.Tile[habilis.Skin]
}

func NewGame() *Game {
	tiles := hjkl.GenTileGrid(80, 24, func(o hjkl.Vector) *hjkl.Tile[habilis.Skin] {
		t := hjkl.NewTile[habilis.Skin](o)
		if hjkl.RandChance(.1) {
			t.Pass = false
			t.Face = hjkl.RandChoice(Tree)
		} else {
			t.Face = hjkl.RandChoice(Ground)
		}
		return t
	})
	hjkl.GenFence(tiles, func(t *hjkl.Tile[habilis.Skin]) {
		t.Face = hjkl.Ch('#')
		t.Pass = false
	})
	open := func(t *hjkl.Tile[habilis.Skin]) bool {
		return t.Pass && t.Occupant == nil
	}

	hero := habilis.NewSkinMob(
		"Grog",
		hjkl.Ch('@'),
		habilis.NewCircle("Core", habilis.StoneCore, 3),
		habilis.NewCircle("Rogok", habilis.StoneDmg, 1),
		habilis.NewCircle("Warrior", habilis.StoneMelee, 1),
		habilis.NewCircle("Tough", habilis.StoneArm, 1),
	)
	hjkl.PlaceMob(hero, hjkl.RandSelect(tiles, open))

	prey := habilis.NewSkinMob(
		"Mammoth",
		hjkl.Ch('M'),
		habilis.NewCircle("Core", habilis.StoneCore, 9),
		habilis.NewCircle("Tusks", habilis.StoneMelee, 1),
		habilis.NewCircle("Tough", habilis.StoneArm, 1),
	)
	hjkl.PlaceMob(prey, hjkl.RandSelect(tiles, open))

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
	hjkl.WithWindow(c, hjkl.Vec(0, 0), hjkl.Vec(80, 24), func(c hjkl.Canvas) {
		hjkl.DisplayTiles(c, g.Tiles)
	})
}

func main() {
	if err := hjkl.Run(NewGame()); err != nil {
		panic(err)
	}
}
