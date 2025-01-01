// Package gen provides proceedurally generated Tile maps.
package gen

import "github.com/jefflund/stones/pkg/hjkl"

// GenTileGrid creates a two-dimensional grid of Tile.
func GenTileGrid(cols, rows int, f func(hjkl.Vector) *hjkl.Tile) []*hjkl.Tile {
	grid := make(map[hjkl.Vector]*hjkl.Tile)
	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			v := hjkl.Vec(x, y)
			grid[v] = f(v)
		}
	}

	for off, src := range grid {
		for _, delta := range hjkl.CompassDirs {
			if dst, ok := grid[off.Add(delta)]; ok {
				src.Adjacent[delta] = dst
			}
		}
	}

	tiles := make([]*hjkl.Tile, 0, len(grid))
	for _, t := range grid {
		tiles = append(tiles, t)
	}
	return tiles
}

// GenFence applies a function to any Tile which is not 8-connected.
func GenFence(tiles []*hjkl.Tile, f func(*hjkl.Tile)) {
	for _, t := range tiles {
		if len(t.Adjacent) != 8 {
			f(t)
		}
	}
}
