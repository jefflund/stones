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
func DisplayTiles[T any](c Canvas, tiles []*Tile[T]) {
	for _, t := range tiles {
		if t.Occupant != nil {
			c.Blit(t.Offset, t.Occupant.Face)
		} else {
			c.Blit(t.Offset, t.Face)
		}
	}
}
