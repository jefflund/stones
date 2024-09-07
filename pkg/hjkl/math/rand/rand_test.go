package rand

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
			Seed(seed)
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

func TestIntn_10(t *testing.T) {
	exp := make([]int, 10)
	for i := 0; i < 10; i++ {
		exp[i] = 100
	}
	RunX2TestCases(t, "Intn(10)", exp, func() int {
		return Intn(10)
	})
}

func TestIntn_100(t *testing.T) {
	exp := make([]int, 100)
	for i := 0; i < 100; i++ {
		exp[i] = 100
	}
	RunX2TestCases(t, "Intn(100)", exp, func() int {
		return Intn(100)
	})
}

func TestFloat64(t *testing.T) {
	exp := make([]int, 100)
	for i := 0; i < 100; i++ {
		exp[i] = 100
	}
	RunX2TestCases(t, "Float64", exp, func() int {
		return int(Float64() * 100)
	})
}

func TestChoice(t *testing.T) {
	a := []rune{'a', 'b', 'c', 'd'}
	index := map[rune]int{'a': 0, 'b': 1, 'c': 2, 'd': 3}
	exp := []int{250, 250, 250, 250}
	RunX2TestCases(t, "Choice", exp, func() int {
		return index[Choice(a)]
	})
}

func TestRange_1D6(t *testing.T) {
	exp := []int{0, 100, 100, 100, 100, 100, 100}
	RunX2TestCases(t, "Range(1, 6)", exp, func() int {
		return Range(1, 6)
	})
}

func TestRange_10_15(t *testing.T) {
	exp := []int{100, 100, 100, 100, 100, 100}
	RunX2TestCases(t, "Range(10, 15)", exp, func() int {
		return Range(10, 15) - 10
	})
}

func TestRange_5_5(t *testing.T) {
	exp := []int{10}
	RunX2TestCases(t, "Range(5, 5)", exp, func() int {
		return Range(5, 5) - 5
	})
}

func TestChance(t *testing.T) {
	cases := []struct {
		Exp []int
		P   float64
	}{
		{[]int{1000, 0}, 0},
		{[]int{990, 10}, .01},
		{[]int{900, 100}, .1},
		{[]int{500, 500}, .5},
		{[]int{100, 900}, .9},
		{[]int{50, 950}, .95},
		{[]int{0, 1000}, 1},
	}
	for _, c := range cases {
		name := fmt.Sprintf("Chance(%.2f)", c.P)
		t.Run(name, func(t *testing.T) {
			RunX2TestCases(t, name, c.Exp, func() int {
				if Chance(c.P) {
					return 1
				}
				return 0
			})
		})
	}
}

func TestIndex(t *testing.T) {
	a := []rune{'a', 'b', 'c', 'd', 'e'}
	w := map[rune]int{
		'a': 1,
		'b': 6,
		'c': 2,
		'd': 0,
		'e': 3,
	}
	exp := []int{100, 600, 200, 0, 300}
	RunX2TestCases(t, "Index", exp, func() int {
		return Index(a, func(r rune) int {
			return w[r]
		})
	})
}

func TestSelect(t *testing.T) {
	cands := make([]int, 100)
	exp := make([]int, 100)
	for i := 0; i < 100; i++ {
		cands[i] = i
		if Chance(0.1) {
			cands[i] *= -1
		} else {
			exp[i] = 100
		}
	}
	positive := func(x int) bool { return x >= 0 }
	RunX2TestCases(t, "Select", exp, func() int {
		return Select(cands, positive)
	})
}

func BenchmarkSelect_Easy(b *testing.B) {
	cands := make([]int, 100)
	for i := range cands {
		cands[i] = i
		if Chance(0.1) {
			cands[i] *= -1
		}
	}
	positive := func(x int) bool { return x >= 0 }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Select(cands, positive)
	}
}

func BenchmarkSelect_Hard(b *testing.B) {
	cands := make([]int, 1000000)
	for i := range cands {
		cands[i] = -i
	}
	positive := func(x int) bool { return x >= 0 }
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Select(cands, positive)
	}
}

func TestIntn_Zero(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Intn(0) failed to panic")
		}
	}()
	Intn(0)
}

func TestIntn_Neg(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Intn(-1) failed to panic")
		}
	}()
	Intn(-1)
}

func TestRange_Invert(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Range failed to panic on b < a")
		}
	}()
	Range(6, 5)
}

func TestChoice_Empty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Choice failed to panic on empty")
		}
	}()
	Choice([]int{})
}

func TestChance_LT0(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Chance(-0.0001) failed to panic")
		}
	}()
	Chance(-0.0001)
}

func TestChance_GT1(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Chance(1.0001) failed to panic")
		}
	}()
	Chance(1.0001)
}

func TestSelect_Empty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Select failed to panic on empty")
		}
	}()
	Select([]int{}, func(x int) bool { return true })
}

func TestSelect_Invalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Select failed to panic on empty")
		}
	}()
	Select([]int{1, 2, 3}, func(x int) bool { return false })
}
