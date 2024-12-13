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
}

// NewTile constructs a new Tile with the given Vector offset.
func NewTile(offset Vector) *Tile {
	return &Tile{
		Offset: offset,
		Face:   Ch('.'),
		Pass:   true,
	}
}
