package habilis

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/jefflund/stones/pkg/hjkl"
)

func TestSkin_Count(t *testing.T) {
	s := Skin{
		"Grog", []Circle{
			NewCircle("Core", StoneCore, 3),
			NewCircle("Rogok", StoneDmg, 1),
			NewCircle("Warrior", StoneMelee, 1),
			NewCircle("Tough", StoneArm, 1),
		},
	}
	cases := []struct {
		Stone    Stone
		Expected int
	}{
		{StoneAny, 6},
		{StoneCore, 3},
		{StoneDmg, 2},
		{StoneArm, 2},
		{StoneHit, 1},
		{StoneEvs, 1},
		{StoneNone, 0},
	}
	for _, c := range cases {
		name := fmt.Sprintf("Skin.Count(%o)", c.Stone)
		t.Run(name, func(t *testing.T) {
			if actual := s.Count(c.Stone); actual != c.Expected {
				t.Errorf("%s got %d", name, actual)
			}
		})
	}
}

func BakeRNG[T any](seed uint64, n int, f func() T) []T {
	hjkl.RandSeed(seed)
	baked := make([]T, n)
	for i := 0; i < n; i++ {
		baked[i] = f()
	}
	hjkl.RandSeed(seed)
	return baked
}

func TestSkin_RollHit(t *testing.T) {
	const n = 5
	s := Skin{
		"Saber-tooth", []Circle{
			NewCircle("Core", StoneCore, 6),
			NewCircle("Bite", StoneMelee, 2),
		},
	}
	expected := BakeRNG(0x12345, n, func() int {
		return hjkl.RandRange(0, 6) + 2
	})
	actual := make([]int, n)
	for i := 0; i < n; i++ {
		actual[i] = s.Roll(StoneHit)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Error("Skin.Roll gave incorrect result")
	}
}

func TestSkin_RollNone(t *testing.T) {
	const n = 5
	s := Skin{
		"Saber-tooth", []Circle{
			NewCircle("Core", StoneCore, 6),
			NewCircle("Bite", StoneMelee, 2),
		},
	}
	expected := BakeRNG(0x12345, n, func() int {
		return hjkl.RandRange(0, 6)
	})
	actual := make([]int, n)
	for i := 0; i < n; i++ {
		actual[i] = s.Roll(StoneNone)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Error("Skin.Roll gave incorrect result")
	}
}
