package tui

import "testing"

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
	mammoth := &Obj{"mammoth"}
	tiger := &Obj{"tiger"}

	cases := []struct {
		Fmt      string
		Args     []any
		Expected string
	}{
		{"%s %v %o", []any{I, "hit", mammoth}, "I hit the mammoth."},
		{"%s %v %o", []any{you, "hit", mammoth}, "You hit the mammoth."},
		{"%s %v %o", []any{Grog, "hit", mammoth}, "Grog hits the mammoth."},
		{"%s %v %o", []any{tiger, "bite", mammoth}, "The tiger bites the mammoth."},
		{"%s <hit> %o", []any{I, mammoth}, "I hit the mammoth."},
		{"%s <hit> %o", []any{you, mammoth}, "You hit the mammoth."},
		{"%s <hit> %o", []any{Grog, mammoth}, "Grog hits the mammoth."},
		{"%s <bite> %o", []any{tiger, mammoth}, "The tiger bites the mammoth."},
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
		{"%s %v %o", []any{tiger, "lick", tiger}, "The tiger licks itself."},
		{"%s %v %o!", []any{I, "hit", mammoth}, "I hit the mammoth!"},
		{"%s %v %o.", []any{I, "hit", mammoth}, "I hit the mammoth."},
		{"%s %v %o", []any{I, "hit"}, "I hit %!o(MISSING)."},
		{"%s %v", []any{you, "rest", tiger}, "You rest.%!(EXTRA tiger)"},
		{"%s %v %o", []any{"", "hit", mammoth}, "%!s(EMPTY) hits the mammoth."},
		{"%s %v %o", []any{I, "", mammoth}, "I %!v(EMPTY) the mammoth."},
		{"%s %v %o", []any{you, "hit", ""}, "You hit %!o(EMPTY)."},
	}
	for _, c := range cases {
		if actual := Log(c.Fmt, c.Args...).Message; actual != c.Expected {
			t.Errorf("Got: \"%s\", Wanted: \"%s\"", actual, c.Expected)
		}
	}
}
