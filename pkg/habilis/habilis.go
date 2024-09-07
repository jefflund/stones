// Package habilis implements the habilis system from Cavemaster RPG.
package habilis

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

// NewCircle creates a Circle.
func NewCircle(name string, stone Stone, count int) *Circle {
	return &Circle{name, stone, count, count}
}

// Skin represents a character as a collection of Circle.
type Skin struct {
	Name    string
	Circles []*Circle
}

// NewSkin creates a Skin.
func NewSkin(name string, circles ...*Circle) *Skin {
	return &Skin{name, circles}
}
