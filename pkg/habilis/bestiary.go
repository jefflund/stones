package habilis

import (
	"github.com/jefflund/stones/pkg/hjkl"
	"github.com/jefflund/stones/pkg/hjkl/rl"
)

// CircleEntry contains the data needed to construct a Circle.
type CircleEntry struct {
	Name  string
	Stone Stone
	Count int
}

// New creates a new Circle using the data from the CircleEntry.
func (c CircleEntry) New() *Circle {
	return &Circle{c.Name, c.Stone, c.Count, c.Count}
}

// BestiaryEntry contains the data needed to construct a Mob with a Skin.
type BestiaryEntry struct {
	Name    string
	Face    hjkl.Glyph
	Circles []CircleEntry
}

// New creates a new Mob using the BestiaryEntry data.
func (b BestiaryEntry) New() *rl.Mob {
	circles := make([]*Circle, len(b.Circles))
	for i, c := range b.Circles {
		circles[i] = c.New()
	}
	mob := rl.NewMob(b.Face)
	mob.AddComponent(&Skin{b.Name, circles})
	return mob
}

// Bestiary contains the data needed to create new creature Mob.
var Bestiary = []BestiaryEntry{
	{
		"mammoth", hjkl.ChFg('M', hjkl.ColorLightBlack), []CircleEntry{
			{"Core", StoneCore, 9},
			{"Tusks", StoneMelee, 1},
			{"Tough", StoneArm, 1},
		},
	},
	{
		"saber-tooth", hjkl.ChFg('t', hjkl.ColorYellow), []CircleEntry{
			{"Core", StoneCore, 6},
			{"Bite", StoneMelee, 2},
		},
	},
}

// NewHero creates a new Mob to represent the player.
func NewHero() *rl.Mob {
	entry := BestiaryEntry{
		"you",
		hjkl.Ch('@'),
		[]CircleEntry{
			{"Core", StoneCore, 4},
			{"Rogok", StoneDmg, 1},
			{"Warrior", StoneMelee, 1},
		},
	}
	return entry.New()
}
