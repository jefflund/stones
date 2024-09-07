package math

import (
	"fmt"
	"testing"
)

func TestMin(t *testing.T) {
	cases := []struct {
		x, y, expected int
	}{
		{4, 7, 4},
		{-5, 3, -5},
		{-4, -3, -4},
		{4, -2, -2},
	}
	for _, c := range cases {
		name := fmt.Sprintf("Min(%d, %d)", c.x, c.y)
		t.Run(name, func(t *testing.T) {
			if actual := Min(c.x, c.y); actual != c.expected {
				t.Errorf("%s == %d != %d", name, actual, c.expected)
			}
		})
	}
}

func TestMax(t *testing.T) {
	cases := []struct {
		x, y, expected int
	}{
		{4, 7, 7},
		{-5, 3, 3},
		{-4, -3, -3},
		{4, -2, 4},
	}
	for _, c := range cases {
		name := fmt.Sprintf("Max(%d, %d)", c.x, c.y)
		t.Run(name, func(t *testing.T) {
			if actual := Max(c.x, c.y); actual != c.expected {
				t.Errorf("%s == %d != %d", name, actual, c.expected)
			}
		})
	}
}

func TestClamp(t *testing.T) {
	cases := []struct {
		a, x, b  int
		expected int
	}{
		{-1, -2, 1, -1},
		{-1, -1, 1, -1},
		{-1, 0, 1, 0},
		{-1, 1, 1, 1},
		{-1, 2, 1, 1},
		{0, -1, 10, 0},
		{0, 0, 10, 0},
		{0, 1, 10, 1},
		{0, 5, 10, 5},
		{0, 9, 10, 9},
		{0, 10, 10, 10},
		{0, 11, 10, 10},
	}
	for _, c := range cases {
		name := fmt.Sprintf("Clamp(%d, %d, %d)", c.a, c.x, c.b)
		t.Run(name, func(t *testing.T) {
			if actual := Clamp(c.a, c.x, c.b); actual != c.expected {
				t.Errorf("%s == %d != %d", name, actual, c.expected)
			}
		})
	}
}

func TestAbs(t *testing.T) {
	cases := []struct {
		x, expected int
	}{
		{1, 1},
		{2, 2},
		{0, 0},
		{-1, 1},
		{-2, 2},
	}
	for _, c := range cases {
		name := fmt.Sprintf("Abs(%d)", c.x)
		t.Run(name, func(t *testing.T) {
			if actual := Abs(c.x); actual != c.expected {
				t.Errorf("%s == %d != %d", name, actual, c.expected)
			}
		})
	}
}

func TestSign(t *testing.T) {
	cases := []struct {
		x, expected int
	}{
		{1, 1},
		{2, 1},
		{0, 0},
		{-1, -1},
		{-2, -1},
	}
	for _, c := range cases {
		name := fmt.Sprintf("Sign(%d", c.x)
		t.Run(name, func(t *testing.T) {
			if actual := Sign(c.x); actual != c.expected {
				t.Errorf("%s == %d != %d", name, actual, c.expected)
			}
		})
	}
}

func TestMod(t *testing.T) {
	cases := []struct {
		x, y     int
		expected int
	}{
		{-11, 5, 4},
		{-10, 5, 0},
		{-9, 5, 1},
		{-8, 5, 2},
		{-7, 5, 3},
		{-6, 5, 4},
		{-5, 5, 0},
		{-4, 5, 1},
		{-3, 5, 2},
		{-2, 5, 3},
		{-1, 5, 4},
		{0, 5, 0},
		{1, 5, 1},
		{2, 5, 2},
		{3, 5, 3},
		{4, 5, 4},
		{5, 5, 0},
		{6, 5, 1},
		{7, 5, 2},
		{8, 5, 3},
		{9, 5, 4},
		{10, 5, 0},
		{11, 5, 1},
	}
	for _, c := range cases {
		name := fmt.Sprintf("Mod(%d, %d)", c.x, c.y)
		t.Run(name, func(t *testing.T) {
			if actual := Mod(c.x, c.y); c.expected != actual {
				t.Errorf("%s == %d != %d", name, actual, c.expected)
			}
		})
	}
}
