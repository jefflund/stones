package hjkl

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

// CompassDirs contains the eight point compass directions as Vector.
var CompassDirs = []Vector{
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

// VIKeys is a mapping of VI Key to CompassDirs.
var VIKeyDirs = map[Key]Vector{
	'h': {-1, 0},
	'j': {0, 1},
	'k': {0, -1},
	'l': {1, 0},
	'n': {1, 1},
	'b': {-1, 1},
	'u': {1, -1},
	'y': {-1, -1},
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

// ChFg is shorthand for Glyph{Ch: ch, Fg: fg, Bg: ColorBlack}.
func ChFg(ch rune, fg Color) Glyph {
	return Glyph{Ch: ch, Fg: fg, Bg: ColorBlack}
}
