package hjkl

import (
	"reflect"
	"testing"
)

type MockComponent struct {
	F  func(Event) bool
	Ok bool
}

func (c *MockComponent) Process(e Entity, v Event) {
	if c.F(v) {
		c.Ok = true
	}
}

func TestMobComponent_OpenTriggerMove(t *testing.T) {
	var mob Entity

	srcUpdate := &MockComponent{F: func(v Event) bool {
		return reflect.DeepEqual(v, &OccupantUpdate{nil})
	}}
	src := NewTile(Vector{0, 0}, srcUpdate)
	dstUpdate := &MockComponent{F: func(v Event) bool {
		return reflect.DeepEqual(v, &OccupantUpdate{mob})
	}}
	dst := NewTile(Vector{1, 1}, dstUpdate)
	src.Handle(&AdjacentUpdate{Vector{1, 1}, dst})

	mobUpdate := &MockComponent{F: func(v Event) bool {
		return reflect.DeepEqual(v, &PosUpdate{dst})
	}}
	mob = NewMob(Ch('@'), mobUpdate)
	mob.Handle(&PosUpdate{src})
	src.Handle(&OccupantUpdate{mob})

	mob.Handle(&MoveTrigger{Vector{1, 1}})

	if !srcUpdate.Ok {
		t.Error("MoveTrigger failed to send OccupantUpdate to src")
	}
	if (OccupantQuery{}).Send(src) != nil {
		t.Error("MoveTrigger failed to update src.Occupant")
	}
	if !dstUpdate.Ok {
		t.Error("MoveTrigger failed to send OccupantUpdate to dst")
	}
	if !reflect.DeepEqual(OccupantQuery{}.Send(dst), mob) {
		t.Error("MoveTrigger failed to update dst.Occupant")
	}
	if !mobUpdate.Ok {
		t.Error("MoveTrigger failed to send PosUpdate to mob")
	}
	if !reflect.DeepEqual(PosQuery{}.Send(mob), dst) {
		t.Error("MoveTrigger failed to update dst.Occupant")
	}
}

func TestMobComponent_CollideTriggerMove(t *testing.T) {
	var mob Entity

	var src Entity = NewTile(Vector{0, 0})
	var dst Entity = NewTile(Vector{1, 1})
	dst.Handle(&PassUpdate{false})
	src.Handle(&AdjacentUpdate{Vector{1, 1}, dst})

	mobUpdate := &MockComponent{F: func(v Event) bool {
		return reflect.DeepEqual(v, &CollisionEvent{dst})
	}}
	mob = NewMob(Ch('@'), mobUpdate)
	mob.Handle(&PosUpdate{src})
	src.Handle(&OccupantUpdate{mob})

	mob.Handle(&MoveTrigger{Vector{1, 1}})

	if !reflect.DeepEqual(OccupantQuery{}.Send(src), mob) {
		t.Error("MoveTrigger erroneously updated src.Occupant")
	}
	if (OccupantQuery{}).Send(dst) != nil {
		t.Error("MoveTrigger erroneously updated dst.Occupant")
	}
	if !mobUpdate.Ok {
		t.Error("MoveTrigger failed to send CollisionEvent to mob")
	}
	if !reflect.DeepEqual(PosQuery{}.Send(mob), src) {
		t.Error("MoveTrigger erroneously updated dst.Occupant")
	}
}

func TestMobComponent_BumpTriggerMove(t *testing.T) {
	var mob Entity

	bumped := NewMob(Ch('D'))

	var src Entity = NewTile(Vector{0, 0})
	var dst Entity = NewTile(Vector{1, 1})
	dst.Handle(&OccupantUpdate{bumped})
	src.Handle(&AdjacentUpdate{Vector{1, 1}, dst})

	mobUpdate := &MockComponent{F: func(v Event) bool {
		return reflect.DeepEqual(v, &BumpEvent{bumped})
	}}
	mob = NewMob(Ch('@'), mobUpdate)
	mob.Handle(&PosUpdate{src})
	src.Handle(&OccupantUpdate{mob})

	mob.Handle(&MoveTrigger{Vector{1, 1}})

	if !reflect.DeepEqual(OccupantQuery{}.Send(src), mob) {
		t.Error("MoveTrigger erroneously updated src.Occupant")
	}
	if !reflect.DeepEqual(OccupantQuery{}.Send(dst), bumped) {
		t.Error("MoveTrigger erroneously updated dst.Occupant")
	}
	if !mobUpdate.Ok {
		t.Error("MoveTrigger failed to send CollisionEvent to mob")
	}
	if !reflect.DeepEqual(PosQuery{}.Send(mob), src) {
		t.Error("MoveTrigger erroneously updated dst.Occupant")
	}
}

func TestMobComponent_FaceQuery(t *testing.T) {
	expected := Ch('@')
	c := &MobComponent{Face: expected}
	v := FaceQuery{}
	c.Process(nil, &v)
	if v.Response != expected {
		t.Error("MobComponent failed to process FaceQuery")
	}
}

