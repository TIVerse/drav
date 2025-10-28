package fuzz

import (
	"testing"

	"github.com/TIVerse/drav/pkg/maya"
)

func FuzzDiff(f *testing.F) {
	// Seed corpus
	f.Add(10, 10)
	f.Add(80, 24)
	f.Add(120, 40)

	f.Fuzz(func(t *testing.T, width, height int) {
		// Clamp to reasonable values
		if width < 1 || width > 1000 {
			return
		}
		if height < 1 || height > 1000 {
			return
		}

		// Create buffers
		buf1 := maya.NewBuffer(width, height)
		buf2 := maya.NewBuffer(width, height)

		// Write some data
		buf2.SetRune(0, 0, 'A')

		// Diff should not panic
		_ = maya.Diff(buf1, buf2)
	})
}
