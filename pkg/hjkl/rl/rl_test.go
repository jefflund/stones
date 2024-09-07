package rl

import (
	"testing"

	"github.com/jefflund/stones/pkg/hjkl"
)

func TestMob_Move(t *testing.T) {
	mob := NewMob(hjkl.Ch('@'))
	src := NewTile(hjkl.Vec(0, 0))
	dst := NewTile(hjkl.Vec(1, 1))

	mob.Pos = src
	src.Adjacent[hjkl.Vec(1, 1)] = dst
	src.Occupant = mob

	mob.Move(hjkl.Vec(1, 1))

	if src.Occupant != nil {
		t.Error("Move failed to update src.Occupant")
	}
	if dst.Occupant != mob {
		t.Error("Move failed to update dst.Occupant")
	}
	if mob.Pos != dst {
		t.Error("Move failed to update mob.Pos")
	}
}

func TestMob_CollideTriggerMove(t *testing.T) {
	mob := NewMob(hjkl.Ch('@'))
	src := NewTile(hjkl.Vec(0, 0))
	dst := NewTile(hjkl.Vec(1, 1))

	collideSent := false
	mob.Pos = src
	mob.AddComponent(ComponentFunc[CollideEvent](func(m *Mob, v *CollideEvent) {
		if m == mob && v.Obstacle == dst {
			collideSent = true
		}
	}))
	src.Adjacent[hjkl.Vec(1, 1)] = dst
	src.Occupant = mob
	dst.Pass = false

	mob.Move(hjkl.Vec(1, 1))

	if !collideSent {
		t.Error("Move failed to send CollideEvent")
	}
	if src.Occupant != mob {
		t.Error("Move erroneously updated src.Occupant on collision")
	}
	if dst.Occupant != nil {
		t.Error("Move erroneously updated dst.Occupant on collision")
	}
	if mob.Pos != src {
		t.Error("Move erroneously updated mob.Pos on collision")
	}
}

func TestMob_Bump(t *testing.T) {
	mob := NewMob(hjkl.Ch('@'))
	bumped := NewMob(hjkl.Ch('D'))
	src := NewTile(hjkl.Vec(0, 0))
	dst := NewTile(hjkl.Vec(1, 1))

	bumpSent := false
	mob.Pos = src
	mob.AddComponent(ComponentFunc[BumpEvent](func(m *Mob, v *BumpEvent) {
		if m == mob && v.Bumped == bumped {
			bumpSent = true
		}
	}))
	src.Adjacent[hjkl.Vec(1, 1)] = dst
	src.Occupant = mob
	dst.Occupant = bumped

	mob.Move(hjkl.Vec(1, 1))

	if !bumpSent {
		t.Error("Move failed to send BumpEvent")
	}
	if src.Occupant != mob {
		t.Error("Move erroneously updated src.Occupant on bump")
	}
	if dst.Occupant != bumped {
		t.Error("Move erroneously updated dst.Occupant on bump")
	}
	if mob.Pos != src {
		t.Error("Move erroneously updated mob.Pos on bump")
	}
}

func TestPlaceMob(t *testing.T) {
	mob := NewMob(hjkl.Ch('@'))
	dst := NewTile(hjkl.Vec(1, 1))
	PlaceMob(mob, dst)
	if mob.Pos != dst {
		t.Error("PlaceMob failed to set mob.Pos")
	}
	if dst.Occupant != mob {
		t.Error("PlaceMob failed to set dst.Occupant")
	}
}
