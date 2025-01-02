package rpg

import "github.com/jefflund/stones/pkg/hjkl"

type BestiaryEntry struct {
	Face       hjkl.Glyph
	Attributes Attributes
}

func (b BestiaryEntry) New() *hjkl.Mob {
	m := hjkl.NewMob(b.Face)
	m.Components.Add(&Character{
		Attributes: b.Attributes,
		Variables: Variables{
			Health: b.Attributes.MaxHealth,
		},
	})
	return m
}

var Bestiary = []BestiaryEntry{
	{
		Face: hjkl.ChFg('U', hjkl.ColorRed),
		Attributes: Attributes{
			MaxHealth: 10,
			Damage:    3,
		},
	},
	{
		Face: hjkl.ChFg('u', hjkl.ColorRed),
		Attributes: Attributes{
			MaxHealth: 3,
			Damage:    1,
		},
	},
	{
		Face: hjkl.ChFg('u', hjkl.ColorLightRed),
		Attributes: Attributes{
			MaxHealth: 3,
			Damage:    2,
		},
	},
	{
		Face: hjkl.ChFg('a', hjkl.ColorLightBlue),
		Attributes: Attributes{
			MaxHealth: 5,
			Damage:    1,
		},
	},
	{
		Face: hjkl.ChFg('A', hjkl.ColorLightBlue),
		Attributes: Attributes{
			MaxHealth: 20,
			Damage:    1,
		},
	},
}

func NewHero() *hjkl.Mob {
	entry := BestiaryEntry{
		Face: hjkl.Ch('@'),
		Attributes: Attributes{
			MaxHealth: 10,
			Damage:    2,
		},
	}
	return entry.New()
}
