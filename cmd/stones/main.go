package main

import "github.com/jefflund/stones/pkg/hjkl"

const (
	COLS = 80
	ROWS = 24
)

type Game struct {
	Hero  *hjkl.Mob
	Level map[hjkl.Vector]*hjkl.Tile
}

func NewGame() *Game {
	hero := hjkl.NewMob(hjkl.Ch('@'))
	level := make(map[hjkl.Vector]*hjkl.Tile)
	for x := 0; x < COLS; x++ {
		for y := 0; y < ROWS; y++ {
			tile := hjkl.NewTile(hjkl.Vec(x, y))
			if x == 0 || x == COLS-1 || y == 0 || y == ROWS-1 {
				tile.Face = hjkl.Ch('#')
				tile.Pass = false
			}
			level[tile.Offset] = tile
		}
	}
	hero.Pos = level[hjkl.Vec(40, 12)]
	hero.Pos.Occupant = hero
	return &Game{hero, level}
}

func (g *Game) Update(ks []hjkl.Key) error {
	var delta hjkl.Vector
	for _, k := range ks {
		switch k {
		case hjkl.KeyEsc, hjkl.KeyCtrlC:
			return hjkl.Termination
		case 'h':
			delta = hjkl.Vec(-1, 0)
		case 'j':
			delta = hjkl.Vec(0, 1)
		case 'k':
			delta = hjkl.Vec(0, -1)
		case 'l':
			delta = hjkl.Vec(1, 0)
		}
	}
	if dst := g.Level[g.Hero.Pos.Offset.Add(delta)]; dst.Pass {
		g.Hero.Pos.Occupant = nil
		dst.Occupant = g.Hero
		g.Hero.Pos = dst
	}
	return nil
}

func (g *Game) Draw(c hjkl.Canvas) {
	for v, g := range g.Level {
		c.Blit(v, g.Face)
	}
	c.Blit(g.Hero.Pos.Offset, g.Hero.Face)
}

func main() {
	if err := hjkl.Run(NewGame()); err != nil {
		panic(err)
	}
}
