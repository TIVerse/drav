package vayu

import (
	"testing"
)

func TestCapabilityChecker_CanReadPath(t *testing.T) {
	caps := Capabilities{
		Filesystem: FSCapability{
			Read: []string{"/data/*", "/config/app.conf"},
		},
	}
	checker := NewCapabilityChecker(caps)

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "exact match",
			path:     "/config/app.conf",
			expected: true,
		},
		{
			name:     "wildcard match - file",
			path:     "/data/file.txt",
			expected: true,
		},
		{
			name:     "wildcard match - nested",
			path:     "/data/subdir/file.txt",
			expected: true,
		},
		{
			name:     "wildcard no match - parent dir",
			path:     "/data",
			expected: false,
		},
		{
			name:     "no match - different path",
			path:     "/etc/passwd",
			expected: false,
		},
		{
			name:     "empty path",
			path:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checker.CanReadPath(tt.path)
			if result != tt.expected {
				t.Errorf("CanReadPath(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestCapabilityChecker_CanWritePath(t *testing.T) {
	caps := Capabilities{
		Filesystem: FSCapability{
			Write: []string{"/tmp/*", "/var/app/data.db"},
		},
	}
	checker := NewCapabilityChecker(caps)

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "exact match",
			path:     "/var/app/data.db",
			expected: true,
		},
		{
			name:     "wildcard match",
			path:     "/tmp/tempfile.txt",
			expected: true,
		},
		{
			name:     "wildcard no match - parent",
			path:     "/tmp",
			expected: false,
		},
		{
			name:     "no match",
			path:     "/etc/shadow",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checker.CanWritePath(tt.path)
			if result != tt.expected {
				t.Errorf("CanWritePath(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestCapabilityChecker_CanAccessDomain(t *testing.T) {
	caps := Capabilities{
		Network: NetworkCapability{
			AllowedDomains: []string{"api.example.com", "*.github.com"},
		},
	}
	checker := NewCapabilityChecker(caps)

	tests := []struct {
		name     string
		domain   string
		expected bool
	}{
		{
			name:     "exact match",
			domain:   "api.example.com",
			expected: true,
		},
		{
			name:     "wildcard match - subdomain",
			domain:   "raw.github.com",
			expected: true,
		},
		{
			name:     "wildcard match - deep subdomain",
			domain:   "api.raw.github.com",
			expected: true,
		},
		{
			name:     "no match - similar domain",
			domain:   "example.com",
			expected: false,
		},
		{
			name:     "no match - different domain",
			domain:   "malicious.com",
			expected: false,
		},
		{
			name:     "empty domain",
			domain:   "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checker.CanAccessDomain(tt.domain)
			if result != tt.expected {
				t.Errorf("CanAccessDomain(%q) = %v, want %v", tt.domain, result, tt.expected)
			}
		})
	}
}

func TestCapabilityChecker_CanAccessPort(t *testing.T) {
	caps := Capabilities{
		Network: NetworkCapability{
			AllowedPorts: []int{80, 443, 8080},
		},
	}
	checker := NewCapabilityChecker(caps)

	tests := []struct {
		name     string
		port     int
		expected bool
	}{
		{
			name:     "allowed port 80",
			port:     80,
			expected: true,
		},
		{
			name:     "allowed port 443",
			port:     443,
			expected: true,
		},
		{
			name:     "allowed port 8080",
			port:     8080,
			expected: true,
		},
		{
			name:     "disallowed port 22",
			port:     22,
			expected: false,
		},
		{
			name:     "disallowed port 3306",
			port:     3306,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checker.CanAccessPort(tt.port)
			if result != tt.expected {
				t.Errorf("CanAccessPort(%d) = %v, want %v", tt.port, result, tt.expected)
			}
		})
	}
}

func TestCapabilityChecker_EmptyPorts(t *testing.T) {
	caps := Capabilities{
		Network: NetworkCapability{
			AllowedPorts: []int{}, // Empty list
		},
	}
	checker := NewCapabilityChecker(caps)

	// With empty list, all ports should be denied
	if checker.CanAccessPort(80) {
		t.Error("Expected port 80 to be denied with empty allowed ports")
	}
}

func TestCapabilityChecker_EmptyDomains(t *testing.T) {
	caps := Capabilities{
		Network: NetworkCapability{
			AllowedDomains: []string{}, // Empty list
		},
	}
	checker := NewCapabilityChecker(caps)

	// With empty list, all domains should be denied
	if checker.CanAccessDomain("example.com") {
		t.Error("Expected domain to be denied with empty allowed domains")
	}
}

func TestCapabilityChecker_CanReadStore(t *testing.T) {
	caps := Capabilities{
		State: StateCapability{
			ReadStores: []string{"user", "settings"},
		},
	}
	checker := NewCapabilityChecker(caps)

	tests := []struct {
		name     string
		store    string
		expected bool
	}{
		{
			name:     "allowed store - user",
			store:    "user",
			expected: true,
		},
		{
			name:     "allowed store - settings",
			store:    "settings",
			expected: true,
		},
		{
			name:     "disallowed store",
			store:    "secrets",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checker.CanReadStore(tt.store)
			if result != tt.expected {
				t.Errorf("CanReadStore(%q) = %v, want %v", tt.store, result, tt.expected)
			}
		})
	}
}

func TestCapabilityChecker_CanWriteStore(t *testing.T) {
	caps := Capabilities{
		State: StateCapability{
			WriteStores: []string{"cache", "temp"},
		},
	}
	checker := NewCapabilityChecker(caps)

	tests := []struct {
		name     string
		store    string
		expected bool
	}{
		{
			name:     "allowed store - cache",
			store:    "cache",
			expected: true,
		},
		{
			name:     "allowed store - temp",
			store:    "temp",
			expected: true,
		},
		{
			name:     "disallowed store",
			store:    "user",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checker.CanWriteStore(tt.store)
			if result != tt.expected {
				t.Errorf("CanWriteStore(%q) = %v, want %v", tt.store, result, tt.expected)
			}
		})
	}
}

