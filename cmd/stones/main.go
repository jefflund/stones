package main

import (
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
	tiles := gen.GenTileGrid(48, 22, func(o hjkl.Vector) *rl.Tile {
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

	hero := rl.NewMob(hjkl.Ch('@'))
	rl.PlaceMob(hero, rand.Select(tiles, open))

	log := tui.NewLog(hjkl.Vec(50, 1), hjkl.Vec(29, 22))
	hero.AddComponent(rl.ComponentFunc[rl.CollideEvent](func(m *rl.Mob, v *rl.CollideEvent) {
		log.Update(tui.Log(
			"%s <bump> %o",
			string(m.Face.Ch),
			string(v.Obstacle.Face.Ch)))
	}))

	screen := tui.TUI{
		tui.NewBorder(hjkl.Vec(0, 0), hjkl.Vec(50, 24)),
		tui.NewTiles(hjkl.Vec(1, 1), hjkl.Vec(48, 22), tiles),
		tui.NewBorder(hjkl.Vec(49, 0), hjkl.Vec(31, 24)),
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
