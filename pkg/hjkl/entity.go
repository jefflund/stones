package hjkl

// Event is a message handled by Event.
type Event any

// Entity is a game object which handles Event.
type Entity interface {
	Handle(Event)
}

// Component processes Event for an Entity.
type Component interface {
	Process(Entity, Event)
}

// ComponentSlice is an Entity composed of Component.
type ComponentSlice []Component

// ComponentFunc is a function which implements Component.
type ComponentFunc func(Entity, Event)

// Process calls the underlying function.
func (c ComponentFunc) Process(e Entity, v Event) {
	c(e, v)
}

// Handle has each constituent Component process the Event.
func (e ComponentSlice) Handle(v Event) {
	for _, c := range e {
		c.Process(e, v)
	}
}

// MobComponent is a Component which lets an Entity work as a mob.
type MobComponent struct {
	Face Glyph
	Pos  Entity
}

// NewMob constructs an Entity with a MobComponent.
func NewMob(face Glyph, cs ...Component) ComponentSlice {
	c := &MobComponent{
		Face: face,
	}
	return append(ComponentSlice{c}, cs...)
}

// Process implements Component for MobComponent.
func (c *MobComponent) Process(e Entity, v Event) {
	switch v := v.(type) {
	case *MoveTrigger:
		dst := Adjacent(c.Pos, v.Delta)
		if !Pass(dst) {
			Collide(e, dst)
		} else if bumped := Occupant(dst); bumped != nil {
			Bump(e, bumped)
		} else {
			SetOccupant(c.Pos, nil)
			SetOccupant(dst, e)
			SetPos(e, dst)
		}
	case *FaceQuery:
		v.Response = c.Face
	case *FaceUpdate:
		c.Face = v.Update
	case *PosQuery:
		v.Response = c.Pos
	case *PosUpdate:
		c.Pos = v.Update
	case *OffsetQuery:
		if c.Pos != nil {
			c.Pos.Handle(v)
		}
	}
}

// TileComponent is a Component which lets an Entity work as a tile.
type TileComponent struct {
	Face     Glyph
	Occupant Entity
	Pass     bool
	Offset   Vector
	Adjacent map[Vector]Entity
}

// NewTile constructs an Entity with a MobComponent.
func NewTile(offset Vector, cs ...Component) ComponentSlice {
	c := &TileComponent{
		Face:     Ch(' '),
		Pass:     true,
		Offset:   offset,
		Adjacent: make(map[Vector]Entity),
	}
	return append(ComponentSlice{c}, cs...)
}

// Process implements Component for TileComponent.
func (c *TileComponent) Process(e Entity, v Event) {
	switch v := v.(type) {
	case *FaceQuery:
		v.Response = c.Face
		if c.Occupant != nil {
			c.Occupant.Handle(v)
		}
	case *FaceUpdate:
		c.Face = v.Update
	case *OccupantQuery:
		v.Response = c.Occupant
	case *OccupantUpdate:
		c.Occupant = v.Update
	case *PassQuery:
		v.Response = c.Pass
	case *PassUpdate:
		c.Pass = v.Update
	case *OffsetQuery:
		v.Response = c.Offset
	case *AdjacentQuery:
		v.Response = c.Adjacent[v.Delta]
	case *AdjacentUpdate:
		c.Adjacent[v.Delta] = v.Update
	}
}

// PlaceMob sends a PosUpdate to a mob and an OccupantUpdate to a tile.
func PlaceMob(mob, tile Entity) {
	SetPos(mob, tile)
	SetOccupant(tile, mob)
}

// MoveTrigger is an Event triggering movement in an Entity.
type MoveTrigger struct {
	Delta Vector
}

// Move sends a MoveTrigger to an Entity.
func Move(e Entity, delta Vector) {
	e.Handle(&MoveTrigger{delta})
}

// CollisionEvent is an Event indicating collision with an Entity obstacle.
type CollisionEvent struct {
	Obstacle Entity
}

// Collide sends a CollisionEvent to an Entity.
func Collide(e, obstacle Entity) {
	e.Handle(&CollisionEvent{obstacle})
}

// BumpEvent is an Event indicating that an Entity bumped another.
type BumpEvent struct {
	Bumped Entity
}

// Bump sends a BumpEvent to an Entity.
func Bump(e, bumped Entity) {
	e.Handle(&BumpEvent{bumped})
}

// FaceQuery is an Event querying an Entity for its Glyph face.
type FaceQuery struct {
	Response Glyph
}

// Face sends a FaceQuery to an Entity, and returns the response.
func Face(e Entity) Glyph {
	v := FaceQuery{}
	e.Handle(&v)
	return v.Response
}

// FaceUpdate is an Event updating the Glyph face of an Entity.
type FaceUpdate struct {
	Update Glyph
}

// SetFace sends a FaceUpdate to an Entity.
func SetFace(e Entity, update Glyph) {
	e.Handle(&FaceUpdate{update})
}

// PosQuery is an Event querying an Entity for its Entity pos.
type PosQuery struct {
	Response Entity
}

// Pos sends a PosQuery to an Entity and returns the response.
func Pos(e Entity) Entity {
	v := PosQuery{}
	e.Handle(&v)
	return v.Response
}

// PosUpdate is an Event updating the Entity pos of an Entity.
type PosUpdate struct {
	Update Entity
}

// SetPos sends a PosUpdate to an Entity.
func SetPos(e, update Entity) {
	e.Handle(&PosUpdate{update})
}

// OccupantQuery is an Event querying an Entity for its Entity occupant.
type OccupantQuery struct {
	Response Entity
}

// Occupant sends an OccupantQuery to an Entity and returns the response.
func Occupant(e Entity) Entity {
	v := OccupantQuery{}
	e.Handle(&v)
	return v.Response
}

// OccupantUpdate is an Event updating the Entity occupant of an Entity.
type OccupantUpdate struct {
	Update Entity
}

// SetOccupant sends an OccupantUpdate to an Entity.
func SetOccupant(e, update Entity) {
	e.Handle(&OccupantUpdate{update})
}

// PassQuery is an Event querying an Entity for its bool passability.
type PassQuery struct {
	Response bool
}

// Pass sends a PassQuery to an Entity and returns the result.
func Pass(e Entity) bool {
	v := PassQuery{}
	e.Handle(&v)
	return v.Response
}

// PassUpdate is an Event updating the bool passability of an Entity.
type PassUpdate struct {
	Update bool
}

// SetPass sends a PassUpdate to an Entity.
func SetPass(e Entity, update bool) {
	e.Handle(&PassUpdate{update})
}

// OffsetQuery is an Event querying an Entity for its Vector offset.
type OffsetQuery struct {
	Response Vector
}

// Offset sends an OffsetQuery to an Entity and returns the response.
func Offset(e Entity) Vector {
	v := OffsetQuery{}
	e.Handle(&v)
	return v.Response
}

// AdjacentQuery is an Event querying an Entity for an adjacent Entity.
type AdjacentQuery struct {
	Delta    Vector
	Response Entity
}

// Adjacent sends an AdjacentQuery to an Event and returns the response.
func Adjacent(e Entity, delta Vector) Entity {
	v := AdjacentQuery{Delta: delta}
	e.Handle(&v)
	return v.Response
}

// AdjacentUpdate is an Event updating an adjacent Entity of an Entity.
type AdjacentUpdate struct {
	Delta  Vector
	Update Entity
}

// SetAdjacent sends an AdjacentUpdate to an Entity.
func SetAdjacent(e Entity, delta Vector, update Entity) {
	e.Handle(&AdjacentUpdate{delta, update})
}
