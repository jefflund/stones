package hjkl

import (
	"testing"
)

func TestNewMob(t *testing.T) {
	m := NewMob(Ch('@'))
	if m.Face != Ch('@') {
		t.Error("NewMob produced incorrect Face")
	}
}

type TestGetter struct {
	Field[string]
}

func TestGet(t *testing.T) {
	m := NewMob(Glyph{})
	var received Event
	m.Components = append(m.Components, Handler(func(e *Mob, v *TestGetter) {
		received = v
		v.Value = "foobar"
	}))
	v := &TestGetter{}
	if got := Get(m, v); got != "foobar" {
		t.Errorf("Get(Mob, TestGetter) returned wrong value %s", got)
	}
	if received != v {
		t.Error("Get(mob, TestGetter) failed to send Event")
	}
	if v.Value != "foobar" {
		t.Error("Get(Mob, TestGetter) failed to update Event")
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

type LogEntry struct {
	Entity any
	Event  any
}

func TestMob_Bump(t *testing.T) {
	m := NewMob(Ch('@'))
	n := NewMob(Ch('D'))
	a := NewTile(Vec(1, 1))
	b := NewTile(Vec(2, 1))
	a.Adjacent[Vec(1, 0)] = b
	PlaceMob(m, a)
	PlaceMob(n, b)

	var got *Bump
	m.Components = append(m.Components, Handler(func(e *Mob, v *Bump) {
		got = v
	}))

	m.Handle(&Move{Vec(1, 0)})

	if got == nil || got.Bumped != n {
		t.Error("Mob.Handle(Move) sent incorrect Bump")
	}
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

	var got *Collide
	m.Components = append(m.Components, Handler(func(e *Mob, v *Collide) {
		got = v
	}))

	m.Handle(&Move{Vec(2, 1)})

	if got == nil || got.Obstacle != b {
		t.Error("Mob.Handle(Move) send incorrect Collide")
	}
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

	var gotA *SetOccupant
	a.Components = append(a.Components, Handler(func(e *Tile, v *SetOccupant) {
		gotA = v
	}))
	var gotB *SetOccupant
	b.Components = append(b.Components, Handler(func(e *Tile, v *SetOccupant) {
		gotB = v
	}))
	var gotM *SetPos
	m.Components = append(m.Components, Handler(func(e *Mob, v *SetPos) {
		gotM = v
	}))

	m.Handle(&Move{Vec(2, 1)})

	if gotA == nil || gotA.Value != nil {
		t.Error("Mob.Handle(Move) sent incorrect SetOccupant clear}")
	}
	if gotB == nil || gotB.Value != m {
		t.Error("Mob.Handle(Move) sent incorrect SetOccupant update")
	}
	if gotM == nil || gotM.Value != b {
		t.Error("Mob.Handle(Move) sent incorrect SetPos update")
	}
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
