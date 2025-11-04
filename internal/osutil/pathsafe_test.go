package osutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPathValidator_Validate(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	allowedDir := filepath.Join(tmpDir, "allowed")
	if err := os.MkdirAll(allowedDir, 0755); err != nil {
		t.Fatal(err)
	}

	validator := NewPathValidator([]string{allowedDir})

	tests := []struct {
		name      string
		path      string
		shouldErr bool
	}{
		{
			name:      "allowed path",
			path:      filepath.Join(allowedDir, "file.txt"),
			shouldErr: false,
		},
		{
			name:      "allowed subdirectory",
			path:      filepath.Join(allowedDir, "subdir", "file.txt"),
			shouldErr: false,
		},
		{
			name:      "path traversal with ..",
			path:      filepath.Join(allowedDir, "..", "forbidden.txt"),
			shouldErr: true,
		},
		{
			name:      "path outside allowed root",
			path:      filepath.Join(tmpDir, "forbidden", "file.txt"),
			shouldErr: true,
		},
		{
			name:      "empty path",
			path:      "",
			shouldErr: true,
		},
		{
			name:      "absolute path outside root",
			path:      "/etc/passwd",
			shouldErr: true,
		},
		{
			name:      "allowed root itself",
			path:      allowedDir,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.path)
			if (err != nil) != tt.shouldErr {
				t.Errorf("Validate(%q) error = %v, shouldErr = %v", tt.path, err, tt.shouldErr)
			}
		})
	}
}

func TestPathValidator_ValidateWithSymlinks(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()
	allowedDir := filepath.Join(tmpDir, "allowed")
	forbiddenDir := filepath.Join(tmpDir, "forbidden")
	
	if err := os.MkdirAll(allowedDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(forbiddenDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Create a symlink from allowed to forbidden
	symlinkPath := filepath.Join(allowedDir, "link_to_forbidden")
	if err := os.Symlink(forbiddenDir, symlinkPath); err != nil {
		t.Skip("Cannot create symlinks, skipping test")
	}

	validator := NewPathValidator([]string{allowedDir})

	// Attempt to access forbidden dir through symlink
	targetPath := filepath.Join(symlinkPath, "secret.txt")
	err := validator.Validate(targetPath)
	
	// Should reject because symlink resolves outside allowed root
	if err == nil {
		t.Error("Expected error for symlink traversal, got nil")
	}
}

func TestIsWithinRoot(t *testing.T) {
	tmpDir := t.TempDir()
	root := filepath.Join(tmpDir, "root")
	if err := os.MkdirAll(root, 0755); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		root     string
		path     string
		expected bool
	}{
		{
			name:     "path within root",
			root:     root,
			path:     filepath.Join(root, "file.txt"),
			expected: true,
		},
		{
			name:     "path within subdirectory",
			root:     root,
			path:     filepath.Join(root, "sub", "file.txt"),
			expected: true,
		},
		{
			name:     "path outside root with ..",
			root:     root,
			path:     filepath.Join(root, "..", "outside.txt"),
			expected: false,
		},
		{
			name:     "completely different path",
			root:     root,
			path:     "/etc/passwd",
			expected: false,
		},
		{
			name:     "empty root",
			root:     "",
			path:     "/some/path",
			expected: false,
		},
		{
			name:     "empty path",
			root:     root,
			path:     "",
			expected: false,
		},
		{
			name:     "root itself",
			root:     root,
			path:     root,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsWithinRoot(tt.root, tt.path)
			if result != tt.expected {
				t.Errorf("IsWithinRoot(%q, %q) = %v, want %v", tt.root, tt.path, result, tt.expected)
			}
		})
	}
}

func TestIsSafePath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "safe path",
			path:     "/home/user/file.txt",
			expected: true,
		},
		{
			name:     "path with null byte",
			path:     "/home/user/file\x00.txt",
			expected: false,
		},
		{
			name:     "just dot",
			path:     ".",
			expected: false,
		},
		{
			name:     "just double dot",
			path:     "..",
			expected: false,
		},
		{
			name:     "relative path with ..",
			path:     "../file.txt",
			expected: true, // Not rejected by IsSafePath, only by validator
		},
		{
			name:     "empty path",
			path:     "",
			expected: true, // Empty is technically safe, validator rejects it
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsSafePath(tt.path)
			if result != tt.expected {
				t.Errorf("IsSafePath(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestPathValidator_MultipleRoots(t *testing.T) {
	tmpDir := t.TempDir()
	root1 := filepath.Join(tmpDir, "root1")
	root2 := filepath.Join(tmpDir, "root2")
	
	if err := os.MkdirAll(root1, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(root2, 0755); err != nil {
		t.Fatal(err)
	}

	validator := NewPathValidator([]string{root1, root2})

	// Should allow paths in either root
	path1 := filepath.Join(root1, "file.txt")
	if err := validator.Validate(path1); err != nil {
		t.Errorf("Expected path in root1 to be allowed: %v", err)
	}

	path2 := filepath.Join(root2, "file.txt")
	if err := validator.Validate(path2); err != nil {
		t.Errorf("Expected path in root2 to be allowed: %v", err)
	}

	// Should reject paths outside both roots
	outside := filepath.Join(tmpDir, "outside", "file.txt")
	if err := validator.Validate(outside); err == nil {
		t.Error("Expected path outside roots to be rejected")
	}
}

func BenchmarkPathValidatorValidate(b *testing.B) {
	tmpDir := b.TempDir()
	root := filepath.Join(tmpDir, "root")
	if err := os.MkdirAll(root, 0755); err != nil {
		b.Fatal(err)
	}

	validator := NewPathValidator([]string{root})
	path := filepath.Join(root, "deep", "nested", "path", "to", "file.txt")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validator.Validate(path)
	}
}

func BenchmarkIsWithinRoot(b *testing.B) {
	tmpDir := b.TempDir()
	root := filepath.Join(tmpDir, "root")
	path := filepath.Join(root, "sub", "file.txt")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsWithinRoot(root, path)
	}
}
