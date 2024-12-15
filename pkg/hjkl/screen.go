package hjkl

// Widget is an object which can draw itself on a Canvas.
type Widget interface {
	Draw(Canvas)
}

// Screen is a Widget composed of other Screen.
type Screen []Widget

// Screen has each constitunet Widget draw itself.
func (s Screen) Draw(c Canvas) {
	for _, w := range s {
		w.Draw(c)
	}
}

// Window serves as a base to Widget which require relative drawing.
type Window struct {
	Pos  Vector
	Size Vector
}

// RelBlit performs a Blit on the Canvas relative to the Window dimensions.
func (w *Window) RelBlit(c Canvas, v Vector, g Glyph) {
	if 0 <= v.X && v.X < w.Size.X && 0 <= v.Y && v.Y < w.Size.Y {
		c.Blit(w.Pos.Add(v), g)
	}
}

// TilesWidget is a Widget which draws a collection of Tile.
type TilesWidget struct {
	Window
	Tiles []*Tile
}

// NewTilesWidget creates a TileWidget with the given collection of Tile.
func NewTilesWidget(pos, size Vector, tiles []*Tile) *TilesWidget {
	return &TilesWidget{Window{pos, size}, tiles}
}

// Draw draws the collection of Tile.
func (w *TilesWidget) Draw(c Canvas) {
	for _, t := range w.Tiles {
		w.RelBlit(c, t.Offset, Get(t, &Face{}))
	}
}
