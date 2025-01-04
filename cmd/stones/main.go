package main

import (
	"github.com/jefflund/stones/pkg/hjkl"
	"github.com/jefflund/stones/pkg/hjkl/clock"
	"github.com/jefflund/stones/pkg/hjkl/gen"
	"github.com/jefflund/stones/pkg/hjkl/rand"
	"github.com/jefflund/stones/pkg/rpg"
)

type Game struct {
	hjkl.Screen
	Hero  *hjkl.Mob
	Level []*hjkl.Tile
	Clock *clock.Clock[*hjkl.Mob]
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

	clock := clock.New[*hjkl.Mob]()

	hero := rpg.NewHero()
	hjkl.PlaceMob(hero, rand.FilteredChoice(level, open))

	for i := 1; i <= 30; i++ {
		mob := rand.Choice(rpg.Bestiary).New()
		hjkl.PlaceMob(mob, rand.FilteredChoice(level, open))
		clock.Schedule(mob, i%10)
	}

	screen := hjkl.Screen{
		hjkl.NewTilesWidget(hjkl.Vec(0, 0), hjkl.Vec(cols, rows), level),
	}

	return &Game{screen, hero, level, clock}
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

	for _, m := range g.Clock.Tick() {
		if m.Pos == nil {
			continue
		}

		delta := rand.Choice(hjkl.CompassDirs)
		m.Handle(&hjkl.Move{Delta: delta})
		g.Clock.Schedule(m, rand.Range(10, 50))
	}

	return nil
}

func main() {
	if err := hjkl.Run(NewGame()); err != nil {
		panic(err)
	}
}
