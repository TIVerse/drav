package ansi

import (
	"strings"
	"testing"
)

func TestSanitize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain text",
			input:    "Hello, World!",
			expected: "Hello, World!",
		},
		{
			name:     "CSI sequence - color",
			input:    "\x1b[31mRed Text\x1b[0m",
			expected: "Red Text",
		},
		{
			name:     "CSI alternative - C1 control (as UTF-8 rune)",
			input:    string(rune(0x9b)) + "31m" + "Red Text" + string(rune(0x9b)) + "0m",
			expected: "Red Text",
		},
		{
			name:     "OSC sequence with BEL terminator",
			input:    "\x1b]0;Window Title\x07Normal Text",
			expected: "Normal Text",
		},
		{
			name:     "OSC sequence with ST terminator",
			input:    "\x1b]0;Window Title\x1b\\Normal Text",
			expected: "Normal Text",
		},
		{
			name:     "DCS sequence",
			input:    "\x1bPDCS data\x1b\\Normal Text",
			expected: "Normal Text",
		},
		{
			name:     "PM sequence",
			input:    "\x1b^PM data\x1b\\Normal Text",
			expected: "Normal Text",
		},
		{
			name:     "APC sequence",
			input:    "\x1b_APC data\x1b\\Normal Text",
			expected: "Normal Text",
		},
		{
			name:     "multiple escape sequences",
			input:    "\x1b[1m\x1b[31mBold Red\x1b[0m Normal",
			expected: "Bold Red Normal",
		},
		{
			name:     "control characters with safe ones",
			input:    "Line1\nLine2\tTabbed\rReturn",
			expected: "Line1\nLine2\tTabbed\rReturn",
		},
		{
			name:     "dangerous control characters",
			input:    "Text\x00\x01\x02\x03More",
			expected: "TextMore",
		},
		{
			name:     "terminal injection attempt",
			input:    "Legit text\x1b[2J\x1b[H\x1b]0;Malicious\x07",
			expected: "Legit text",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "lone ESC at end",
			input:    "Text\x1b",
			expected: "Text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Sanitize(tt.input)
			if result != tt.expected {
				t.Errorf("Sanitize(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStripANSI(t *testing.T) {
	input := "\x1b[1;31mBold Red\x1b[0m"
	expected := "Bold Red"
	result := StripANSI(input)
	if result != expected {
		t.Errorf("StripANSI(%q) = %q, want %q", input, result, expected)
	}
}

func TestIsSafe(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "safe plain text",
			input:    "Hello, World!",
			expected: true,
		},
		{
			name:     "safe with newline",
			input:    "Line1\nLine2",
			expected: true,
		},
		{
			name:     "safe with tab",
			input:    "Col1\tCol2",
			expected: true,
		},
		{
			name:     "unsafe with ESC",
			input:    "Text\x1b[31m",
			expected: false,
		},
		{
			name:     "unsafe with C1 CSI (as UTF-8 rune)",
			input:    "Text" + string(rune(0x9b)) + "31m",
			expected: false,
		},
		{
			name:     "unsafe with null byte",
			input:    "Text\x00More",
			expected: false,
		},
		{
			name:     "unsafe with control char",
			input:    "Text\x01More",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSafe(tt.input)
			if result != tt.expected {
				t.Errorf("IsSafe(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeLongInput(t *testing.T) {
	// Test with a very long input to ensure no performance issues
	var sb strings.Builder
	for i := 0; i < 10000; i++ {
		sb.WriteString("\x1b[31mRed\x1b[0m ")
	}
	input := sb.String()
	result := Sanitize(input)
	
	// Should only contain "Red " repeated
	expected := strings.Repeat("Red ", 10000)
	if result != expected {
		t.Errorf("Sanitize long input failed")
	}
}

func BenchmarkSanitize(b *testing.B) {
	input := "\x1b[1;31mBold Red\x1b[0m Normal \x1b[32mGreen\x1b[0m"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Sanitize(input)
	}
}

func BenchmarkIsSafe(b *testing.B) {
	input := "This is a safe string with no escape sequences"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsSafe(input)
	}
}
