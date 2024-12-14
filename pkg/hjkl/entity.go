package hjkl

// Event is a message handled by Entity.
type Event any

// Entity is a game object which handles Event.
type Entity interface {
	Handle(Event)
}

// Getter is an Event which gets a value.
type Getter[T any] interface {
	Get() T
}

// Field allows an Event to function as a Getter.
type Field[T any] struct {
	Value T
}

// Get gets the Field value.
func (v Field[T]) Get() T {
	return v.Value
}

// Get sends a Getter Event to an Entity, then returns the value.
func Get[T any](e Entity, v Getter[T]) T {
	e.Handle(v)
	return v.Get()
}

// Face is an Event which gets a Glyph face.
type Face struct {
	Field[Glyph]
}

// SetPos is an Event updating a Tile position.
type SetPos struct {
	Value *Tile
}

// SetOccupant is an Event updating a Mob occupant.
type SetOccupant struct {
	Value *Mob
}

// Move is an Event which triggers movement.
type Move struct {
	Delta Vector
}

// Bump is an Event sent upon bumping a Mob.
type Bump struct {
	Bumped *Mob
}

// Collide is an Event sent upon colliding with a Tile.
type Collide struct {
	Obstacle *Tile
}

// Mob represents a creature capable of occupying Tile.
type Mob struct {
	Face Glyph
	Pos  *Tile
}

// NewMob constructs a new Mob with the given Glyph face.
func NewMob(face Glyph) *Mob {
	return &Mob{Face: face}
}

// Handle implements Entity for Mob.
func (e *Mob) Handle(v Event) {
	switch v := v.(type) {
	case *Face:
		v.Value = e.Face
	case *SetPos:
		e.Pos = v.Value
	case *Move:
		if dst, ok := e.Pos.Adjacent[v.Delta]; ok {
			if dst.Occupant != nil {
				e.Handle(&Bump{dst.Occupant})
			} else if !dst.Pass {
				e.Handle(&Collide{dst})
			} else {
				e.Pos.Handle(&SetOccupant{nil})
				dst.Handle(&SetOccupant{e})
				e.Handle(&SetPos{dst})
			}
		}
	}
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

// Handle implements Entity for Tile.
func (e *Tile) Handle(v Event) {
	switch v := v.(type) {
	case *Face:
		v.Value = e.Face
		if e.Occupant != nil {
			e.Occupant.Handle(v)
		}
	case *SetOccupant:
		e.Occupant = v.Value
	}
}

// PlaceMob places a Mob on a Tile.
func PlaceMob(m *Mob, t *Tile) {
	t.Occupant = m
	m.Pos = t
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
