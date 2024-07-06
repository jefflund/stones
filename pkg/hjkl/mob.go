package hjkl

// Mob is a game object which occupies Tile.
type Mob[T any] struct {
	Face Glyph
	Pos  *Tile[T]

	Data T

	OnCollide func(*Mob[T], *Tile[T])
	OnBump    func(*Mob[T], *Mob[T])
}

// NewMob constructs a Mob with the given Glyph face.
func NewMob[T any](face Glyph, data T) *Mob[T] {
	return &Mob[T]{
		Face: face,
		Data: data,
	}
}

// Move attempts to move the Mob to a new Tile.
func (m *Mob[T]) Move(delta Vector) {
	dst := m.Pos.Adjacent[delta]
	if !dst.Pass {
		m.OnCollide(m, dst)
	} else if dst.Occupant != nil {
		m.OnBump(m, dst.Occupant)
	} else {
		m.Pos.Occupant = nil
		dst.Occupant = m
		m.Pos = dst
	}
}

// Tile is a square in the game map.
type Tile[T any] struct {
	Face     Glyph
	Pass     bool
	Occupant *Mob[T]

	Offset   Vector
	Adjacent map[Vector]*Tile[T]
}

// NewTile creates a new Tile with the given Vector offset.
func NewTile[T any](offset Vector) *Tile[T] {
	return &Tile[T]{
		Face:     Ch('.'),
		Pass:     true,
		Offset:   offset,
		Adjacent: make(map[Vector]*Tile[T]),
	}
}
