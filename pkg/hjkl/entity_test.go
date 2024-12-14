package hjkl

import "testing"

func TestNewMob(t *testing.T) {
	m := NewMob(Ch('@'))
	if m.Face != Ch('@') {
		t.Error("NewMob produced incorrect Face")
	}
}

func TestMob_Face(t *testing.T) {
	m := NewMob(Ch('@'))
	if got := Get(m, &Face{}); got != Ch('@') {
		t.Errorf("Get(Mob, Face) got %v", got)
	}
}

func TestMob_SetPos(t *testing.T) {
	m := NewMob(Ch('@'))
	a := NewTile(Vector{})
	m.Handle(&SetPos{Value: a})
	if m.Pos != a {
		t.Error("Mob.Handle(SetPos) failed to set Pos")
	}
}

func TestMob_Bump(t *testing.T) {
	m := NewMob(Ch('@'))
	n := NewMob(Ch('D'))
	a := NewTile(Vec(1, 1))
	b := NewTile(Vec(2, 0))
	a.Adjacent[Vec(2, 1)] = b
	PlaceMob(m, a)
	PlaceMob(n, b)
	m.Handle(&Move{Vec(2, 1)})
	if a.Occupant != m {
		t.Error("Mob.Handle(Move) incorrectly set old Occupant")
	}
	if b.Occupant != n {
		t.Error("Mob.Handle(Move) incorrectly set new Occupant")
	}
	if m.Pos != a {
		t.Error("Mob.Handle(Move) incorrectly set Pos")
	}
}

func TestMob_Collide(t *testing.T) {
	m := NewMob(Ch('@'))
	a := NewTile(Vec(1, 1))
	b := NewTile(Vec(2, 0))
	a.Adjacent[Vec(2, 1)] = b
	PlaceMob(m, a)
	b.Pass = false
	m.Handle(&Move{Vec(2, 1)})
	if a.Occupant != m {
		t.Error("Mob.Handle(Move) incorrectly set old Occupant")
	}
	if b.Occupant != nil {
		t.Error("Mob.Handle(Move) incorrectly set new Occupant")
	}
	if m.Pos != a {
		t.Error("Mob.Handle(Move) incorrectly set Pos")
	}
}

func TestMob_Move(t *testing.T) {
	m := NewMob(Ch('@'))
	a := NewTile(Vec(1, 1))
	b := NewTile(Vec(2, 0))
	a.Adjacent[Vec(2, 1)] = b
	PlaceMob(m, a)
	m.Handle(&Move{Vec(2, 1)})
	if a.Occupant != nil {
		t.Error("Mob.Handle(Move) failed to set old Occupant")
	}
	if b.Occupant != m {
		t.Error("Mob.Handle(Move) failed to set new Occupant")
	}
	if m.Pos != b {
		t.Error("Mob.Handle(Move) failed to update Pos")
	}
}

func TestNewTile(t *testing.T) {
	a := NewTile(Vec(4, 2))
	if a.Offset != Vec(4, 2) {
		t.Error("NewTile produced incorrect Offset")
	}
	if a.Face != Ch('.') {
		t.Error("NewTile produced incorrect Face")
	}
	if !a.Pass {
		t.Error("NewTile produced incorrect Pass")
	}
	if a.Adjacent == nil {
		t.Error("NewTile failed to initialize Adjacent")
	}
}

func TestTile_Face(t *testing.T) {
	a := NewTile(Vector{})
	if got := Get(a, &Face{}); got != Ch('.') {
		t.Errorf("Get(Mob, Face) got %v", got)
	}
}

func TestTile_OccupiedFace(t *testing.T) {
	m := NewMob(Ch('@'))
	a := NewTile(Vector{})
	PlaceMob(m, a)
	if got := Get(a, &Face{}); got != Ch('@') {
		t.Errorf("Get(Mob, Face) got %v while occupied", got)
	}
}

func TestPlaceMob(t *testing.T) {
	m := NewMob(Ch('@'))
	a := NewTile(Vector{})
	PlaceMob(m, a)
	if m.Pos != a {
		t.Error("PlaceMob failed to update Pos")
	}
	if a.Occupant != m {
		t.Error("PlaceMob failed to update Occupant")
	}
}
