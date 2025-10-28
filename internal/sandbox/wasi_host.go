package sandbox

// WASIHost provides WASI host functions for WASM plugins.
// This is a stub implementation for future WASM support.
type WASIHost struct {
	allowedPaths []string
}

// NewWASIHost creates a new WASI host.
func NewWASIHost(allowedPaths []string) *WASIHost {
	return &WASIHost{
		allowedPaths: allowedPaths,
	}
}

// CanAccessPath checks if a path can be accessed.
func (wh *WASIHost) CanAccessPath(path string) bool {
	for _, allowed := range wh.allowedPaths {
		if path == allowed {
			return true
		}
	}
	return false
}
