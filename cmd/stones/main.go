package main

import (
	"fmt"
	"strings"

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
	Hero  *hjkl.Mob
	Log   *hjkl.LogWidget
	Tiles []*hjkl.Tile
}

func NewGame() *Game {
	tiles := hjkl.GenTileGrid(60, 22, func(o hjkl.Vector) *hjkl.Tile {
		t := hjkl.NewTile(o)
		if hjkl.RandChance(.1) {
			t.Pass = false
			t.Face = hjkl.RandChoice(Tree)
		} else {
			t.Face = hjkl.RandChoice(Ground)
		}
		return t
	})
	open := func(t *hjkl.Tile) bool {
		return t.Pass && t.Occupant == nil
	}

	log := &hjkl.LogWidget{MaxLen: 5}

	hero := hjkl.NewMob(hjkl.Ch('@'))
	hero.AddComponent(&habilis.Skin{
		Name: "Grog",
		Circles: []habilis.Circle{
			habilis.NewCircle("Core", habilis.StoneCore, 3),
			habilis.NewCircle("Rogok", habilis.StoneDmg, 1),
			habilis.NewCircle("Warrior", habilis.StoneMelee, 1),
			habilis.NewCircle("Tough", habilis.StoneArm, 1),
		},
	})
	hero.AddComponent(log)
	hjkl.PlaceMob(hero, hjkl.RandSelect(tiles, open))

	hjkl.PlaceMob(
		habilis.NewBestiaryMob("Mammoth"),
		hjkl.RandSelect(tiles, open),
	)
	hjkl.PlaceMob(
		habilis.NewBestiaryMob("Sabertooth"),
		hjkl.RandSelect(tiles, open),
	)
	return &Game{hero, log, tiles}
}

func (g *Game) Update(ks []hjkl.Key) error {
	for _, k := range ks {
		if k == 'q' || k == hjkl.KeyEsc {
			return hjkl.Termination
		}
		if delta, ok := KeyMap[k]; ok {
			if g.Hero.Pos.Adjacent[delta] != nil {
				g.Hero.Move(delta)
			}
		}
	}
	return nil
}

func (g *Game) Status() string {
	s := habilis.GetSkin(g.Hero)
	lines := []string{
		s.Name,
		fmt.Sprintf("Stones: %d", s.Count(habilis.StoneAny)),
		fmt.Sprintf("Pos: %v", g.Hero.Pos.Offset),
	}
	return strings.Join(lines, "\n")
}

func (g *Game) Draw(c hjkl.Canvas) {
	hjkl.DisplayBorder(c, 80, 24)
	hjkl.WithWindow(c, hjkl.Vec(1, 1), hjkl.Vec(60, 22), func(c hjkl.Canvas) {
		hjkl.DisplayTiles(c, g.Tiles)
	})
	hjkl.WithWindow(c, hjkl.Vec(61, 1), hjkl.Vec(18, 22), func(c hjkl.Canvas) {
		hjkl.DisplayString(c, g.Status())
	})
	hjkl.WithWindow(c, hjkl.Vec(0, 24), hjkl.Vec(80, g.Log.MaxLen), g.Log.Display)
}

func main() {
	if err := hjkl.Run(NewGame()); err != nil {
		panic(err)
	}
}