func TestPathMatches(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		pattern  string
		expected bool
	}{
		{
			name:     "exact match",
			path:     "/data/file.txt",
			pattern:  "/data/file.txt",
			expected: true,
		},
		{
			name:     "wildcard prefix match",
			path:     "/data/subdir/file.txt",
			pattern:  "/data/*",
			expected: true,
		},
		{
			name:     "wildcard no match - parent dir itself",
			path:     "/data",
			pattern:  "/data/*",
			expected: false,
		},
		{
			name:     "no match",
			path:     "/etc/passwd",
			pattern:  "/data/*",
			expected: false,
		},
		{
			name:     "empty path",
			path:     "",
			pattern:  "/data/*",
			expected: false,
		},
		{
			name:     "empty pattern",
			path:     "/data/file.txt",
			pattern:  "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pathMatches(tt.path, tt.pattern)
			if result != tt.expected {
				t.Errorf("pathMatches(%q, %q) = %v, want %v", tt.path, tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestDomainMatches(t *testing.T) {
	tests := []struct {
		name     string
		domain   string
		pattern  string
		expected bool
	}{
		{
			name:     "exact match",
			domain:   "example.com",
			pattern:  "example.com",
			expected: true,
		},
		{
			name:     "wildcard subdomain match",
			domain:   "api.example.com",
			pattern:  "*.example.com",
			expected: true,
		},
		{
			name:     "wildcard deep subdomain match",
			domain:   "v1.api.example.com",
			pattern:  "*.example.com",
			expected: true,
		},
		{
			name:     "wildcard no match - base domain",
			domain:   "example.com",
			pattern:  "*.example.com",
			expected: false,
		},
		{
			name:     "no match - different domain",
			domain:   "example.org",
			pattern:  "*.example.com",
			expected: false,
		},
		{
			name:     "empty domain",
			domain:   "",
			pattern:  "*.example.com",
			expected: false,
		},
		{
			name:     "empty pattern",
			domain:   "example.com",
			pattern:  "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := domainMatches(tt.domain, tt.pattern)
			if result != tt.expected {
				t.Errorf("domainMatches(%q, %q) = %v, want %v", tt.domain, tt.pattern, result, tt.expected)
			}
		})
	}
}

func BenchmarkCanReadPath(b *testing.B) {
	caps := Capabilities{
		Filesystem: FSCapability{
			Read: []string{"/data/*", "/config/*", "/tmp/*"},
		},
	}
	checker := NewCapabilityChecker(caps)
	path := "/data/subdir/file.txt"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		checker.CanReadPath(path)
	}
}

func BenchmarkCanAccessDomain(b *testing.B) {
	caps := Capabilities{
		Network: NetworkCapability{
			AllowedDomains: []string{"*.example.com", "*.github.com", "api.service.io"},
		},
	}
	checker := NewCapabilityChecker(caps)
	domain := "api.example.com"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		checker.CanAccessDomain(domain)
	}
}
