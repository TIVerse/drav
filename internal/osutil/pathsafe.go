package osutil

import (
	"fmt"
	"path/filepath"
	"strings"
)

// PathValidator validates paths against security constraints.
type PathValidator struct {
	allowedRoots []string
}

// NewPathValidator creates a new path validator.
func NewPathValidator(allowedRoots []string) *PathValidator {
	return &PathValidator{
		allowedRoots: allowedRoots,
	}
}

// Validate validates a path against allowed roots.
func (pv *PathValidator) Validate(path string) error {
	cleanPath := filepath.Clean(path)
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Check against allowed roots
	for _, root := range pv.allowedRoots {
		absRoot, err := filepath.Abs(root)
		if err != nil {
			continue
		}

		if strings.HasPrefix(absPath, absRoot) {
			return nil
		}
	}

	return fmt.Errorf("path not allowed: %s", path)
}

// IsWithinRoot checks if a path is within a root directory.
func IsWithinRoot(root, path string) bool {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return false
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	return strings.HasPrefix(absPath, absRoot)
}
