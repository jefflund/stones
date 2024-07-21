package habilis

import "github.com/jefflund/stones/pkg/hjkl"

// BestiaryEntry stores the data needed to create a Mob[Skin].
type BestiaryEntry struct {
	Face    hjkl.Glyph
	Circles []Circle
}

// Bestiary contains BestiaryEntry with data for various named Mob[Skin].
var Bestiary = map[string]BestiaryEntry{
	"Mammoth": {
		hjkl.Ch('M'), []Circle{
			NewCircle("Core", StoneCore, 9),
			NewCircle("Run", StoneMisc, 1),
			NewCircle("Intelligent", StoneMisc, 1),
			NewCircle("Tusk", StoneMelee, 1),
			NewCircle("Tough", StoneArm, 1),
		},
	},
	"Sabertooth": {
		hjkl.Glyph{Ch: 't', Fg: hjkl.ColorYellow}, []Circle{
			NewCircle("Core", StoneCore, 6),
			NewCircle("Run", StoneMisc, 1),
			NewCircle("Bite", StoneMelee, 2),
		},
	},
}

// NewBestiaryMob creates a Mob[Skin] from a named BestiaryEntry.
func NewBestiaryMob(name string) *hjkl.Mob[Skin] {
	e := Bestiary[name]
	return NewSkinMob(name, e.Face, e.Circles...)
}
