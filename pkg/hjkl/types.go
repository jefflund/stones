package hjkl

import "math"

// Key represents a single keypress.
type Key rune

// Key constants which normally require escapes.
const (
	KeyEsc   Key = 0x1B
	KeyEnter Key = 0x0D
	KeyCtrlC Key = 0x03
)

// Vector is a two-dimension int vector.
type Vector struct {
	X, Y int
}

// Vec is shorthand for Vector{X: x, Y: y}.
func Vec(x, y int) Vector {
	return Vector{X: x, Y: y}
}

// Add returns the sum of adding another Vector to this one.
func (a Vector) Add(b Vector) Vector {
	return Vector{a.X + b.X, a.Y + b.Y}
}

// Sub returns the sum of subtracting another Vector from this one.
func (a Vector) Sub(b Vector) Vector {
	return Vector{a.X - b.X, a.Y - b.Y}
}

// Neg returns the result of negating the Vector.
func (a Vector) Neg() Vector {
	return Vector{-a.X, -a.Y}
}

// Manhattan returns the L_1 norm of the Vector.
func (a Vector) Manhattan() int {
	return Abs(a.X) + Abs(a.Y)
}

// Euclidean returns the L_2 norm of the Vector.
func (a Vector) Euclidean() float64 {
	return math.Hypot(float64(a.X), float64(a.Y))
}

// Chebyshev returns the L_inf norm of the Vector.
func (a Vector) Chebyshev() int {
	return Max(Abs(a.X), Abs(a.Y))
}

// Color describes the color of a Glyph as a uint8 ANSI 256 color code.
type Color uint8

// Color constants for use with Glyph.
const (
	ColorBlack Color = iota
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
	ColorLightBlack
	ColorLightRed
	ColorLightGreen
	ColorLightYellow
	ColorLightBlue
	ColorLightMagenta
	ColorLightCyan
	ColorLightWhite
)

// Glyph represents a single onscreen character.
type Glyph struct {
	Ch rune
	Fg Color
	Bg Color
}

// Ch is shorthand for Glyph{Ch: ch, Fg: ColorWhite, Bg: ColorBlack}.
func Ch(ch rune) Glyph {
	return Glyph{Ch: ch, Fg: ColorWhite, Bg: ColorBlack}
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
