package hjkl

import (
	"fmt"
	"testing"
)

type Obj struct {
	Name string
}

func (o *Obj) String() string {
	return o.Name
}

func TestLog(t *testing.T) {
	I := &Obj{"I"}
	you := &Obj{"you"}
	Grog := &Obj{"Grog"}
	orc := &Obj{"orc"}
	goblin := &Mob{}
	goblin.Components = append(goblin.Components, Handler(func(m *Mob, v *NameQuery) {
		v.Value = "goblin"
	}))

	cases := []struct {
		Fmt  string
		Args []any
		Want string
	}{
		{"%s %v %o", []any{I, "hit", orc}, "I hit the orc."},
		{"%s %v %o", []any{you, "hit", orc}, "You hit the orc."},
		{"%s %v %o", []any{Grog, "hit", orc}, "Grog hits the orc."},
		{"%s %v %o", []any{goblin, "bite", orc}, "The goblin bites the orc."},
		{"%s <hit> %o", []any{I, orc}, "I hit the orc."},
		{"%s <hit> %o", []any{you, orc}, "You hit the orc."},
		{"%s <hit> %o", []any{Grog, orc}, "Grog hits the orc."},
		{"%s <bite> %o", []any{goblin, orc}, "The goblin bites the orc."},
		{"%s %v %o", []any{"a tick", "bite", "you"}, "A tick bites you."},
		{"%s %v %o", []any{you, "eat", "an apple"}, "You eat an apple."},
		{"%s %v %o", []any{you, "eat", "the apple"}, "You eat the apple."},
		{"%s %v %o", []any{Grog, "deny", you}, "Grog denies you."},
		{"%s %v", []any{"snake", "hiss"}, "The snake hisses."},
		{"%s %v hungry", []any{I, "be"}, "I am hungry."},
		{"%s %v hungry", []any{you, "be"}, "You are hungry."},
		{"%s %v hungry", []any{Grog, "be"}, "Grog is hungry."},
		{"%s %v %o", []any{I, "heal", I}, "I heal myself."},
		{"%s %v %o", []any{you, "heal", you}, "You heal yourself."},
		{"%s %v %o", []any{Grog, "heal", Grog}, "Grog heals themself."},
		{"%s %v %o", []any{goblin, "lick", goblin}, "The goblin licks itself."},
		{"%s %v %o!", []any{I, "hit", orc}, "I hit the orc!"},
		{"%s %v %o.", []any{I, "hit", orc}, "I hit the orc."},
		{"%s %v %o", []any{I, "hit"}, "I hit %!o(MISSING)."},
		{"%s %v", []any{you, "rest", goblin}, "You rest.%!(EXTRA goblin)"},
		{"%s %v %o", []any{"", "hit", orc}, "%!s(EMPTY) hits the orc."},
		{"%s %v %o", []any{I, "", orc}, "I %!v(EMPTY) the orc."},
		{"%s %v %o", []any{you, "hit", ""}, "You hit %!o(EMPTY)."},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			if got := Log(c.Fmt, c.Args...); got != c.Want {
				t.Errorf("Got: \"%s\", Wanted: \"%s\"", got, c.Want)
			}
		})
	}
}
