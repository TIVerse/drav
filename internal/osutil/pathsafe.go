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
// It normalizes and validates all allowed roots.
func NewPathValidator(allowedRoots []string) *PathValidator {
	normalized := make([]string, 0, len(allowedRoots))
	for _, root := range allowedRoots {
		// Clean and get absolute path
		clean := filepath.Clean(root)
		abs, err := filepath.Abs(clean)
		if err != nil {
			continue
		}
		// Evaluate symlinks to get real path
		real, err := filepath.EvalSymlinks(abs)
		if err != nil {
			// If symlink evaluation fails, use abs path
			real = abs
		}
		normalized = append(normalized, real)
	}
	return &PathValidator{
		allowedRoots: normalized,
	}
}

// Validate validates a path against allowed roots.
// It prevents directory traversal attacks by:
// 1. Cleaning and normalizing the path
// 2. Resolving symlinks
// 3. Checking if the resolved path is within allowed roots
func (pv *PathValidator) Validate(path string) error {
	if path == "" {
		return fmt.Errorf("empty path not allowed")
	}

	// Clean the path (removes .., ., etc.)
	cleanPath := filepath.Clean(path)

	// Get absolute path
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}

	// Resolve symlinks to get the real path
	// This prevents symlink-based traversal attacks
	realPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		// If file doesn't exist yet, use absolute path
		// But check parent directory
		dir := filepath.Dir(absPath)
		if dir != absPath {
			if err := pv.Validate(dir); err != nil {
				return fmt.Errorf("parent directory not allowed: %w", err)
			}
			return nil
		}
		realPath = absPath
	}

	// Check against allowed roots
	for _, root := range pv.allowedRoots {
		// Use filepath.Rel to check if path is under root
		rel, err := filepath.Rel(root, realPath)
		if err != nil {
			continue
		}

		// If relative path doesn't start with "..", it's within root
		if !strings.HasPrefix(rel, "..") && !strings.HasPrefix(rel, string(filepath.Separator)+"..") {
			return nil
		}
	}

	return fmt.Errorf("path not allowed: %s (resolved to %s)", path, realPath)
}

// IsWithinRoot checks if a path is within a root directory.
// It properly handles symlinks and path traversal attempts.
func IsWithinRoot(root, path string) bool {
	if root == "" || path == "" {
		return false
	}

	// Clean and normalize root
	cleanRoot := filepath.Clean(root)
	absRoot, err := filepath.Abs(cleanRoot)
	if err != nil {
		return false
	}
	realRoot, err := filepath.EvalSymlinks(absRoot)
	if err != nil {
		realRoot = absRoot
	}

	// Clean and normalize path
	cleanPath := filepath.Clean(path)
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return false
	}
	realPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		// If file doesn't exist, check parent
		dir := filepath.Dir(absPath)
		if dir != absPath {
			return IsWithinRoot(realRoot, dir)
		}
		realPath = absPath
	}

	// Use filepath.Rel to check relationship
	rel, err := filepath.Rel(realRoot, realPath)
	if err != nil {
		return false
	}

	// Path is within root if relative path doesn't escape upward
	return !strings.HasPrefix(rel, "..") && !strings.HasPrefix(rel, string(filepath.Separator)+"..")
}

// IsSafePath performs basic safety checks on a path string.
// It rejects paths with null bytes and other dangerous patterns.
func IsSafePath(path string) bool {
	// Reject paths with null bytes
	if strings.ContainsRune(path, 0) {
		return false
	}

	// Reject paths that are just "." or ".."
	if path == "." || path == ".." {
		return false
	}

	return true
}
