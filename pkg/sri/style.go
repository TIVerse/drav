package sri

// Style defines text styling attributes.
type Style struct {
	Foreground Color
	Background Color
	Bold       bool
	Italic     bool
	Underline  bool
	Strikethrough bool
	Reverse    bool
}

// Merge merges another style into this one.
func (s Style) Merge(other Style) Style {
	result := s

	// Only override if the other style has non-zero values
	if other.Foreground != (Color{}) {
		result.Foreground = other.Foreground
	}
	if other.Background != (Color{}) {
		result.Background = other.Background
	}
	if other.Bold {
		result.Bold = true
	}
	if other.Italic {
		result.Italic = true
	}
	if other.Underline {
		result.Underline = true
	}
	if other.Strikethrough {
		result.Strikethrough = true
	}
	if other.Reverse {
		result.Reverse = true
	}

	return result
}

// WithForeground returns a copy with the foreground color set.
func (s Style) WithForeground(color Color) Style {
	s.Foreground = color
	return s
}

// WithBackground returns a copy with the background color set.
func (s Style) WithBackground(color Color) Style {
	s.Background = color
	return s
}

// WithBold returns a copy with bold set.
func (s Style) WithBold(bold bool) Style {
	s.Bold = bold
	return s
}

// WithItalic returns a copy with italic set.
func (s Style) WithItalic(italic bool) Style {
	s.Italic = italic
	return s
}

// WithUnderline returns a copy with underline set.
func (s Style) WithUnderline(underline bool) Style {
	s.Underline = underline
	return s
}

// WithReverse returns a copy with reverse set.
func (s Style) WithReverse(reverse bool) Style {
	s.Reverse = reverse
	return s
}

// StyleBuilder provides a fluent API for building styles.
type StyleBuilder struct {
	style Style
}

// NewStyleBuilder creates a new style builder.
func NewStyleBuilder() *StyleBuilder {
	return &StyleBuilder{}
}

// Foreground sets the foreground color.
func (sb *StyleBuilder) Foreground(color Color) *StyleBuilder {
	sb.style.Foreground = color
	return sb
}

// Background sets the background color.
func (sb *StyleBuilder) Background(color Color) *StyleBuilder {
	sb.style.Background = color
	return sb
}

// Bold sets bold.
func (sb *StyleBuilder) Bold() *StyleBuilder {
	sb.style.Bold = true
	return sb
}

// Italic sets italic.
func (sb *StyleBuilder) Italic() *StyleBuilder {
	sb.style.Italic = true
	return sb
}

// Underline sets underline.
func (sb *StyleBuilder) Underline() *StyleBuilder {
	sb.style.Underline = true
	return sb
}

// Reverse sets reverse.
func (sb *StyleBuilder) Reverse() *StyleBuilder {
	sb.style.Reverse = true
	return sb
}

// Build returns the built style.
func (sb *StyleBuilder) Build() Style {
	return sb.style
}
