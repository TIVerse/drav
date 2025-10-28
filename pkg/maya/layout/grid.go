package layout

// GridLayout implements a grid layout system.
type GridLayout struct {
	Rows    int
	Columns int
	RowGap  int
	ColGap  int
}

// GridItem represents an item in the grid.
type GridItem struct {
	Row    int
	Column int
	RowSpan    int
	ColSpan    int
}

// Layout computes grid layout.
func (gl *GridLayout) Layout(items []GridItem, width, height int) []Rect {
	if len(items) == 0 || gl.Rows == 0 || gl.Columns == 0 {
		return nil
	}

	// Calculate cell dimensions
	totalRowGap := gl.RowGap * (gl.Rows - 1)
	totalColGap := gl.ColGap * (gl.Columns - 1)

	cellWidth := (width - totalColGap) / gl.Columns
	cellHeight := (height - totalRowGap) / gl.Rows

	rects := make([]Rect, len(items))

	for i, item := range items {
		// Calculate position
		x := item.Column * (cellWidth + gl.ColGap)
		y := item.Row * (cellHeight + gl.RowGap)

		// Calculate size with span
		w := cellWidth*item.ColSpan + gl.ColGap*(item.ColSpan-1)
		h := cellHeight*item.RowSpan + gl.RowGap*(item.RowSpan-1)

		rects[i] = Rect{
			X:      x,
			Y:      y,
			Width:  w,
			Height: h,
		}
	}

	return rects
}

// NewGrid creates a new grid layout.
func NewGrid(rows, columns int) *GridLayout {
	return &GridLayout{
		Rows:    rows,
		Columns: columns,
		RowGap:  0,
		ColGap:  0,
	}
}
