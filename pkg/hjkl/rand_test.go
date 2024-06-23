package hjkl

import (
	"fmt"
	"math"
	"testing"
)

// X2Crit99 contains critical values of the chi-squared distribution for
// various degrees of freedom and a p-value of 0.99 for use by X2Test.
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

// X2Test runs Pearson's chi-square test with a p-value of 0.99. It returns
// true if we sustain the null hypothesis that the observed counts are
// consistent with the expected counts., and false if we reject.
func X2Test(exp, obs []int) bool {
	X2 := 0.0
	df := 0
	for i := 0; i < len(exp); i++ {
		e, o := exp[i], obs[i]
		if e == 0 {
			// We should never observe something outside support.
			if o != 0 {
				return false
			}
		} else {
			X2 += math.Pow(float64(o-e), 2) / float64(e)
			df++
		}
	}
	return X2 < X2Crit99[df]
}

// Although we're testing randomness, we want the tests to be deterministic, so
// we use a set of fixed but arbitrarily chosen seed values.
var Seeds = []uint64{
	0xBAAAAAAD,
	0x1337C0DE,
	0xFEE15BAD,
	0xABCDEF01,
	0x12345678,
}

// RunX2TestCases generates observed data using the given function, and runs
// X2Test on the observed and expected data. This process is repeated for each
// of the seeds. This isn't a thourough or rigorous statistical test, but it
// does serve as a basic check to make sure that the results are reasonable.
func RunX2TestCases(t *testing.T, name string, exp []int, f func() int) {
	n := 0
	for _, e := range exp {
		n += e
	}
	for _, seed := range Seeds {
		t.Run(fmt.Sprintf("%s (seed: %x)", name, seed), func(t *testing.T) {
			RandSeed(seed)
			obs := make([]int, len(exp))
			for i := 0; i < n; i++ {
				obs[f()]++
			}
			if !X2Test(exp, obs) {
				t.Errorf("X2Test failed. Observed: %v", obs)
			}
		})
	}
}

func TestRandIntn_10(t *testing.T) {
	exp := make([]int, 10)
	for i := 0; i < 10; i++ {
		exp[i] = 100
	}
	RunX2TestCases(t, "RandIntn(10)", exp, func() int {
		return RandIntn(10)
	})
}

func TestRandIntn_100(t *testing.T) {
	exp := make([]int, 100)
	for i := 0; i < 100; i++ {
		exp[i] = 100
	}
	RunX2TestCases(t, "RandIntn(100)", exp, func() int {
		return RandIntn(100)
	})
}

func TestRandChoice(t *testing.T) {
	a := []rune{'a', 'b', 'c', 'd'}
	index := map[rune]int{'a': 0, 'b': 1, 'c': 2, 'd': 3}
	exp := []int{250, 250, 250, 250}
	RunX2TestCases(t, "RandChoice", exp, func() int {
		return index[RandChoice(a)]
	})
}
