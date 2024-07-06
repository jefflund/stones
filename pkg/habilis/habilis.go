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

	StoneAny Stone = 0

	StoneMelee = StoneHit | StoneEvs | StoneDmg | StoneArm
)

// Circle is a collection of Stone on a Skin.
type Circle struct {
	Name     string
	Stone    Stone
	Count    int
	MaxCount int
}

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

// OnCollide is currently a noop.
func (s *Skin) OnCollide(m *hjkl.Mob[Skin], t *hjkl.Tile[Skin]) {
}

// OnBump is currently a noop.
func (s *Skin) OnBump(m, b *hjkl.Mob[Skin]) {
}
