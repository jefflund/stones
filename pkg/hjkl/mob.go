package hjkl

// Mob is a game object which occupies Tile.
type Mob struct {
	Face Glyph
	Pos  *Tile
}

// NewMob constructs a Mob with the given Glyph face.
func NewMob(face Glyph) *Mob {
	return &Mob{
		Face: face,
	}
}

// Move attempts to move the Mob to a new Tile.
func (m *Mob) Move(delta Vector) {
	dst := m.Pos.Adjacent[delta]
	if !dst.Pass {
		m.Collide(dst)
	} else if dst.Occupant != nil {
		m.Bump(dst.Occupant)
	} else {
		m.Pos.Occupant = nil
		dst.Occupant = m
		m.Pos = dst
	}
}

// Collide is called when the Mob tries to move to an impassable Tile.
func (m *Mob) Collide(obstacle *Tile) {
}

// Bump is called when the Mob tries to move to an occupied Tile.
func (m *Mob) Bump(bumped *Mob) {
}

// Tile is a square in the game map.
type Tile struct {
	Face     Glyph
	Pass     bool
	Occupant *Mob

	Offset   Vector
	Adjacent map[Vector]*Tile
}

// NewTile creates a new Tile with the given Vector offset.
func NewTile(offset Vector) *Tile {
	return &Tile{
		Face:     Ch('.'),
		Pass:     true,
		Offset:   offset,
		Adjacent: make(map[Vector]*Tile),
	}
}
