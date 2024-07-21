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

// NewSkinMob creates a Mob[Skin].
func NewSkinMob(name string, face hjkl.Glyph, circles ...Circle) *hjkl.Mob[Skin] {
	s := Skin{name, circles}
	m := hjkl.NewMob(face, s)
	m.OnCollide = s.OnCollide
	m.OnBump = s.OnBump
	return m
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
	s.Circles[index].Count--
}

// OnCollide is currently a noop.
func (s *Skin) OnCollide(m *hjkl.Mob[Skin], t *hjkl.Tile[Skin]) {
}

// OnBump is currently a noop.
func (s *Skin) OnBump(m, b *hjkl.Mob[Skin]) {
}
