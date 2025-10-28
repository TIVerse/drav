package layout

// StackLayout stacks items on top of each other (z-index based).
type StackLayout struct {
	Align   Alignment
	Justify Justification
}

// StackItem represents a stacked item.
type StackItem struct {
	Width  int
	Height int
	ZIndex int
}

// Layout arranges items in a stack.
func (sl *StackLayout) Layout(items []StackItem, width, height int) []Rect {
	if len(items) == 0 {
		return nil
	}

	rects := make([]Rect, len(items))

	for i, item := range items {
		// Center by default (can be customized based on Align/Justify)
		x := (width - item.Width) / 2
		y := (height - item.Height) / 2

		if x < 0 {
			x = 0
		}
		if y < 0 {
			y = 0
		}

		rects[i] = Rect{
			X:      x,
			Y:      y,
			Width:  item.Width,
			Height: item.Height,
		}
	}

	// Sort by z-index (items are rendered in order)
	sortByZIndex(items, rects)

	return rects
}

// sortByZIndex sorts items by z-index using bubble sort.
func sortByZIndex(items []StackItem, rects []Rect) {
	n := len(items)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if items[j].ZIndex > items[j+1].ZIndex {
				items[j], items[j+1] = items[j+1], items[j]
				rects[j], rects[j+1] = rects[j+1], rects[j]
			}
		}
	}
}
