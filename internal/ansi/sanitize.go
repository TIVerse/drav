package ansi

import (
	"strings"
)

// Sanitize removes potentially dangerous ANSI escape sequences.
func Sanitize(input string) string {
	// Remove all escape sequences for security
	var result strings.Builder
	inEscape := false

	for _, ch := range input {
		if ch == '\x1b' || ch == '\x9b' {
			inEscape = true
			continue
		}

		if inEscape {
			// Skip until we find the terminator
			if (ch >= 0x40 && ch <= 0x5F) || (ch >= 0x61 && ch <= 0x7E) {
				inEscape = false
			}
			continue
		}

		result.WriteRune(ch)
	}

	return result.String()
}

// StripANSI removes all ANSI escape codes from a string.
func StripANSI(input string) string {
	return Sanitize(input)
}

// IsSafe checks if a string contains only safe characters.
func IsSafe(input string) bool {
	for _, ch := range input {
		if ch == '\x1b' || ch == '\x9b' {
			return false
		}
		if ch < 32 && ch != '\n' && ch != '\t' && ch != '\r' {
			return false
		}
	}
	return true
}
