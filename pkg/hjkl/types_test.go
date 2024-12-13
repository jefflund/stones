package hjkl

import (
	"fmt"
	"testing"
)

func TestVec(t *testing.T) {
	want := Vector{4, 2}
	got := Vec(4, 2)
	if want != got {
		t.Errorf("Vec(4, 2) = %v != %v", got, want)
	}
}

func TestVector_Add(t *testing.T) {
	cases := []struct {
		a, b, want Vector
	}{
		{Vector{4, 2}, Vector{2, -1}, Vector{6, 1}},
		{Vector{4, -2}, Vector{2, -1}, Vector{6, -3}},
		{Vector{0, 2}, Vector{2, -1}, Vector{2, 1}},
		{Vector{-3, 2}, Vector{2, -1}, Vector{-1, 1}},
		{Vector{-3, 2}, Vector{0, 0}, Vector{-3, 2}},
		{Vector{0, 0}, Vector{4, 9}, Vector{4, 9}},
	}
	for _, c := range cases {
		name := fmt.Sprintf("%v.Add(%v)", c.a, c.b)
		t.Run(name, func(t *testing.T) {
			if got := c.a.Add(c.b); got != c.want {
				t.Errorf("%s = %v != %v", name, got, c.want)
			}
		})
	}
}

func TestVector_Sub(t *testing.T) {
	cases := []struct {
		a, b, want Vector
	}{
		{Vector{4, 2}, Vector{2, -1}, Vector{2, 3}},
		{Vector{4, -2}, Vector{2, -1}, Vector{2, -1}},
		{Vector{0, 2}, Vector{2, -1}, Vector{-2, 3}},
		{Vector{-3, 2}, Vector{2, -1}, Vector{-5, 3}},
		{Vector{-3, 2}, Vector{0, 0}, Vector{-3, 2}},
		{Vector{0, 0}, Vector{4, 9}, Vector{-4, -9}},
	}
	for _, c := range cases {
		name := fmt.Sprintf("%v.Sub(%v)", c.a, c.b)
		t.Run(name, func(t *testing.T) {
			if got := c.a.Sub(c.b); got != c.want {
				t.Errorf("%s = %v != %v", name, got, c.want)
			}
		})
	}
}

func TestVector_Neg(t *testing.T) {
	cases := []struct {
		vec, want Vector
	}{
		{Vector{4, 2}, Vector{-4, -2}},
		{Vector{4, -2}, Vector{-4, 2}},
		{Vector{0, 2}, Vector{0, -2}},
		{Vector{-3, 2}, Vector{3, -2}},
		{Vector{0, 0}, Vector{0, 0}},
	}
	for _, c := range cases {
		name := fmt.Sprintf("%v.Neg()", c.vec)
		t.Run(name, func(t *testing.T) {
			if got := c.vec.Neg(); got != c.want {
				t.Errorf("%s = %v != %v", name, got, c.want)
			}
		})
	}
}

func TestCh(t *testing.T) {
	want := Glyph{'@', ColorWhite, ColorBlack}
	got := Ch('@')
	if want != got {
		t.Errorf("Ch('@') = %v != %v", got, want)
	}
}

func TestChFg(t *testing.T) {
	want := Glyph{'D', ColorRed, ColorBlack}
	got := ChFg('D', ColorRed)
	if want != got {
		t.Errorf("ChFg('@', ColorRed) = %v != %v", got, want)
	}
}
