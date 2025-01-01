package main

import (
	"github.com/jefflund/stones/pkg/hjkl"
	"github.com/jefflund/stones/pkg/hjkl/gen"
	"github.com/jefflund/stones/pkg/hjkl/rand"
	"github.com/jefflund/stones/pkg/rpg"
)

type Game struct {
	hjkl.Screen
	Hero  *hjkl.Mob
	Level []*hjkl.Tile
}

func NewGame() *Game {
	cols, rows := 80, 24

	level := gen.GenTileGrid(cols, rows, func(o hjkl.Vector) *hjkl.Tile {
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
	})
	gen.GenFence(level, func(t *hjkl.Tile) {
		t.Face = hjkl.Ch('#')
		t.Pass = false
	})

	open := func(t *hjkl.Tile) bool {
		return t.Pass && t.Occupant == nil
	}

	hero := hjkl.NewMob(hjkl.Ch('@'))
	hero.Components.Add(&rpg.Stats{
		Health: 10,
		Damage: 2,
	})
	hjkl.PlaceMob(hero, rand.FilteredChoice(level, open))
	for i := 1; i <= 10; i++ {
		face := rand.Choice([]hjkl.Glyph{
			hjkl.ChFg('u', hjkl.ColorRed),
			hjkl.ChFg('U', hjkl.ColorRed),
			hjkl.ChFg('d', hjkl.ColorRed),
			hjkl.ChFg('D', hjkl.ColorRed),
		})
		mob := hjkl.NewMob(face)
		mob.Components.Add(&rpg.Stats{
			Health: i,
			Damage: 1,
		})
		hjkl.PlaceMob(mob, rand.FilteredChoice(level, open))
	}

	screen := hjkl.Screen{
		hjkl.NewTilesWidget(hjkl.Vec(0, 0), hjkl.Vec(cols, rows), level),
	}

	return &Game{screen, hero, level}
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

func main() {
	if err := hjkl.Run(NewGame()); err != nil {
		panic(err)
	}
}
