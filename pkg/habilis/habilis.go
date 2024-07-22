// Package habilis implements the habilis system from Cavemaster RPG.
package habilis

import "github.com/jefflund/stones/pkg/hjkl"

// Stone is the core mechanic for habilis, representing health and abilities.
type Stone uint

// Stone constants, which function as bit flags.
const (
	StoneCore Stone = 1 << iota
	StoneHit
	StoneEvs
	StoneDmg
	StoneArm
	StoneMisc

	StoneAny  Stone = 0
	StoneNone Stone = ^StoneAny

	StoneMelee = StoneHit | StoneEvs | StoneDmg | StoneArm
)

// Circle is a collection of Stone on a Skin.
type Circle struct {
	Name     string
	Stone    Stone
	Count    int
	MaxCount int
}

// NewCircle creates a Circle.
func NewCircle(name string, stone Stone, count int) Circle {
	return Circle{name, stone, count, count}
}

// Skin represents a character as a collection of Circle.
type Skin struct {
	Name    string
	Circles []Circle
}

type SkinQuery struct {
	Response *Skin
}

func GetSkin(m *hjkl.Mob) *Skin {
	q := SkinQuery{}
	m.Handle(&q)
	return q.Response
}

func (s *Skin) Process(m *hjkl.Mob, v hjkl.Event) {
	switch v := v.(type) {
	case *hjkl.CollideEvent:
		s.Hurt()
	case *hjkl.BumpEvent:
		bs := GetSkin(v.Bumped)
		core := s.Roll(StoneNone) - bs.Roll(StoneNone)
		hit := core + s.Count(StoneHit) - bs.Count(StoneEvs)
		if hit > 0 {
			dmg := core + s.Count(StoneDmg) - bs.Count(StoneArm)
			for x := hjkl.Min(0, dmg); x < dmg; x++ {
				bs.Hurt()
			}
		}
	case *SkinQuery:
		v.Response = s
	}
}

// Count gets the total count for a Stone type.
func (s *Skin) Count(stone Stone) int {
	total := 0
	for _, c := range s.Circles {
		if c.Stone&stone == stone {
			total += c.Count
		}
	}
	return total
}

// Roll gets a roll for a particular Stone type.
func (s *Skin) Roll(bonus Stone) int {
	core := s.Count(StoneCore)
	return hjkl.RandRange(0, core) + s.Count(bonus)
}

// Hurt removes a pseudo-random Stone from the Skin. It panics if there is no
// Stone to remove.
func (s *Skin) Hurt() {
	index := hjkl.RandIndex(s.Circles, func(c Circle) int { return c.Count })
	s.Circles[index].Count--
}

// Heal restores a psueod-random missing Stone to the Skin. It panics if there
// is no missing Stone to restore.
func (s *Skin) Heal() {
	index := hjkl.RandIndex(s.Circles, func(c Circle) int {
		return c.MaxCount - c.Count
	})
	s.Circles[index].Count++
}
