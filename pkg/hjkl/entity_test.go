package hjkl

import (
	"testing"
)

func TestMob_Move(t *testing.T) {
	mob := NewMob(Ch('@'))
	src := NewTile(Vec(0, 0))
	dst := NewTile(Vec(1, 1))

	mob.Pos = src
	src.Adjacent[Vec(1, 1)] = dst
	src.Occupant = mob

	mob.Move(Vec(1, 1))

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
	mob := NewMob(Ch('@'))
	src := NewTile(Vec(0, 0))
	dst := NewTile(Vec(1, 1))

	mob.Pos = src
	src.Adjacent[Vec(1, 1)] = dst
	src.Occupant = mob
	dst.Pass = false

	mob.Move(Vec(1, 1))

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
	mob := NewMob(Ch('@'))
	bumped := NewMob(Ch('D'))
	src := NewTile(Vec(0, 0))
	dst := NewTile(Vec(1, 1))

	mob.Pos = src
	src.Adjacent[Vec(1, 1)] = dst
	src.Occupant = mob
	dst.Occupant = bumped

	mob.Move(Vec(1, 1))

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
