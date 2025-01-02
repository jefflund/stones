// Package rpg provides the rpg mechanics for stones.
package rpg

import "github.com/jefflund/stones/pkg/hjkl"

type Damage struct {
	Amount int
}

type Attributes struct {
	MaxHealth int
	Damage    int
}

type Variables struct {
	Health int
}

type Character struct {
	Attributes
	Variables
}

func (c *Character) Handle(e *hjkl.Mob, v hjkl.Event) {
	switch v := v.(type) {
	case *hjkl.Bump:
		v.Bumped.Handle(&Damage{c.Damage})
	case *Damage:
		c.Health -= v.Amount
		if c.Health <= 0 {
			e.Pos.Occupant = nil
			e.Pos = nil
		}
	}
}
