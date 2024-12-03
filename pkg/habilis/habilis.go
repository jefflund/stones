// Package habilis implements the habilis system from Cavemaster RPG.
package habilis

import (
	"github.com/jefflund/stones/pkg/hjkl/math/rand"
	"github.com/jefflund/stones/pkg/hjkl/rl"
	"github.com/jefflund/stones/pkg/hjkl/tui"
)

// Stone is the core mechanic of habilis, representing health and abilities.
type Stone uint64

// Stone constants which function as bit flags.
const (
	StoneAny  Stone = 0
	StoneNone Stone = ^StoneAny

	StoneCore Stone = 1 << iota
	StoneHit
	StoneEvs
	StoneDmg
	StoneArm
	StoneMisc

	StoneMelee = StoneHit | StoneEvs | StoneDmg | StoneArm
)

// Circle is a collection of Stone on a Skin.
type Circle struct {
	Name     string
	Stone    Stone
	Count    int
	MaxCount int
}

// Matches returns true if the Circle Stone matches the given Stone.
func (c *Circle) Matches(t Stone) bool {
	return c.Stone&t == t
}

// CountQuery is an Event querying a Mob for the current count of a Stone.
type CountQuery struct {
	Stone Stone
	Count int
}

// Response implements Query for CountQuery.
func (q *CountQuery) Response() int {
	return q.Count
}

// RollQuery is an Event querying a Mob for a core roll.
type RollQuery struct {
	Roll int
}

// Response implements Query for RollQuery.
func (q *RollQuery) Response() int {
	return q.Roll
}

// Skin represents a character as a collection of Circle.
type Skin struct {
	Name    string
	Circles []*Circle
}

// Process implements Component for Skin.
func (s *Skin) Process(m *rl.Mob, v rl.Event) {
	switch v := v.(type) {
	case *tui.NameQuery:
		v.Name = s.Name
	case *CountQuery:
		v.Count = s.Count(v.Stone)
	case *RollQuery:
		v.Roll = s.Roll()
	case *rl.BumpEvent:
		Melee(m, v.Bumped)
	case *rl.CollideEvent:
		m.Handle(tui.Log("%s <collide> with %o", m, string(v.Obstacle.Face.Ch)))
	}
}

func Melee(a, b *rl.Mob) {
	roll := rl.Send(a, &RollQuery{}) - rl.Send(b, &RollQuery{})
	tohit := rl.Send(a, &CountQuery{Stone: StoneHit})
	toevs := rl.Send(b, &CountQuery{Stone: StoneEvs})
	if hit := roll + tohit - toevs; hit > 0 {
		todmg := rl.Send(a, &CountQuery{Stone: StoneDmg})
		toarm := rl.Send(b, &CountQuery{Stone: StoneArm})
		if dmg := roll + todmg - toarm; dmg > 0 {
			a.Handle(tui.Log("%s <hit> %o for %x", a, b, dmg))
		} else {
			a.Handle(tui.Log("%s <graze> %o", a, b))
		}
	} else {
		a.Handle(tui.Log("%s <miss> %o", a, b))
	}
}

// Count gets the current count of a Stone type on Skin.
func (s *Skin) Count(t Stone) int {
	count := 0
	for _, c := range s.Circles {
		if c.Matches(t) {
			count += c.Count
		}
	}
	return count
}

// Roll makes a uniform random core roll for the Skin.
func (s *Skin) Roll() int {
	return rand.Range(0, s.Count(StoneCore))
}
