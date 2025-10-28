package sri

// Theme represents a UI theme.
type Theme struct {
	Name    string
	Palette Palette
	Styles  map[string]Style
}

// NewTheme creates a new theme.
func NewTheme(name string, palette Palette) *Theme {
	return &Theme{
		Name:    name,
		Palette: palette,
		Styles:  make(map[string]Style),
	}
}

// GetStyle retrieves a style by name.
func (t *Theme) GetStyle(name string) (Style, bool) {
	style, exists := t.Styles[name]
	return style, exists
}

// SetStyle sets a style.
func (t *Theme) SetStyle(name string, style Style) {
	t.Styles[name] = style
}

// DefaultDark returns the default dark theme.
func DefaultDark() *Theme {
	palette := Palette{
		Primary:    RGB(99, 102, 241),   // Indigo
		Secondary:  RGB(139, 92, 246),   // Purple
		Success:    RGB(34, 197, 94),    // Green
		Warning:    RGB(251, 191, 36),   // Amber
		Error:      RGB(239, 68, 68),    // Red
		Info:       RGB(59, 130, 246),   // Blue
		Background: RGB(17, 24, 39),     // Dark gray
		Surface:    RGB(31, 41, 55),     // Slightly lighter
		Text:       RGB(243, 244, 246),  // Light gray
		TextMuted:  RGB(156, 163, 175),  // Medium gray
		Border:     RGB(75, 85, 99),     // Border gray
	}

	theme := NewTheme("default-dark", palette)

	// Define common styles
	theme.SetStyle("normal", Style{
		Foreground: palette.Text,
		Background: palette.Background,
	})

	theme.SetStyle("primary", Style{
		Foreground: palette.Primary,
		Background: palette.Background,
		Bold:       true,
	})

	theme.SetStyle("error", Style{
		Foreground: palette.Error,
		Background: palette.Background,
		Bold:       true,
	})

	theme.SetStyle("success", Style{
		Foreground: palette.Success,
		Background: palette.Background,
	})

	theme.SetStyle("focused", Style{
		Foreground: palette.Text,
		Background: palette.Primary,
		Bold:       true,
	})

	return theme
}

// DefaultLight returns the default light theme.
func DefaultLight() *Theme {
	palette := Palette{
		Primary:    RGB(79, 70, 229),    // Indigo
		Secondary:  RGB(124, 58, 237),   // Purple
		Success:    RGB(22, 163, 74),    // Green
		Warning:    RGB(217, 119, 6),    // Amber
		Error:      RGB(220, 38, 38),    // Red
		Info:       RGB(37, 99, 235),    // Blue
		Background: RGB(255, 255, 255),  // White
		Surface:    RGB(249, 250, 251),  // Light gray
		Text:       RGB(17, 24, 39),     // Dark gray
		TextMuted:  RGB(107, 114, 128),  // Medium gray
		Border:     RGB(209, 213, 219),  // Border gray
	}

	theme := NewTheme("default-light", palette)

	// Define common styles
	theme.SetStyle("normal", Style{
		Foreground: palette.Text,
		Background: palette.Background,
	})

	theme.SetStyle("primary", Style{
		Foreground: palette.Primary,
		Background: palette.Background,
		Bold:       true,
	})

	theme.SetStyle("error", Style{
		Foreground: palette.Error,
		Background: palette.Background,
		Bold:       true,
	})

	theme.SetStyle("success", Style{
		Foreground: palette.Success,
		Background: palette.Background,
	})

	theme.SetStyle("focused", Style{
		Foreground: palette.Background,
		Background: palette.Primary,
		Bold:       true,
	})

	return theme
}
