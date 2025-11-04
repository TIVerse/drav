package ansi

import (
	"strings"
)

// Sanitize removes potentially dangerous ANSI escape sequences.
// It handles CSI, OSC, DCS, and other control sequences to prevent
// terminal injection attacks.
func Sanitize(input string) string {
	var result strings.Builder
	i := 0
	runes := []rune(input)

	for i < len(runes) {
		ch := runes[i]

		// ESC (0x1b) starts an escape sequence
		if ch == '\x1b' {
			if i+1 < len(runes) {
				next := runes[i+1]
				switch next {
				case '[':
					// CSI sequence: ESC [ ... [final byte]
					i = skipCSISequence(runes, i+2)
					continue
				case ']':
					// OSC sequence: ESC ] ... BEL or ESC \
					i = skipOSCSequence(runes, i+2)
					continue
				case 'P':
					// DCS sequence: ESC P ... ESC \
					i = skipDCSSequence(runes, i+2)
					continue
				case '^', '_':
					// PM/APC sequences: ESC ^ ... ESC \ or ESC _ ... ESC \
					i = skipStringSequence(runes, i+2)
					continue
				default:
					// Two-byte escape sequence
					if next >= 0x40 && next <= 0x5F {
						i += 2
						continue
					}
				}
			}
			// Skip lone ESC
			i++
			continue
		}

		// C1 control codes (0x80-0x9F) - CSI alternative (0x9b)
		if ch == '\x9b' {
			i = skipCSISequence(runes, i+1)
			continue
		}

		// Filter other dangerous control characters (except safe ones)
		if ch < 32 {
			if ch == '\n' || ch == '\t' || ch == '\r' {
				result.WriteRune(ch)
			}
			// Skip other control characters
			i++
			continue
		}

		// Safe character
		result.WriteRune(ch)
		i++
	}

	return result.String()
}

// skipCSISequence skips a CSI (Control Sequence Introducer) sequence.
// Format: CSI [parameters] [intermediate bytes] [final byte]
func skipCSISequence(runes []rune, start int) int {
	i := start
	for i < len(runes) {
		ch := runes[i]
		// Parameter bytes: 0x30-0x3F
		// Intermediate bytes: 0x20-0x2F
		// Final byte: 0x40-0x7E
		if ch >= 0x40 && ch <= 0x7E {
			return i + 1 // Found final byte
		}
		if ch < 0x20 || ch > 0x7E {
			return i // Invalid sequence, stop here
		}
		i++
	}
	return i
}

// skipOSCSequence skips an OSC (Operating System Command) sequence.
// Format: OSC ... BEL or OSC ... ESC \
func skipOSCSequence(runes []rune, start int) int {
	i := start
	for i < len(runes) {
		ch := runes[i]
		// Terminated by BEL (0x07)
		if ch == '\x07' {
			return i + 1
		}
		// Terminated by ST (ESC \)
		if ch == '\x1b' && i+1 < len(runes) && runes[i+1] == '\\' {
			return i + 2
		}
		i++
	}
	return i
}

// skipDCSSequence skips a DCS (Device Control String) sequence.
// Format: DCS ... ESC \
func skipDCSSequence(runes []rune, start int) int {
	return skipStringSequence(runes, start)
}

// skipStringSequence skips a string-type escape sequence (DCS, PM, APC).
// Terminated by ST (ESC \)
func skipStringSequence(runes []rune, start int) int {
	i := start
	for i < len(runes) {
		ch := runes[i]
		// Terminated by ST (ESC \)
		if ch == '\x1b' && i+1 < len(runes) && runes[i+1] == '\\' {
			return i + 2
		}
		i++
	}
	return i
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
