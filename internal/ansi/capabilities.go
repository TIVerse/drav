package ansi

import (
	"os"
	"strings"
)

// ColorMode represents the terminal color capability.
type ColorMode int

const (
	// ColorModeNone means no color support.
	ColorModeNone ColorMode = iota
	// ColorMode16 means 16-color support.
	ColorMode16
	// ColorMode256 means 256-color support.
	ColorMode256
	// ColorModeTrueColor means 24-bit color support.
	ColorModeTrueColor
)

// DetectColorMode detects the terminal's color capability.
func DetectColorMode() ColorMode {
	colorTerm := os.Getenv("COLORTERM")
	if colorTerm == "truecolor" || colorTerm == "24bit" {
		return ColorModeTrueColor
	}

	term := os.Getenv("TERM")
	if strings.Contains(term, "256color") {
		return ColorMode256
	}

	if strings.Contains(term, "color") || term == "xterm" {
		return ColorMode16
	}

	return ColorModeNone
}

// SupportsColor returns whether the terminal supports color.
func SupportsColor() bool {
	return DetectColorMode() != ColorModeNone
}

// SupportsTrueColor returns whether the terminal supports true color.
func SupportsTrueColor() bool {
	return DetectColorMode() == ColorModeTrueColor
}

// Supports256Color returns whether the terminal supports 256 colors.
func Supports256Color() bool {
	mode := DetectColorMode()
	return mode == ColorMode256 || mode == ColorModeTrueColor
}
