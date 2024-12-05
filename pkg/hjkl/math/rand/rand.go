// Package rand provides pseudo-random number generation.
package rand

import (
	"time"
)

// sm64 is a SplitMix64 pseudo-random number generator. This isn't a good
// generator per se, but is good enough for hjkl and very fast.
var sm64 = uint64(time.Now().UnixNano())

// Seed seeds the pseudo-random number generator.
func Seed(seed uint64) {
	sm64 = seed
}

// Seed seeds the pseudo-random number generator using the current time.
func SeedTime() {
	sm64 = uint64(time.Now().UnixNano())
}

// Uint64 gets a pseudo-random uint64.
func Uint64() uint64 {
	sm64 += 0x9e3779b97f4a7c15
	x := sm64
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	return x ^ (x >> 31)
}

// Int gets a positive pseudo-random int.
func Int() int {
	return int(uint(Uint64() << 1 >> 1))
}

// Intn gets a positive pseudo-random int in [0, n). It panics if n <= 0.
func Intn(n int) int {
	if n <= 0 {
		panic("n must be positive")
	}
	return Int() % n
}

// Range gets a pseudo-random int in [a, b]. In panics if b < a.
func Range(a, b int) int {
	if b < a {
		panic("Cannot have b < a")
	}
	return Intn(b-a+1) + a
}

// Float64 gets a pseudo-random float64 in [0, 1).
func Float64() float64 {
	// Floating point numbers do not have a uniform distribution of values. We
	// get around this by only working with the 53-bit mantissa.
	return float64(Uint64()>>11) / (1 << 53)
}

// Chance gets a pseudo-random bool which is true with probability p. It
// panics if p is not in [0, 1].
func Chance(p float64) bool {
	if p < 0 || p > 1 {
		panic("p must be in [0, 1].")
	}
	return Float64() < p
}

// Choice gets a pseudo-random element of a. It panics if len(a) == 0.
func Choice[T any](a []T) T {
	if len(a) == 0 {
		panic("a cannot be empty")
	}
	return a[Intn(len(a))]
}

// Index gets a pseudo-random index into a, with probability proportional
// to the result of a weighting function. It panics if f returns a negative
// weight, or if no element of a has positive weight.
func Index[T any](a []T, f func(T) int) int {
	weights := make([]int, len(a))
	n := 0
	for i, t := range a {
		w := f(t)
		if w < 0 {
			panic("invalid negative weight")
		}
		weights[i] = w
		n += w
	}
	if n <= 0 {
		panic("no valid weights")
	}

	sample := Intn(n)
	for i, w := range weights {
		if sample < w {
			return i
		}
		sample -= w
	}

	panic("invalid state") // Should be impossible.
}

// Select returns a pseudo-random element of a for which f returns true. It
// panics if a is empty or if f returns false for all elements of a.
func Select[T any](a []T, f func(T) bool) T {
	n := len(a)
	if n == 0 {
		panic("a cannot be empty")
	}

	// Try rejection sampling for 1000 rounds.
	for i := 0; i < 1000; i++ {
		x := a[Intn(n)]
		if f(x) {
			return x
		}
	}

	// If rejection sampling fails, find every valid index.
	cands := make([]int, 0, n)
	for i, x := range a {
		if f(x) {
			cands = append(cands, i)
		}
	}

	// Return a valid element if there is one.
	if len(cands) == 0 {
		panic("no valid values in a")
	}
	return a[cands[Intn(len(cands))]]
}
