package hjkl

// Event is a message handled by Entity.
type Event any

// Entity is a game object which handles Event.
type Entity interface {
	Handle(Event)
}

// Component handles Event for Entity.
type Component[E Entity] interface {
	Handle(E, Event)
}

// ComponentFunc is a function which acts as a Component.
type ComponentFunc[E Entity] func(E, Event)

// Handle calls the underlying function.
func (c ComponentFunc[E]) Handle(e E, v Event) {
	c(e, v)
}

// ComponentSlice is a Component composed of other Component.
type ComponentSlice[E Entity] []Component[E]

// Add appends a Component to the ComponentSlice.
func (s *ComponentSlice[E]) Add(c Component[E]) {
	*s = append(*s, c)
}

// Handle has each constiutent Component handle the Event.
func (s ComponentSlice[E]) Handle(e E, v Event) {
	for _, c := range s {
		c.Handle(e, v)
	}
}

// Handler creates a ComponentFunc which handles a specific Event type.
func Handler[E Entity, V Event](f func(E, V)) ComponentFunc[E] {
	return ComponentFunc[E](func(e E, v Event) {
		if v, ok := v.(V); ok {
			f(e, v)
		}
	})
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

	Components ComponentSlice[*Mob]
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

	e.Components.Handle(e, v)
}

// Tile represents a single square in the game mpa.
type Tile struct {
	Offset   Vector
	Face     Glyph
	Pass     bool
	Occupant *Mob
	Adjacent map[Vector]*Tile

	Components ComponentSlice[*Tile]
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

	e.Components.Handle(e, v)
}

// PlaceMob places a Mob on a Tile.
func PlaceMob(m *Mob, t *Tile) {
	t.Occupant = m
	m.Pos = t
}
