package main

import "github.com/jefflund/stones/pkg/hjkl"

var KeyMap = map[hjkl.Key]hjkl.Vector{
	'h': hjkl.Vec(-1, 0),
	'j': hjkl.Vec(0, 1),
	'k': hjkl.Vec(0, -1),
	'l': hjkl.Vec(1, 0),
}

type Game struct {
	Hero  *hjkl.Mob
	Tiles []*hjkl.Tile
}

func NewGame() *Game {
	hero := hjkl.NewMob(hjkl.Ch('@'))
	tiles := hjkl.GenTileGrid(80, 24, hjkl.NewTile)
	hjkl.GenFence(tiles, func(t *hjkl.Tile) {
		t.Face = hjkl.Ch('#')
		t.Pass = false
	})
	hero.Pos = tiles[1000]
	hero.Pos.Occupant = hero
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
	for _, t := range g.Tiles {
		c.Blit(t.Offset, t.Face)
	}
	c.Blit(g.Hero.Pos.Offset, g.Hero.Face)
}

func main() {
	if err := hjkl.Run(NewGame()); err != nil {
		panic(err)
	}
}
