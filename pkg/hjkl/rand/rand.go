package rand

import "time"

// The 64-bit state of the SplitMix64 generator used by Uint64.
var state uint64

// init ensures that users don't have to manually seed the generator.
func init() {
	SeedTime()
}

// Seed seeds the random number generator to a deterministic state.
func Seed(seed uint64) {
	state = seed
}

// Seed seeds the random number generator using the current time.
func SeedTime() {
	Seed(uint64(time.Now().UnixNano()))
}

// Uint64 returns a uniform random uint64.
func Uint64() uint64 {
	// The SplitMix64 algorithm from Java 8's SplittableRandom class. This
	// isn't a good generator per se, but it is good enough to pass BigCrush,
	// and extremely fast with only 64 bits of state. We could of course reuse
	// the generator from math/rand, but we provide a different API for hjkl so
	// we might as well use a less clunky random source.
	state += 0x9e3779b97f4a7c15
	z := state
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

// Float64 returns a uniform random float64 in [0, 1).
func Float64() float64 {
	// Floating point values are not uniformly distributed, so we can't just
	// divide by 1<<64. Instead, we use just the 53 mantiass bits of float64.
	return float64(Uint64()>>11) / float64(1<<53)
}

// Intn returns a uniform random int in [0, n). It panics if n <= 0.
func Intn(n int) int {
	if n <= 0 {
		panic("Invalid argument to Intn")
	}
	return int(Uint64() % uint64(n))
}

// Chance returns true with probability p. It panics of p < 0 or p > 1.
func Chance(p float64) bool {
	if p < 0 || p > 1 {
		panic("Invalid argument to Chance")
	}
	return Float64() < p
}
