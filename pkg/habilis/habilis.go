// Package habilis implements the habilis system from Cavemaster RPG.
package habilis

import (
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
	case *rl.BumpEvent:
		m.Handle(tui.Log("%s <bump> %o", m, v.Bumped))
	case *rl.CollideEvent:
		m.Handle(tui.Log("%s <collide> with %o", m, string(v.Obstacle.Face.Ch)))
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

// Count gets the maximum count of a Stone type on Skin.
func (s *Skin) MaxCount(t Stone) int {
	count := 0
	for _, c := range s.Circles {
		if c.Matches(t) {
			count += c.MaxCount
		}
	}
	return count
}
