// Package gen provides procedural generation for a roguelike game.
package gen

import (
	"github.com/jefflund/stones/pkg/hjkl"
	"github.com/jefflund/stones/pkg/hjkl/rl"
)

// GenTile is function which generates a new Tile.
type GenTile func(hjkl.Vector) *rl.Tile

// ModTile is a function which modifies an existing Tile.
type ModTile func(*rl.Tile)

// GenTileGrid creates a new eight-connected grid of Tile.
func GenTileGrid(cols, rows int, f GenTile) []*rl.Tile {
	grid := make(map[hjkl.Vector]*rl.Tile)
	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			grid[hjkl.Vec(x, y)] = f(hjkl.Vec(x, y))
		}
	}

	for off, src := range grid {
		for _, delta := range hjkl.Dirs8 {
			if dst, ok := grid[off.Add(delta)]; ok {
				src.Adjacent[delta] = dst
			}
		}
	}

	tiles := make([]*rl.Tile, 0, len(grid))
	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			tiles = append(tiles, grid[hjkl.Vec(x, y)])
		}
	}
	return tiles
}

// GenFence modifies edge Tile.
func GenFence(tiles []*rl.Tile, f ModTile) {
	for _, t := range tiles {
		if len(t.Adjacent) < 8 {
			f(t)
		}
	}
}
