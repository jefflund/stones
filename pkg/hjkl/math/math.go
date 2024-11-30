// Package math provides mathematical functions.
package math

import "math"

// Hypot returns the sqrt(x * x + y * y).
func Hypot(x, y int) float64 {
	return math.Hypot(float64(x), float64(y))
}

// Min returns the minimum of two ints.
func Min(x, y int) int {
	if y < x {
		return y
	}
	return x
}

// Max returns the maximum of two ints.
func Max(x, y int) int {
	if y > x {
		return y
	}
	return x
}

// Clamp limits x to the range [a, b].
func Clamp(a, x, b int) int {
	if x < a {
		return a
	}
	if x > b {
		return b
	}
	return x
}

// Abs returns the absolute value of x.
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Sign returns 1 if x > 0, -1 if x < 0, and 0 if x = 0.
func Sign(x int) int {
	return Clamp(-1, x, 1)
}

// Mod returns x modulo y (not the same as x % y, which is remainder).
func Mod(x, y int) int {
	z := x % y
	if z < 0 {
		z += y
	}
	return z
}
