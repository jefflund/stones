package hjkl

// Window is a translated and cropped view of another Canvas.
type Window struct {
	Canvas Canvas
	Pos    Vector
	Size   Vector
}

// WithWindow creates a Window and passes it to the given function.
func WithWindow(c Canvas, pos, size Vector, f func(Canvas)) {
	f(Window{c, pos, size})
}

// Blit sends a Glyph to the underlying Canvas, but relative to the Window.
func (w Window) Blit(v Vector, g Glyph) {
	if 0 <= v.X && v.X < w.Size.X && 0 <= v.Y && v.Y < w.Size.Y {
		w.Canvas.Blit(w.Pos.Add(v), g)
	}
}

// DisplayTiles Blits each Tile face to a Canvas.
func DisplayTiles(c Canvas, tiles []*Tile) {
	for _, t := range tiles {
		if t.Occupant != nil {
			c.Blit(t.Offset, t.Occupant.Face)
		} else {
			c.Blit(t.Offset, t.Face)
		}
	}
}

// DisplayBorder Blits a border with the given size.
func DisplayBorder(c Canvas, cols, rows int) {
	for x := 1; x < cols-1; x++ {
		c.Blit(Vec(x, 0), Ch('-'))
		c.Blit(Vec(x, rows-1), Ch('-'))
	}
	for y := 1; y < rows-1; y++ {
		c.Blit(Vec(0, y), Ch('|'))
		c.Blit(Vec(cols-1, y), Ch('|'))
	}
	c.Blit(Vec(0, 0), Ch('+'))
	c.Blit(Vec(cols-1, 0), Ch('+'))
	c.Blit(Vec(0, rows-1), Ch('+'))
	c.Blit(Vec(cols-1, rows-1), Ch('+'))
}

// DisplayString Blits a string.
func DisplayString(c Canvas, s string) {
	x, y := 0, 0
	for _, ch := range s {
		if ch == '\n' {
			x, y = 0, y+1
		} else {
			c.Blit(Vec(x, y), Ch(ch))
			x++
		}
	}
}

type LogMessage struct {
	Message string
}

type LogWidget struct {
	MaxLen  int
	History []string
}

func (w *LogWidget) Process(m *Mob, v Event) {
	if v, ok := v.(*LogMessage); ok {
		w.History = append(w.History, v.Message)
	}
	if len(w.History) > w.MaxLen {
		w.History = w.History[1:]
	}
}

func (w *LogWidget) Display(c Canvas) {
	for y, s := range w.History {
		for x, ch := range s {
			c.Blit(Vec(x, y), Ch(ch))
		}
	}
}
