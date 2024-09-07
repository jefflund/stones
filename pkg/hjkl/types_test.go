package hjkl

import (
	"fmt"
	"math"
	"testing"
)

func TestVector_Add(t *testing.T) {
	cases := []struct {
		a, b, expected Vector
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
			if actual := c.a.Add(c.b); actual != c.expected {
				t.Errorf("%s = %v != %v", name, actual, c.expected)
			}
		})
	}
}

func TestVector_Sub(t *testing.T) {
	cases := []struct {
		a, b, expected Vector
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
			if actual := c.a.Sub(c.b); actual != c.expected {
				t.Errorf("%s = %v != %v", name, actual, c.expected)
			}
		})
	}
}

func TestVector_Neg(t *testing.T) {
	cases := []struct {
		vec, expected Vector
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
			if actual := c.vec.Neg(); actual != c.expected {
				t.Errorf("%s = %v != %v", name, actual, c.expected)
			}
		})
	}
}

func TestVector_Manhattan(t *testing.T) {
	cases := []struct {
		vec      Vector
		expected int
	}{
		{Vector{4, 6}, 10},
		{Vector{-4, 6}, 10},
		{Vector{4, -6}, 10},
		{Vector{-4, -6}, 10},
	}
	for _, c := range cases {
		name := fmt.Sprintf("%v.Manhattan()", c.vec)
		t.Run(name, func(t *testing.T) {
			if actual := c.vec.Manhattan(); actual != c.expected {
				t.Errorf("%s == %v != %v", name, actual, c.expected)
			}
		})
	}
}

func TestVector_Euclidean(t *testing.T) {
	cases := []struct {
		vec      Vector
		expected float64
	}{
		{Vector{3, 4}, 5},
		{Vector{-3, 4}, 5},
		{Vector{3, -4}, 5},
		{Vector{-3, -4}, 5},
		{Vector{2, 4}, math.Sqrt(20)},
		{Vector{-2, 4}, math.Sqrt(20)},
		{Vector{2, -4}, math.Sqrt(20)},
		{Vector{-2, -4}, math.Sqrt(20)},
	}
	for _, c := range cases {
		name := fmt.Sprintf("%v.Euclidean()", c.vec)
		t.Run(name, func(t *testing.T) {
			if actual := c.vec.Euclidean(); actual != c.expected {
				t.Errorf("%s == %v != %v", name, actual, c.expected)
			}
		})
	}
}

func TestVector_Chebyshev(t *testing.T) {
	cases := []struct {
		vec      Vector
		expected int
	}{
		{Vector{4, 6}, 6},
		{Vector{-4, 6}, 6},
		{Vector{4, -6}, 6},
		{Vector{-4, -6}, 6},
	}
	for _, c := range cases {
		name := fmt.Sprintf("%v.Chebyshev()", c.vec)
		t.Run(name, func(t *testing.T) {
			if actual := c.vec.Chebyshev(); actual != c.expected {
				t.Errorf("%s == %v != %v", name, actual, c.expected)
			}
		})
	}
}
