package main

import (
	"github.com/jefflund/stones/pkg/habilis"
	"github.com/jefflund/stones/pkg/hjkl"
	"github.com/jefflund/stones/pkg/hjkl/math/rand"
	"github.com/jefflund/stones/pkg/hjkl/rl"
	"github.com/jefflund/stones/pkg/hjkl/rl/gen"
	"github.com/jefflund/stones/pkg/hjkl/tui"
)

type Game struct {
	tui.TUI
	Hero  *rl.Mob
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
			t.Face = hjkl.Glyph{Ch: ch, Fg: hjkl.ColorGreen}
			t.Pass = false
		}
		return t
	})
	open := func(t *rl.Tile) bool {
		return t.Pass && t.Occupant == nil
	}

	hero := habilis.NewHero()
	rl.PlaceMob(hero, rand.Select(tiles, open))

	for i := 0; i < 20; i++ {
		mob := rand.Choice(habilis.Bestiary).New()
		rl.PlaceMob(mob, rand.Select(tiles, open))
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

	return &Game{screen, hero, tiles}
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

func main() {
	if err := hjkl.Run(NewGame()); err != nil {
		panic(err)
	}
}
