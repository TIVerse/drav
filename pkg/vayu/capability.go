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
func pathMatches(path, pattern string) bool {
	// Simple prefix matching for now
	// TODO: Implement glob pattern matching
	if len(pattern) > 0 && pattern[len(pattern)-1] == '*' {
		prefix := pattern[:len(pattern)-1]
		return len(path) >= len(prefix) && path[:len(prefix)] == prefix
	}
	return path == pattern
}

// domainMatches checks if a domain matches an allowed pattern.
func domainMatches(domain, pattern string) bool {
	// Simple wildcard matching
	if len(pattern) > 0 && pattern[0] == '*' {
		suffix := pattern[1:]
		return len(domain) >= len(suffix) && domain[len(domain)-len(suffix):] == suffix
	}
	return domain == pattern
}
