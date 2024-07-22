package hjkl

// GenTile is function which generates a new Tile.
type GenTile func(Vector) *Tile

// ModTile is a function which modifies an existing Tile.
type ModTile func(*Tile)

// GenTileGrid creates a new eight-connected grid of Tile.
func GenTileGrid(cols, rows int, f GenTile) []*Tile {
	grid := make(map[Vector]*Tile)
	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			grid[Vector{x, y}] = f(Vector{x, y})
		}
	}

	for off, src := range grid {
		for _, delta := range dirs8 {
			if dst, ok := grid[off.Add(delta)]; ok {
				src.Adjacent[delta] = dst
			}
		}
	}

	tiles := make([]*Tile, 0, len(grid))
	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			tiles = append(tiles, grid[Vector{x, y}])
		}
	}
	return tiles
}

// GenFence modifies edge Tile.
func GenFence(tiles []*Tile, f ModTile) {
	for _, t := range tiles {
		if len(t.Adjacent) < 8 {
			f(t)
		}
	}
}
