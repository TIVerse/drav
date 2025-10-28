package fuzz

import (
	"testing"

	"github.com/TIVerse/drav/pkg/vak"
)

func FuzzCommandParser(f *testing.F) {
	// Seed corpus
	f.Add("hello world")
	f.Add("cmd --flag value")
	f.Add("test -f \"quoted string\"")

	f.Fuzz(func(t *testing.T, input string) {
		parser := vak.NewParser()
		// Parser should not panic on any input
		_, _ = parser.Parse(input)
	})
}