func TestMobComponent_FaceUpdate(t *testing.T) {
	expected := Ch('@')
	c := &MobComponent{Face: Ch('O')}
	c.Process(nil, &FaceUpdate{expected})
	if c.Face != expected {
		t.Error("MobComponent failed to process FaceUpdate")
	}
}

func TestMobComponent_PosQuery(t *testing.T) {
	expected := &ComponentSlice{}
	c := &MobComponent{Pos: expected}
	v := PosQuery{}
	c.Process(nil, &v)
	if v.Response != expected {
		t.Error("MobComponent failed to process PosQuery")
	}
}

func TestMobComponent_PosUpdate(t *testing.T) {
	expected := &ComponentSlice{}
	c := &MobComponent{Pos: nil}
	c.Process(nil, &PosUpdate{expected})
	if c.Pos != expected {
		t.Error("MobComponent failed to process PosUpdate")
	}
}

func TestMobComponent_OffsetQuery(t *testing.T) {
	expected := Vector{420, 69}
	c := &MobComponent{
		Pos: ComponentSlice{
			ComponentFunc(func(_ Entity, v Event) {
				if v, ok := v.(*OffsetQuery); ok {
					v.Response = expected
				}
			}),
		},
	}
	v := OffsetQuery{}
	c.Process(nil, &v)
	if v.Response != expected {
		t.Error("MobComponent failed to process OffsetQuery")
	}
}

func TestTileComponent_FaceQuery(t *testing.T) {
	expected := Ch('#')
	c := &TileComponent{Face: expected}
	v := FaceQuery{}
	c.Process(nil, &v)
	if v.Response != expected {
		t.Error("TileComponent failed to process FaceQuery")
	}
}

func TestTileComponent_OccupiedFaceQuery(t *testing.T) {
	expected := Ch('@')
	c := &TileComponent{
		Face: Ch('#'),
		Occupant: ComponentSlice{
			ComponentFunc(func(_ Entity, v Event) {
				if v, ok := v.(*FaceQuery); ok {
					v.Response = expected
				}
			}),
		},
	}
	v := FaceQuery{}
	c.Process(nil, &v)
	if v.Response != expected {
		t.Error("TileComponent failed to process FaceQuery while occupied")
	}
}

func TestTileComponent_FaceUpdate(t *testing.T) {
	expected := Ch('#')
	c := &TileComponent{Face: Ch('.')}
	c.Process(nil, &FaceUpdate{expected})
	if c.Face != expected {
		t.Error("TileComponent failed to process FaceUpdate")
	}
}

func TestTileComponent_OccupantQuery(t *testing.T) {
	expected := &ComponentSlice{}
	c := &TileComponent{Occupant: expected}
	v := OccupantQuery{}
	c.Process(nil, &v)
	if v.Response != expected {
		t.Error("TileComponent failed to process OccupantQuery")
	}
}

func TestTileComponent_OccupantUpdate(t *testing.T) {
	expected := &ComponentSlice{}
	c := &TileComponent{Occupant: nil}
	c.Process(nil, &OccupantUpdate{expected})
	if c.Occupant != expected {
		t.Error("TileComponent failed to process OccupantUpdate")
	}
}

func TestTileComponent_PassQuery(t *testing.T) {
	c := &TileComponent{Pass: true}
	v := PassQuery{}
	c.Process(nil, &v)
	if !v.Response {
		t.Error("TileComponent failed to process PassQuery")
	}
}

func TestTileComponent_PassUpdate(t *testing.T) {
	c := &TileComponent{Pass: false}
	c.Process(nil, &PassUpdate{true})
	if !c.Pass {
		t.Error("TileComponent failed to process PassUpdate")
	}
}

func TestTileComponent_OffsetQuery(t *testing.T) {
	expected := Vector{420, 69}
	c := &TileComponent{Offset: expected}
	v := OffsetQuery{}
	c.Process(nil, &v)
	if v.Response != expected {
		t.Error("TileComponent failed to process OffsetQuery")
	}
}

func TestTileComponent_AdjacentQuery(t *testing.T) {
	delta := Vector{1, 1}
	expected := &ComponentSlice{}
	c := &TileComponent{
		Adjacent: map[Vector]Entity{
			delta: expected,
		},
	}
	v := AdjacentQuery{Delta: delta}
	c.Process(nil, &v)
	if v.Response != expected {
		t.Error("TileComponent failed to process AdjacentQuery")
	}
}

func TestTileComponent_AdjacentUpdate(t *testing.T) {
	delta := Vector{1, 1}
	expected := &ComponentSlice{}
	c := &TileComponent{
		Adjacent: map[Vector]Entity{
			delta: nil,
		},
	}
	c.Process(nil, &AdjacentUpdate{delta, expected})
	if c.Adjacent[delta] != expected {
		t.Error("TileComponent failed to process AdjacentUpdate")
	}
}
