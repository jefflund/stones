package rand

import (
	"fmt"
	"math"
	"testing"
)

// X2Crit99 gives the critical values from the chi-squared distribution for a
// p-value of .99 for various degrees of freedom.
var X2Crit99 = []float64{
	0, 6.6349, 9.2104, 11.3449, 13.2767, 15.0863, 16.8119, 18.4753, 20.0902,
	21.666, 23.2093, 24.725, 26.217, 27.6882, 29.1412, 30.578, 31.9999,
	33.4087, 34.8052, 36.1908, 37.5663, 38.9322, 40.2894, 41.6383, 42.9798,
	44.314, 45.6416, 46.9628, 48.2782, 49.5878, 50.8922, 52.1914, 53.4857,
	54.7754, 56.0609, 57.342, 58.6192, 59.8926, 61.162, 62.4281, 63.6908,
	64.95, 66.2063, 67.4593, 68.7096, 69.9569, 71.2015, 72.4432, 73.6826,
	74.9194, 76.1538, 77.386, 78.6156, 79.8434, 81.0688, 82.292, 83.5136,
	84.7327, 85.9501, 87.1658, 88.3794, 89.5912, 90.8015, 92.0099, 93.2167,
	94.422, 95.6256, 96.8277, 98.0283, 99.2274, 100.4251, 101.6214, 102.8163,
	104.0098, 105.2019, 106.3929, 107.5824, 108.7709, 109.9582, 111.144,
	112.3288, 113.5123, 114.6948, 115.8762, 117.0566, 118.2356, 119.4137,
	120.5909, 121.7672, 122.9422, 124.1162, 125.2893, 126.4616, 127.633,
	128.8032, 129.9725, 131.1411, 132.3089, 133.4756, 134.6415, 135.8069,
}

// X2Test runs Pearson's chi-squared test with a p-value of .99, returning true
// if we sustain the null hypothesis that the observed counts is consistent
// with the expected counts, and false if we reject.
func X2Test(exp, obs []int) bool {
	n := len(obs)
	X2 := 0.0
	df := 0
	for i := 0; i < n; i++ {
		if exp[i] == 0 {
			// We should never observe something outside support.
			if obs[i] != 0 {
				return false
			}
		} else {
			X2 += math.Pow(float64(obs[i]-exp[i]), 2) / float64(exp[i])
			df++ // Only increment if we're inside support.
		}
	}
	return X2 < X2Crit99[df-1]
}

// In order to ensure that the tests are deterministic, we use a set of fixed,
// but arbitrarily chosen seed values.
var Seeds = []uint64{
	0xBAAAAAAD,
	0x1337C0DE,
	0xFEE15BAD,
	0xABCDEF01,
	0x12345678,
}

// RunX2TestCases generates observed data using the given function, and runs a
// X2Test on the observed data and the expected data. This process is repeated
// for each of the seeds. This isn' a super rigorous statistical test, but it
// serves as a basic confidence test of the pseudo-random funtionality.
func RunX2TestCases(name string, t *testing.T, exp []int, f func() int) {
	n := 0
	for _, count := range exp {
		n += count
	}
	for _, seed := range Seeds {
		t.Run(fmt.Sprintf("%s (seed:%X)", name, seed), func(t *testing.T) {
			Seed(seed)
			obs := make([]int, len(exp))
			for i := 0; i < n; i++ {
				obs[f()]++
			}
			if !X2Test(exp, obs) {
				t.Errorf("Chi-squared test failed")
			}
		})
	}
}

func TestIntn_10(t *testing.T) {
	exp := make([]int, 10)
	for i := 0; i < 10; i++ {
		exp[i] = 100
	}
	RunX2TestCases("RandIntn(10)", t, exp, func() int {
		return Intn(10)
	})
}

func TestIntn_100(t *testing.T) {
	exp := make([]int, 100)
	for i := 0; i < 100; i++ {
		exp[i] = 100
	}
	RunX2TestCases("RandIntn(10)", t, exp, func() int {
		return Intn(100)
	})
}

func TestChance_75(t *testing.T) {
	exp := []int{75, 25}
	RunX2TestCases("Chance(.75)", t, exp, func() int {
		if Chance(.75) {
			return 0
		}
		return 1
	})
}

func TestChance_35(t *testing.T) {
	exp := []int{350, 650}
	RunX2TestCases("Chance(.35)", t, exp, func() int {
		if Chance(.35) {
			return 0
		}
		return 1
	})
}

func TestChance_50(t *testing.T) {
	exp := []int{500, 500}
	RunX2TestCases("Chance(.5)", t, exp, func() int {
		if Chance(.5) {
			return 0
		}
		return 1
	})
}

func TestInvalidArgs(t *testing.T) {
	cases := []struct {
		name string
		call func()
	}{
		{"Intn(0)", func() { Intn(0) }},
		{"Intn(-1)", func() { Intn(-1) }},
		{"Chance(-0.0001)", func() { Chance(-0.0001) }},
		{"Chance(1.0001)", func() { Chance(1.0001) }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("%s did not panic", c.name)
				}
			}()
			c.call()
		})
	}
}

func TestValidArgs(t *testing.T) {
	cases := []struct {
		name string
		call func()
	}{
		{"Intn(1)", func() { Intn(1) }},
		{"Chance(0)", func() { Chance(0) }},
		{"Chance(1)", func() { Chance(1) }},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s panicked", c.name)
				}
			}()
			c.call()
		})
	}
}
