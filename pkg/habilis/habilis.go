// Package habilis implements the habilis system from Cavemaster RPG.
package habilis

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

// Skin represents a character as a collection of Circle.
type Skin struct {
	Name    string
	Circles []Circle
}
