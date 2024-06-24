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
		dst := AdjacentQuery{Delta: v.Delta}.Send(c.Pos)
		if !(PassQuery{}).Send(dst) {
			e.Handle(&CollisionEvent{dst})
		} else if bumped := (OccupantQuery{}).Send(dst); bumped != nil {
			e.Handle(&BumpEvent{bumped})
		} else {
			c.Pos.Handle(&OccupantUpdate{nil})
			dst.Handle(&OccupantUpdate{e})
			e.Handle(&PosUpdate{dst})
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

// MoveTrigger is an Event triggering movement in an Entity.
type MoveTrigger struct {
	Delta Vector
}

// CollisionEvent is an Event indicating collision with an Entity obstacle.
type CollisionEvent struct {
	Obstacle Entity
}

// BumpEvent is an Event indicating that an Entity bumped another.
type BumpEvent struct {
	Bumped Entity
}

// FaceQuery is an Event querying an Entity for its Glyph face.
type FaceQuery struct {
	Response Glyph
}

// Send sends the FaceQuery to an Entity and returns the Glyph response.
func (v FaceQuery) Send(e Entity) Glyph {
	e.Handle(&v)
	return v.Response
}

// FaceUpdate is an Event updating the Glyph face of an Entity.
type FaceUpdate struct {
	Update Glyph
}

// PosQuery is an Event querying an Entity for its Entity pos.
type PosQuery struct {
	Response Entity
}

// Send sends the PosQuery to an Entity and returns the Entity response.
func (v PosQuery) Send(e Entity) Entity {
	e.Handle(&v)
	return v.Response
}

// PosUpdate is an Event updating the Entity pos of an Entity.
type PosUpdate struct {
	Update Entity
}

// OccupantQuery is an Event querying an Entity for its Entity occupant.
type OccupantQuery struct {
	Response Entity
}

// Send sends the OccupantQuery to an Entity and returns the Entity response.
func (v OccupantQuery) Send(e Entity) Entity {
	e.Handle(&v)
	return v.Response
}

// OccupantUpdate is an Event updating the Entity occupant of an Entity.
type OccupantUpdate struct {
	Update Entity
}

// PassQuery is an Event querying an Entity for its bool passability.
type PassQuery struct {
	Response bool
}

// Send sends the PassQuery to an Entity and returns the bool response.
func (v PassQuery) Send(e Entity) bool {
	e.Handle(&v)
	return v.Response
}

// PassUpdate is an Event updating the bool passability of an Entity.
type PassUpdate struct {
	Update bool
}

// OffsetQuery is an Event querying an Entity for its Vector offset.
type OffsetQuery struct {
	Response Vector
}

// Send sends the OffsetQuery to an Entity and returns the Vector response.
func (v OffsetQuery) Send(e Entity) Vector {
	e.Handle(&v)
	return v.Response
}

// AdjacentQuery is an Event querying an Entity for an adjacent Entity.
type AdjacentQuery struct {
	Delta    Vector
	Response Entity
}

// Send sends the AdjacentQuery to an Entity and returns the Entity response.
func (v AdjacentQuery) Send(e Entity) Entity {
	e.Handle(&v)
	return v.Response
}

// AdjacentUpdate is an Event updating an adjacent Entity of an Entity.
type AdjacentUpdate struct {
	Delta  Vector
	Update Entity
}
