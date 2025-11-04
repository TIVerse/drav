package prana

import "sync"

// RenderCallback is called when an observable changes to trigger a re-render.
type RenderCallback func()

var (
	globalRenderCallback RenderCallback
	globalCallbackMu     sync.RWMutex
)

// SetGlobalRenderCallback sets the global render callback.
// This should be called by the App during initialization.
func SetGlobalRenderCallback(callback RenderCallback) {
	globalCallbackMu.Lock()
	defer globalCallbackMu.Unlock()
	globalRenderCallback = callback
}

// GetGlobalRenderCallback returns the global render callback.
func GetGlobalRenderCallback() RenderCallback {
	globalCallbackMu.RLock()
	defer globalCallbackMu.RUnlock()
	return globalRenderCallback
}

// requestRender triggers a re-render if a callback is set.
func requestRender() {
	callback := GetGlobalRenderCallback()
	if callback != nil {
		callback()
	}
}
