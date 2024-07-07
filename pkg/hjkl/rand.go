package hjkl

import "time"

// lcg is a linear congruential generator. This isn't a good generator per se,
// but its good enough for hjkl and very fast. More to the point, we'll use it
// to support a different API for hjkl than math/rand provides.
var lcg = uint64(time.Now().UnixNano())

// RandSeed seeds the pseudo-random number generator.
func RandSeed(seed uint64) {
	lcg = seed
}

// RandSeed seeds the pseudo-random number generator using the current time.
func RandSeedTime() {
	lcg = uint64(time.Now().UnixNano())
}

// RandUint64 gets a pseudo-random uint64.
func RandUint64() uint64 {
	// Constants borrowed from MMIX by Donald Knuth.
	lcg = lcg*6364136223846793005 + 1442695040888963407
	return lcg
}

// RandInt gets a positive pseudo-random int.
func RandInt() int {
	return int(uint(RandUint64() << 1 >> 1))
}

// RandIntn gets a positive pseudo-random int in [0, n). It panics if n <= 0.
func RandIntn(n int) int {
	if n <= 0 {
		panic("n must be positive")
	}
	return RandInt() % n
}

// RandFloat64 gets a pseudo-random float64 in [0, 1).
func RandFloat64() float64 {
	for {
		f := float64(RandInt()) / (1 << 63)
		if f < 1 {
			return f
		}
	}
}

// RandChance gets a pseudo-random bool which is true with probability p. It
// panics if p is not in [0, 1].
func RandChance(p float64) bool {
	if p < 0 || p > 1 {
		panic("p must be in [0, 1].")
	}
	return RandFloat64() < p
}

// RandChoice gets a pseudo-random element of a. It panics if len(a) == 0.
func RandChoice[T any](a []T) T {
	if len(a) == 0 {
		panic("a cannot be empty")
	}
	return a[RandIntn(len(a))]
}

// RandSelect returns a pseudo-random element of a for which f returns true. It
// panics if a is empty or if f returns false for all elements of a.
func RandSelect[T any](a []T, f func(T) bool) T {
	n := len(a)
	if n == 0 {
		panic("a cannot be empty")
	}

	// Try rejection sampling for 1000 rounds.
	for i := 0; i < 1000; i++ {
		x := a[RandIntn(n)]
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
	return a[cands[RandIntn(len(cands))]]
}
