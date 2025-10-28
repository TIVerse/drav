package layout

// FlexLayout implements flexbox-inspired layout.
type FlexLayout struct {
	Direction  Direction
	Wrap       bool
	Justify    Justification
	Align      Alignment
	Gap        int
}

// Direction defines layout direction.
type Direction int

const (
	Row Direction = iota
	Column
)

// Justification defines how items are distributed along the main axis.
type Justification int

const (
	JustifyStart Justification = iota
	JustifyEnd
	JustifyCenter
	JustifySpaceBetween
	JustifySpaceAround
	JustifySpaceEvenly
)

// Alignment defines how items are aligned along the cross axis.
type Alignment int

const (
	AlignStart Alignment = iota
	AlignEnd
	AlignCenter
	AlignStretch
)

// Item represents a flex item.
type Item struct {
	Grow   int
	Shrink int
	Basis  int
	MinWidth  int
	MinHeight int
	MaxWidth  int
	MaxHeight int
}

// Layout computes the layout for flex items.
func (fl *FlexLayout) Layout(items []Item, containerWidth, containerHeight int) []Rect {
	if len(items) == 0 {
		return nil
	}

	rects := make([]Rect, len(items))

	if fl.Direction == Row {
		fl.layoutRow(items, containerWidth, containerHeight, rects)
	} else {
		fl.layoutColumn(items, containerWidth, containerHeight, rects)
	}

	return rects
}

// layoutRow performs horizontal layout.
func (fl *FlexLayout) layoutRow(items []Item, width, height int, rects []Rect) {
	// Calculate total basis and grow/shrink factors
	totalBasis := 0
	totalGrow := 0
	totalShrink := 0

	for _, item := range items {
		totalBasis += item.Basis
		totalGrow += item.Grow
		totalShrink += item.Shrink
	}

	// Add gaps
	totalGap := fl.Gap * (len(items) - 1)
	availableWidth := width - totalGap

	// Determine if we need to grow or shrink
	extra := availableWidth - totalBasis

	x := 0
	for i, item := range items {
		itemWidth := item.Basis

		if extra > 0 && totalGrow > 0 {
			// Distribute extra space based on grow factor
			itemWidth += (extra * item.Grow) / totalGrow
		} else if extra < 0 && totalShrink > 0 {
			// Shrink based on shrink factor
			itemWidth += (extra * item.Shrink) / totalShrink
		}

		// Apply constraints
		if item.MinWidth > 0 && itemWidth < item.MinWidth {
			itemWidth = item.MinWidth
		}
		if item.MaxWidth > 0 && itemWidth > item.MaxWidth {
			itemWidth = item.MaxWidth
		}

		rects[i] = Rect{
			X:      x,
			Y:      0,
			Width:  itemWidth,
			Height: height,
		}

		x += itemWidth + fl.Gap
	}
}

// layoutColumn performs vertical layout.
func (fl *FlexLayout) layoutColumn(items []Item, width, height int, rects []Rect) {
	// Calculate total basis and grow/shrink factors
	totalBasis := 0
	totalGrow := 0
	totalShrink := 0

	for _, item := range items {
		totalBasis += item.Basis
		totalGrow += item.Grow
		totalShrink += item.Shrink
	}

	// Add gaps
	totalGap := fl.Gap * (len(items) - 1)
	availableHeight := height - totalGap

	// Determine if we need to grow or shrink
	extra := availableHeight - totalBasis

	y := 0
	for i, item := range items {
		itemHeight := item.Basis

		if extra > 0 && totalGrow > 0 {
			// Distribute extra space based on grow factor
			itemHeight += (extra * item.Grow) / totalGrow
		} else if extra < 0 && totalShrink > 0 {
			// Shrink based on shrink factor
			itemHeight += (extra * item.Shrink) / totalShrink
		}

		// Apply constraints
		if item.MinHeight > 0 && itemHeight < item.MinHeight {
			itemHeight = item.MinHeight
		}
		if item.MaxHeight > 0 && itemHeight > item.MaxHeight {
			itemHeight = item.MaxHeight
		}

		rects[i] = Rect{
			X:      0,
			Y:      y,
			Width:  width,
			Height: itemHeight,
		}

		y += itemHeight + fl.Gap
	}
}

// Rect represents a rectangle.
type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}
