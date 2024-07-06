package hjkl

// GenTile is function which generates a new Tile.
type GenTile[T any] func(Vector) *Tile[T]

// ModTile is a function which modifies an existing Tile.
type ModTile[T any] func(*Tile[T])

// GenTileGrid creates a new eight-connected grid of Tile.
func GenTileGrid[T any](cols, rows int, f GenTile[T]) []*Tile[T] {
	grid := make(map[Vector]*Tile[T])
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

	tiles := make([]*Tile[T], 0, len(grid))
	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			tiles = append(tiles, grid[Vector{x, y}])
		}
	}
	return tiles
}

// GenFence modifies edge Tile.
func GenFence[T any](tiles []*Tile[T], f ModTile[T]) {
	for _, t := range tiles {
		if len(t.Adjacent) < 8 {
			f(t)
		}
	}
}
