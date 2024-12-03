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
	// TODO: Add computed constants for the various widget dimensions.
	tiles := gen.GenTileGrid(43, 22, func(o hjkl.Vector) *rl.Tile {
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

	log := tui.NewLog(hjkl.Vec(45, 1), hjkl.Vec(34, 22))
	hero.AddComponent(rl.ComponentFunc(func(m *rl.Mob, v rl.Event) {
		switch v := v.(type) {
		case *tui.LogEvent:
			log.Update(v.Message)
		}
	}))

	screen := tui.TUI{
		tui.NewBorder(hjkl.Vec(0, 0), hjkl.Vec(45, 24)),
		tui.NewTiles(hjkl.Vec(1, 1), hjkl.Vec(43, 22), tiles),
		tui.NewBorder(hjkl.Vec(44, 0), hjkl.Vec(36, 24)),
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
