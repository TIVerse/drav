package maya

import (
	"sync"
)

// Cell represents a single terminal cell.
type Cell struct {
	Rune       rune
	Foreground Color
	Background Color
	Bold       bool
	Italic     bool
	Underline  bool
	Reverse    bool
	Dirty      bool
}

// Buffer represents a 2D grid of cells.
type Buffer struct {
	mu     sync.RWMutex
	width  int
	height int
	cells  [][]Cell
}

// NewBuffer creates a new buffer.
func NewBuffer(width, height int) *Buffer {
	cells := make([][]Cell, height)
	for i := range cells {
		cells[i] = make([]Cell, width)
		for j := range cells[i] {
			cells[i][j] = Cell{Rune: ' ', Foreground: DefaultColor(), Background: DefaultColor()}
		}
	}
	return &Buffer{
		width:  width,
		height: height,
		cells:  cells,
	}
}

// Width returns the buffer width.
func (b *Buffer) Width() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.width
}

// Height returns the buffer height.
func (b *Buffer) Height() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.height
}

// Get returns the cell at the given position.
func (b *Buffer) Get(x, y int) Cell {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if y < 0 || y >= b.height || x < 0 || x >= b.width {
		return Cell{}
	}
	return b.cells[y][x]
}

// Set sets the cell at the given position.
func (b *Buffer) Set(x, y int, cell Cell) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if y < 0 || y >= b.height || x < 0 || x >= b.width {
		return
	}
	b.cells[y][x] = cell
	b.cells[y][x].Dirty = true
}

// SetRune sets the rune at the given position.
func (b *Buffer) SetRune(x, y int, r rune) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if y < 0 || y >= b.height || x < 0 || x >= b.width {
		return
	}
	b.cells[y][x].Rune = r
	b.cells[y][x].Dirty = true
}

// SetColors sets the colors at the given position.
func (b *Buffer) SetColors(x, y int, fg, bg Color) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if y < 0 || y >= b.height || x < 0 || x >= b.width {
		return
	}
	b.cells[y][x].Foreground = fg
	b.cells[y][x].Background = bg
	b.cells[y][x].Dirty = true
}

// Clear clears the buffer.
func (b *Buffer) Clear() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			b.cells[y][x] = Cell{Rune: ' ', Foreground: DefaultColor(), Background: DefaultColor(), Dirty: true}
		}
	}
}

// Resize resizes the buffer.
func (b *Buffer) Resize(width, height int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if width == b.width && height == b.height {
		return
	}

	newCells := make([][]Cell, height)
	for i := range newCells {
		newCells[i] = make([]Cell, width)
		for j := range newCells[i] {
			if i < b.height && j < b.width {
				newCells[i][j] = b.cells[i][j]
			} else {
				newCells[i][j] = Cell{Rune: ' ', Foreground: DefaultColor(), Background: DefaultColor()}
			}
			newCells[i][j].Dirty = true
		}
	}

	b.cells = newCells
	b.width = width
	b.height = height
}

// MarkAllDirty marks all cells as dirty.
func (b *Buffer) MarkAllDirty() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			b.cells[y][x].Dirty = true
		}
	}
}

// ClearDirty clears the dirty flag for all cells.
func (b *Buffer) ClearDirty() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			b.cells[y][x].Dirty = false
		}
	}
}

// Clone creates a copy of the buffer.
func (b *Buffer) Clone() *Buffer {
	b.mu.RLock()
	defer b.mu.RUnlock()

	newBuffer := NewBuffer(b.width, b.height)
	for y := 0; y < b.height; y++ {
		copy(newBuffer.cells[y], b.cells[y])
	}
	return newBuffer
}

// WriteString writes a string at the given position.
func (b *Buffer) WriteString(x, y int, s string, fg, bg Color) {
	b.mu.Lock()
	defer b.mu.Unlock()

	col := x
	for _, r := range s {
		if col >= b.width {
			break
		}
		if y >= 0 && y < b.height {
			b.cells[y][col].Rune = r
			b.cells[y][col].Foreground = fg
			b.cells[y][col].Background = bg
			b.cells[y][col].Dirty = true
		}
		col++
	}
}
