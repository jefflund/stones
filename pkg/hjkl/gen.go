package hjkl

type GenTile func(Vector) Entity

type ModTile func(Entity)

func GenTileGrid(cols, rows int, f GenTile) []Entity {
	grid := make(map[Vector]Entity)
	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			grid[Vector{x, y}] = f(Vector{x, y})
		}
	}

	for off, src := range grid {
		for _, delta := range dirs8 {
			if dst, ok := grid[off.Add(delta)]; ok {
				SetAdjacent(src, delta, dst)
			}
		}
	}

	tiles := make([]Entity, 0, len(grid))
	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			tiles = append(tiles, grid[Vector{x, y}])
		}
	}
	return tiles
}

func GneFence(tiles []Entity, f ModTile) {
	for _, t := range tiles {
		if Outdegree(t) < 8 {
			f(t)
		}
	}
}
