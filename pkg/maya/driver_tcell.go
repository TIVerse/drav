package maya

import (
	"github.com/gdamore/tcell/v2"
)

// TcellDriver implements the Driver interface using tcell.
type TcellDriver struct {
	screen tcell.Screen
}

// NewTcellDriver creates a new tcell driver.
func NewTcellDriver() *TcellDriver {
	return &TcellDriver{}
}

// Init initializes the tcell screen.
func (d *TcellDriver) Init() error {
	screen, err := tcell.NewScreen()
	if err != nil {
		return err
	}

	if err := screen.Init(); err != nil {
		return err
	}

	screen.EnableMouse()
	screen.EnablePaste()
	screen.Clear()

	d.screen = screen
	return nil
}

// Size returns the terminal size.
func (d *TcellDriver) Size() (width, height int) {
	if d.screen == nil {
		return 80, 24
	}
	return d.screen.Size()
}

// SetCell sets a cell at the given position.
func (d *TcellDriver) SetCell(x, y int, cell Cell) error {
	if d.screen == nil {
		return nil
	}

	style := tcell.StyleDefault

	// Set foreground color
	if !cell.Foreground.Default {
		style = style.Foreground(tcell.NewRGBColor(
			int32(cell.Foreground.R),
			int32(cell.Foreground.G),
			int32(cell.Foreground.B),
		))
	}

	// Set background color
	if !cell.Background.Default {
		style = style.Background(tcell.NewRGBColor(
			int32(cell.Background.R),
			int32(cell.Background.G),
			int32(cell.Background.B),
		))
	}

	// Apply attributes
	if cell.Bold {
		style = style.Bold(true)
	}
	if cell.Italic {
		style = style.Italic(true)
	}
	if cell.Underline {
		style = style.Underline(true)
	}
	if cell.Reverse {
		style = style.Reverse(true)
	}

	d.screen.SetContent(x, y, cell.Rune, nil, style)
	return nil
}

// Show syncs the screen.
func (d *TcellDriver) Show() error {
	if d.screen == nil {
		return nil
	}
	d.screen.Show()
	return nil
}

// Clear clears the screen.
func (d *TcellDriver) Clear() error {
	if d.screen == nil {
		return nil
	}
	d.screen.Clear()
	return nil
}

// Shutdown shuts down the screen.
func (d *TcellDriver) Shutdown() error {
	if d.screen != nil {
		d.screen.Fini()
		d.screen = nil
	}
	return nil
}

// Screen returns the underlying tcell screen.
func (d *TcellDriver) Screen() tcell.Screen {
	return d.screen
}
