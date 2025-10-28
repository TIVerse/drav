package osutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SafePath validates and cleans a file path to prevent directory traversal.
func SafePath(base, path string) (string, error) {
	// Clean the path
	cleanPath := filepath.Clean(path)

	// Resolve to absolute path
	absBase, err := filepath.Abs(base)
	if err != nil {
		return "", fmt.Errorf("invalid base path: %w", err)
	}

	absPath := filepath.Join(absBase, cleanPath)

	// Check if the path is within the base directory
	if !strings.HasPrefix(absPath, absBase) {
		return "", fmt.Errorf("path traversal detected: %s", path)
	}

	return absPath, nil
}

// FileExists checks if a file exists.
func FileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if a directory exists.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// EnsureDir ensures a directory exists, creating it if necessary.
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// IsReadable checks if a file is readable.
func IsReadable(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

// IsWritable checks if a file is writable.
func IsWritable(path string) bool {
	// Try to open for writing
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return false
	}
	file.Close()
	return true
}
