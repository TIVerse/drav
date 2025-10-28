package layout

// RowLayout arranges items horizontally.
type RowLayout struct {
	Gap    int
	Align  Alignment
}

// Layout arranges items in a row.
func (rl *RowLayout) Layout(itemWidths []int, height int, containerWidth int) []Rect {
	if len(itemWidths) == 0 {
		return nil
	}

	rects := make([]Rect, len(itemWidths))
	x := 0

	for i, width := range itemWidths {
		rects[i] = Rect{
			X:      x,
			Y:      0,
			Width:  width,
			Height: height,
		}
		x += width + rl.Gap
	}

	return rects
}
