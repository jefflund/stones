package hjkl

// Mob represents a creature capable of occupying Tile.
type Mob struct {
	Face Glyph
	Pos  *Tile
}

// NewMob constructs a new Mob with the given Glyph face.
func NewMob(face Glyph) *Mob {
	return &Mob{Face: face}
}

// Tile represents a single square in the game mpa.
type Tile struct {
	Offset   Vector
	Face     Glyph
	Pass     bool
	Occupant *Mob
	Adjacent map[Vector]*Tile
}

// NewTile constructs a new Tile with the given Vector offset.
func NewTile(offset Vector) *Tile {
	return &Tile{
		Offset:   offset,
		Face:     Ch('.'),
		Pass:     true,
		Adjacent: make(map[Vector]*Tile),
	}
}

// GenTileGrid creates a new eight-connected grid of Tile.
func GenTileGrid(cols, rows int, f func(Vector) *Tile) []*Tile {
	grid := make(map[Vector]*Tile)
	for x := 0; x < cols; x++ {
		for y := 0; y < rows; y++ {
			grid[Vec(x, y)] = f(Vec(x, y))
		}
	}

	for off, src := range grid {
		for _, delta := range CompassDirs {
			if dst, ok := grid[off.Add(delta)]; ok {
				src.Adjacent[delta] = dst
			}
		}
	}

	tiles := make([]*Tile, 0, len(grid))
	for _, t := range grid {
		tiles = append(tiles, t)
	}
	return tiles
}

// GenFence modifies edge Tile in a slice of Tile.
func GenFence(tiles []*Tile, f func(*Tile)) {
	for _, t := range tiles {
		if len(t.Adjacent) < 8 {
			f(t)
		}
	}
}
