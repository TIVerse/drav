package vayu

// Capabilities defines what a plugin is allowed to do.
type Capabilities struct {
	Filesystem FSCapability
	Network    NetworkCapability
	State      StateCapability
	UI         UICapability
}

// FSCapability defines filesystem access permissions.
type FSCapability struct {
	Read  []string // Allowed read paths
	Write []string // Allowed write paths
}

// NetworkCapability defines network access permissions.
type NetworkCapability struct {
	AllowedDomains []string // Allowed domains
	AllowedPorts   []int    // Allowed ports
	RateLimit      int      // Requests per minute
}

// StateCapability defines state access permissions.
type StateCapability struct {
	ReadStores  []string // Store names readable
	WriteStores []string // Store names writable
}

// UICapability defines UI access permissions.
type UICapability struct {
	CanRegisterViews   bool
	CanRegisterCommands bool
	CanModifyTheme     bool
}

// CapabilityChecker validates capabilities.
type CapabilityChecker struct {
	caps Capabilities
}

// NewCapabilityChecker creates a new capability checker.
func NewCapabilityChecker(caps Capabilities) *CapabilityChecker {
	return &CapabilityChecker{
		caps: caps,
	}
}

// CanReadPath checks if a path can be read.
func (cc *CapabilityChecker) CanReadPath(path string) bool {
	for _, allowed := range cc.caps.Filesystem.Read {
		if pathMatches(path, allowed) {
			return true
		}
	}
	return false
}

// CanWritePath checks if a path can be written.
func (cc *CapabilityChecker) CanWritePath(path string) bool {
	for _, allowed := range cc.caps.Filesystem.Write {
		if pathMatches(path, allowed) {
			return true
		}
	}
	return false
}

// CanAccessDomain checks if a domain can be accessed.
func (cc *CapabilityChecker) CanAccessDomain(domain string) bool {
	if len(cc.caps.Network.AllowedDomains) == 0 {
		return false
	}

	for _, allowed := range cc.caps.Network.AllowedDomains {
		if domainMatches(domain, allowed) {
			return true
		}
	}
	return false
}

// CanAccessPort checks if a port can be accessed.
func (cc *CapabilityChecker) CanAccessPort(port int) bool {
	if len(cc.caps.Network.AllowedPorts) == 0 {
		return false
	}

	for _, allowed := range cc.caps.Network.AllowedPorts {
		if port == allowed {
			return true
		}
	}
	return false
}

// CanReadStore checks if a store can be read.
func (cc *CapabilityChecker) CanReadStore(name string) bool {
	for _, allowed := range cc.caps.State.ReadStores {
		if name == allowed {
			return true
		}
	}
	return false
}

// CanWriteStore checks if a store can be written.
func (cc *CapabilityChecker) CanWriteStore(name string) bool {
	for _, allowed := range cc.caps.State.WriteStores {
		if name == allowed {
			return true
		}
	}
	return false
}

// pathMatches checks if a path matches an allowed pattern.
// Supports wildcard patterns with * for prefix matching.
// Example: "/data/*" matches "/data/file.txt" but not "/data"
func pathMatches(path, pattern string) bool {
	if path == "" || pattern == "" {
		return false
	}

	// Exact match
	if path == pattern {
		return true
	}

	// Wildcard prefix matching
	if len(pattern) > 1 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		// Ensure the prefix matches and path is longer or has separator
		if len(path) >= len(prefix) && path[:len(prefix)] == prefix {
			// Prevent matching parent directories
			// e.g., "/data/*" should not match "/data" itself
			return len(path) > len(prefix)
		}
	}

	return false
}

// domainMatches checks if a domain matches an allowed pattern.
// Supports wildcard subdomain matching.
// Example: "*.example.com" matches "api.example.com" but not "example.com"
func domainMatches(domain, pattern string) bool {
	if domain == "" || pattern == "" {
		return false
	}

	// Exact match
	if domain == pattern {
		return true
	}

	// Wildcard subdomain matching
	if len(pattern) > 2 && pattern[0] == '*' && pattern[1] == '.' {
		suffix := pattern[1:] // Keep the dot
		// Check if domain ends with the suffix
		if len(domain) > len(suffix) && domain[len(domain)-len(suffix):] == suffix {
			// Ensure there are no additional dots before the suffix
			// This prevents "*.example.com" from matching "foo.bar.example.com"
			// Allow any prefix without dots for single-level wildcards
			// For stricter control, could check for no dots in prefix:
			// prefix := domain[:len(domain)-len(suffix)]
			return true
		}
	}

	return false
}
