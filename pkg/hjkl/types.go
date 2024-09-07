package hjkl

import "github.com/jefflund/stones/pkg/hjkl/math"

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
	return math.Abs(a.X) + math.Abs(a.Y)
}

// Euclidean returns the L_2 norm of the Vector.
func (a Vector) Euclidean() float64 {
	return math.Hypot(a.X, a.Y)
}

// Chebyshev returns the L_inf norm of the Vector.
func (a Vector) Chebyshev() int {
	return math.Max(math.Abs(a.X), math.Abs(a.Y))
}

// Dirs8 contains the eight point compass directions as Vector.
var Dirs8 = []Vector{
	{-1, -1},
	{-1, 0},
	{-1, 1},
	{0, -1},
	{0, 1},
	{1, -1},
	{1, 0},
	{1, 1},
}

// Key represents a single keypress.
type Key rune

// Key constants which normally require escapes.
const (
	KeyEsc   Key = 0x1B
	KeyEnter Key = 0x0D
	KeyCtrlC Key = 0x03
)

var VIKeyMap = map[Key]Vector{
	'h': Vec(-1, 0),
	'j': Vec(0, 1),
	'k': Vec(0, -1),
	'l': Vec(1, 0),
	'n': Vec(1, 1),
	'b': Vec(-1, 1),
	'u': Vec(1, -1),
	'y': Vec(-1, -1),
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
