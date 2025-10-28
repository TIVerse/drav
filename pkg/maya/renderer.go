package maya

import (
	"context"
	"sync"
)

// Renderer handles screen rendering with double-buffering.
type Renderer struct {
	mu           sync.RWMutex
	driver       Driver
	frontBuffer  *Buffer
	backBuffer   *Buffer
	width        int
	height       int
	initialized  bool
	damageTrack  map[string]bool
}

// Driver is the terminal backend interface.
type Driver interface {
	Init() error
	Size() (width, height int)
	SetCell(x, y int, cell Cell) error
	Show() error
	Clear() error
	Shutdown() error
}

// NewRenderer creates a new renderer with the given driver.
func NewRenderer(driver Driver) *Renderer {
	return &Renderer{
		driver:      driver,
		damageTrack: make(map[string]bool),
	}
}

// Init initializes the renderer.
func (r *Renderer) Init() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.initialized {
		return nil
	}

	if err := r.driver.Init(); err != nil {
		return err
	}

	width, height := r.driver.Size()
	r.width = width
	r.height = height
	r.frontBuffer = NewBuffer(width, height)
	r.backBuffer = NewBuffer(width, height)
	r.initialized = true

	return nil
}

// Render renders a view to the screen.
func (r *Renderer) Render(ctx context.Context, view View) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.initialized {
		return nil
	}

	// Clear back buffer
	r.backBuffer.Clear()

	// Render view to back buffer
	r.renderView(view, 0, 0, r.width, r.height)

	// Diff buffers
	diff := Diff(r.frontBuffer, r.backBuffer)

	// Apply changes to terminal
	for _, change := range diff.Changes {
		if err := r.driver.SetCell(change.X, change.Y, change.Cell); err != nil {
			return err
		}
	}

	// Show changes
	if err := r.driver.Show(); err != nil {
		return err
	}

	// Swap buffers
	r.frontBuffer, r.backBuffer = r.backBuffer, r.frontBuffer

	return nil
}

// renderView renders a view to a specific region of the buffer.
func (r *Renderer) renderView(view View, x, y, width, height int) {
	if view == nil {
		return
	}

	attrs := view.Attrs()

	// Apply padding
	x += attrs.Padding.Left
	y += attrs.Padding.Top
	width -= attrs.Padding.Left + attrs.Padding.Right
	height -= attrs.Padding.Top + attrs.Padding.Bottom

	if width <= 0 || height <= 0 {
		return
	}

	switch v := view.(type) {
	case *TextView:
		r.renderText(v, x, y, width, height)
	case *ContainerView:
		r.renderContainer(v, x, y, width, height)
	default:
		// Render children
		for _, child := range view.Children() {
			r.renderView(child, x, y, width, height)
		}
	}
}

// renderText renders a text view.
func (r *Renderer) renderText(tv *TextView, x, y, width, height int) {
	attrs := tv.Attrs()
	fg := attrs.Foreground
	bg := attrs.Background

	if fg.Default {
		fg = DefaultColor()
	}
	if bg.Default {
		bg = DefaultColor()
	}

	// Simple text rendering without wrapping for now
	col := x
	row := y
	for _, ch := range tv.Content {
		if ch == '\n' {
			row++
			col = x
			if row >= y+height {
				break
			}
			continue
		}

		if col >= x+width {
			break
		}

		r.backBuffer.Set(col, row, Cell{
			Rune:       ch,
			Foreground: fg,
			Background: bg,
		})
		col++
	}
}

// renderContainer renders a container view with children.
func (r *Renderer) renderContainer(cv *ContainerView, x, y, width, height int) {
	children := cv.Children()
	if len(children) == 0 {
		return
	}

	if cv.Direction == DirectionRow {
		// Horizontal layout
		childWidth := width / len(children)
		for i, child := range children {
			childX := x + i*childWidth
			r.renderView(child, childX, y, childWidth, height)
		}
	} else {
		// Vertical layout
		childHeight := height / len(children)
		for i, child := range children {
			childY := y + i*childHeight
			r.renderView(child, x, childY, width, childHeight)
		}
	}
}

// Clear clears the screen.
func (r *Renderer) Clear() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.initialized {
		return nil
	}

	r.frontBuffer.Clear()
	r.backBuffer.Clear()
	return r.driver.Clear()
}

// Resize handles terminal resize events.
func (r *Renderer) Resize(width, height int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.initialized {
		return
	}

	r.width = width
	r.height = height
	r.frontBuffer.Resize(width, height)
	r.backBuffer.Resize(width, height)
}

// Shutdown shuts down the renderer.
func (r *Renderer) Shutdown() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.initialized {
		return nil
	}

	r.initialized = false
	return r.driver.Shutdown()
}

// Size returns the current terminal size.
func (r *Renderer) Size() (width, height int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.width, r.height
}

// MarkDirty marks a component as dirty for re-rendering.
func (r *Renderer) MarkDirty(componentID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.damageTrack[componentID] = true
}

// ClearDirty clears dirty marks.
func (r *Renderer) ClearDirty() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.damageTrack = make(map[string]bool)
}
