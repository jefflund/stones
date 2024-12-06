package main

import (
	"github.com/jefflund/stones/pkg/habilis"
	"github.com/jefflund/stones/pkg/hjkl"
	"github.com/jefflund/stones/pkg/hjkl/math/rand"
	"github.com/jefflund/stones/pkg/hjkl/rl"
	"github.com/jefflund/stones/pkg/hjkl/rl/gen"
	"github.com/jefflund/stones/pkg/hjkl/tui"
)

type TriggerWander struct{}

type Game struct {
	tui.TUI
	Hero  *rl.Mob
	Mobs  []*rl.Mob
	Tiles []*rl.Tile
}

func NewGame() *Game {
	const (
		SCR_COLS = 80
		MAP_COLS = 45
		LOG_COLS = SCR_COLS - MAP_COLS - 3
		SCR_ROWS = 24
		MAP_ROWS = SCR_ROWS - 2
		LOG_ROWS = SCR_ROWS - 2
	)

	tiles := gen.GenTileGrid(MAP_COLS, MAP_ROWS, func(o hjkl.Vector) *rl.Tile {
		t := rl.NewTile(o)
		if rand.Chance(0.1) {
			ch := rand.Choice([]rune{'%', '%', '&', '|'})
			t.Face = hjkl.ChFg(ch, hjkl.ColorGreen)
			t.Pass = false
		}
		return t
	})
	gen.GenFence(tiles, func(t *rl.Tile) {
		t.Pass = false
		t.Face = hjkl.ChFg('#', hjkl.ColorLightBlack)
	})
	open := func(t *rl.Tile) bool {
		return t.Pass && t.Occupant == nil
	}

	hero := habilis.NewHero()
	rl.PlaceMob(hero, rand.Select(tiles, open))

	mobs := make([]*rl.Mob, 20)
	for i := 0; i < len(mobs); i++ {
		mobs[i] = rand.Choice(habilis.Bestiary).New()
		rl.PlaceMob(mobs[i], rand.Select(tiles, open))
	}

	log := tui.NewLog(hjkl.Vec(MAP_COLS+2, 1), hjkl.Vec(LOG_COLS, LOG_ROWS))
	hero.AddComponent(rl.ComponentFunc(func(m *rl.Mob, v rl.Event) {
		switch v := v.(type) {
		case *tui.LogEvent:
			log.Update(v.Message)
		}
	}))

	screen := tui.TUI{
		tui.NewBorder(hjkl.Vec(0, 0), hjkl.Vec(SCR_COLS, SCR_ROWS)),
		tui.NewBorder(hjkl.Vec(MAP_COLS+1, 0), hjkl.Vec(1, SCR_ROWS)),
		tui.NewTiles(hjkl.Vec(1, 1), hjkl.Vec(MAP_COLS, MAP_ROWS), tiles),
		log,
	}

	return &Game{screen, hero, mobs, tiles}
}

func (g *Game) Update(ks []hjkl.Key) error {
	for _, k := range ks {
		if k == 'q' || k == hjkl.KeyEsc {
			return hjkl.Termination
		}
	}

	trigger := &habilis.ActTrigger{Keys: ks}
	g.Hero.Handle(trigger)
	if len(ks) > 0 {
		for _, m := range g.Mobs {
			m.Handle(trigger)
		}
	}
	return nil
}

func main() {
	if err := hjkl.Run(NewGame()); err != nil {
		panic(err)
	}
}
