package layout

// ColumnLayout arranges items vertically.
type ColumnLayout struct {
	Gap   int
	Align Alignment
}

// Layout arranges items in a column.
func (cl *ColumnLayout) Layout(itemHeights []int, width int, containerHeight int) []Rect {
	if len(itemHeights) == 0 {
		return nil
	}

	rects := make([]Rect, len(itemHeights))
	y := 0

	for i, height := range itemHeights {
		rects[i] = Rect{
			X:      0,
			Y:      y,
			Width:  width,
			Height: height,
		}
		y += height + cl.Gap
	}

	return rects
}
