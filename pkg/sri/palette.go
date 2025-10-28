package sri

// Palette defines a color palette for a theme.
type Palette struct {
	Primary    Color
	Secondary  Color
	Success    Color
	Warning    Color
	Error      Color
	Info       Color
	Background Color
	Surface    Color
	Text       Color
	TextMuted  Color
	Border     Color
}

// Color represents an RGB color.
type Color struct {
	R, G, B uint8
}

// RGB creates an RGB color.
func RGB(r, g, b uint8) Color {
	return Color{R: r, G: g, B: b}
}

// Hex creates a color from a hex string (e.g., "#FF5733").
func Hex(hex string) Color {
	if len(hex) == 7 && hex[0] == '#' {
		r := hexToByte(hex[1:3])
		g := hexToByte(hex[3:5])
		b := hexToByte(hex[5:7])
		return Color{R: r, G: g, B: b}
	}
	return Color{R: 0, G: 0, B: 0}
}

// ToHex converts the color to a hex string.
func (c Color) ToHex() string {
	return "#" + byteToHex(c.R) + byteToHex(c.G) + byteToHex(c.B)
}

// Lighten lightens the color by the given percentage (0-100).
func (c Color) Lighten(percent int) Color {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	factor := float64(percent) / 100.0
	r := uint8(float64(c.R) + (255-float64(c.R))*factor)
	g := uint8(float64(c.G) + (255-float64(c.G))*factor)
	b := uint8(float64(c.B) + (255-float64(c.B))*factor)

	return Color{R: r, G: g, B: b}
}

// Darken darkens the color by the given percentage (0-100).
func (c Color) Darken(percent int) Color {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	factor := 1.0 - (float64(percent) / 100.0)
	r := uint8(float64(c.R) * factor)
	g := uint8(float64(c.G) * factor)
	b := uint8(float64(c.B) * factor)

	return Color{R: r, G: g, B: b}
}

// hexToByte converts a hex string to a byte.
func hexToByte(hex string) uint8 {
	var result uint8
	for i := 0; i < len(hex) && i < 2; i++ {
		result <<= 4
		if hex[i] >= '0' && hex[i] <= '9' {
			result |= uint8(hex[i] - '0')
		} else if hex[i] >= 'a' && hex[i] <= 'f' {
			result |= uint8(hex[i] - 'a' + 10)
		} else if hex[i] >= 'A' && hex[i] <= 'F' {
			result |= uint8(hex[i] - 'A' + 10)
		}
	}
	return result
}

// byteToHex converts a byte to a hex string.
func byteToHex(b uint8) string {
	const hexChars = "0123456789ABCDEF"
	return string([]byte{hexChars[b>>4], hexChars[b&0xF]})
}
