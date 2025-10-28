package maya

import (
	"hash/fnv"
)

// DiffResult represents the result of diffing two buffers.
type DiffResult struct {
	DirtyLines map[int]bool
	Changes    []CellChange
}

// CellChange represents a single cell change.
type CellChange struct {
	X    int
	Y    int
	Cell Cell
}

// Diff computes the differences between two buffers using a two-pass algorithm.
// Pass 1: Line-level hashing to detect dirty lines quickly.
// Pass 2: Cell-level diffing only on dirty lines.
func Diff(oldBuf, newBuf *Buffer) *DiffResult {
	if oldBuf == nil || newBuf == nil {
		return &DiffResult{
			DirtyLines: make(map[int]bool),
			Changes:    []CellChange{},
		}
	}

	height := newBuf.Height()
	width := newBuf.Width()

	// Resize check
	if oldBuf.Width() != width || oldBuf.Height() != height {
		// Full redraw needed
		return fullDiff(newBuf)
	}

	result := &DiffResult{
		DirtyLines: make(map[int]bool),
		Changes:    make([]CellChange, 0),
	}

	// Pass 1: Line-level hashing
	for y := 0; y < height; y++ {
		if lineHash(oldBuf, y) != lineHash(newBuf, y) {
			result.DirtyLines[y] = true
		}
	}

	// Pass 2: Cell-level diffing on dirty lines
	for y := range result.DirtyLines {
		for x := 0; x < width; x++ {
			oldCell := oldBuf.Get(x, y)
			newCell := newBuf.Get(x, y)

			if !cellsEqual(oldCell, newCell) {
				result.Changes = append(result.Changes, CellChange{
					X:    x,
					Y:    y,
					Cell: newCell,
				})
			}
		}
	}

	return result
}

// lineHash computes a hash for a line in the buffer.
func lineHash(buf *Buffer, y int) uint64 {
	h := fnv.New64a()
	width := buf.Width()
	
	for x := 0; x < width; x++ {
		cell := buf.Get(x, y)
		// Hash the rune and colors
		h.Write([]byte{
			byte(cell.Rune >> 24),
			byte(cell.Rune >> 16),
			byte(cell.Rune >> 8),
			byte(cell.Rune),
			cell.Foreground.R,
			cell.Foreground.G,
			cell.Foreground.B,
			cell.Background.R,
			cell.Background.G,
			cell.Background.B,
		})
		if cell.Bold {
			h.Write([]byte{1})
		}
		if cell.Italic {
			h.Write([]byte{2})
		}
		if cell.Underline {
			h.Write([]byte{3})
		}
		if cell.Reverse {
			h.Write([]byte{4})
		}
	}
	
	return h.Sum64()
}

// cellsEqual checks if two cells are equal.
func cellsEqual(a, b Cell) bool {
	return a.Rune == b.Rune &&
		colorsEqual(a.Foreground, b.Foreground) &&
		colorsEqual(a.Background, b.Background) &&
		a.Bold == b.Bold &&
		a.Italic == b.Italic &&
		a.Underline == b.Underline &&
		a.Reverse == b.Reverse
}

// colorsEqual checks if two colors are equal.
func colorsEqual(a, b Color) bool {
	if a.Default && b.Default {
		return true
	}
	if a.Default != b.Default {
		return false
	}
	return a.R == b.R && a.G == b.G && a.B == b.B
}

// fullDiff returns a diff result that marks everything as changed.
func fullDiff(buf *Buffer) *DiffResult {
	height := buf.Height()
	width := buf.Width()
	
	result := &DiffResult{
		DirtyLines: make(map[int]bool),
		Changes:    make([]CellChange, 0, height*width),
	}

	for y := 0; y < height; y++ {
		result.DirtyLines[y] = true
		for x := 0; x < width; x++ {
			result.Changes = append(result.Changes, CellChange{
				X:    x,
				Y:    y,
				Cell: buf.Get(x, y),
			})
		}
	}

	return result
}

// OptimizeChanges optimizes the list of changes by merging adjacent cells.
func OptimizeChanges(changes []CellChange) [][]CellChange {
	if len(changes) == 0 {
		return nil
	}

	// Group changes by line
	lineGroups := make(map[int][]CellChange)
	for _, change := range changes {
		lineGroups[change.Y] = append(lineGroups[change.Y], change)
	}

	// Merge adjacent cells on each line
	optimized := make([][]CellChange, 0, len(lineGroups))
	for _, group := range lineGroups {
		if len(group) == 0 {
			continue
		}

		// Sort by X coordinate
		sortChangesByX(group)

		// Merge adjacent cells
		merged := make([]CellChange, 0, len(group))
		merged = append(merged, group[0])

		for i := 1; i < len(group); i++ {
			lastIdx := len(merged) - 1
			if group[i].X == merged[lastIdx].X+1 {
				// Adjacent cell, can be merged in a single write
				merged = append(merged, group[i])
			} else {
				// Gap detected, start new group
				optimized = append(optimized, merged)
				merged = make([]CellChange, 0, len(group)-i)
				merged = append(merged, group[i])
			}
		}
		if len(merged) > 0 {
			optimized = append(optimized, merged)
		}
	}

	return optimized
}

// sortChangesByX sorts changes by X coordinate (simple bubble sort for small arrays).
func sortChangesByX(changes []CellChange) {
	n := len(changes)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if changes[j].X > changes[j+1].X {
				changes[j], changes[j+1] = changes[j+1], changes[j]
			}
		}
	}
}
