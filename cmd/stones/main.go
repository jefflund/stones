package main

import "github.com/jefflund/stones/pkg/hjkl"

const (
	COLS = 80
	ROWS = 24
)

var (
	WALL  = hjkl.Ch('#')
	FLOOR = hjkl.Ch('.')
)

type Game struct {
	Hero  hjkl.Vector
	Level map[hjkl.Vector]hjkl.Glyph
}

func NewGame() *Game {
	hero := hjkl.Vec(40, 12)
	level := make(map[hjkl.Vector]hjkl.Glyph)
	for x := 0; x < COLS; x++ {
		for y := 0; y < ROWS; y++ {
			tile := FLOOR
			if x == 0 || x == COLS-1 || y == 0 || y == ROWS-1 {
				tile = WALL
			}
			level[hjkl.Vec(x, y)] = tile
		}
	}
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
	if dst := g.Hero.Add(delta); g.Level[dst] == FLOOR {
		g.Hero = dst
	}
	return nil
}

func (g *Game) Draw(c hjkl.Canvas) {
	for v, g := range g.Level {
		c.Blit(v, g)
	}
	c.Blit(g.Hero, hjkl.Ch('@'))
}

func main() {
	if err := hjkl.Run(NewGame()); err != nil {
		panic(err)
	}
}
