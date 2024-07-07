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

// RandChoice gets a pseudo-random element of a. It panics if len(a) == 0.
func RandChoice[T any](a []T) T {
	if len(a) == 0 {
		panic("a cannot be empty")
	}
	return a[RandIntn(len(a))]
}
