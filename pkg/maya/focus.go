package maya

import (
	"sync"
)

// FocusManager manages focus state across components.
type FocusManager struct {
	mu            sync.RWMutex
	focusedID     string
	focusableIDs  []string
	focusListeners map[string][]FocusListener
}

// FocusListener is called when focus state changes.
type FocusListener func(focused bool)

// NewFocusManager creates a new focus manager.
func NewFocusManager() *FocusManager {
	return &FocusManager{
		focusableIDs:   make([]string, 0),
		focusListeners: make(map[string][]FocusListener),
	}
}

// Register registers a component as focusable.
func (fm *FocusManager) Register(id string) {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	
	// Check if already registered
	for _, existingID := range fm.focusableIDs {
		if existingID == id {
			return
		}
	}
	
	fm.focusableIDs = append(fm.focusableIDs, id)
	
	// If this is the first component, focus it
	if len(fm.focusableIDs) == 1 {
		fm.focusedID = id
		fm.notifyListeners(id, true)
	}
}

// Unregister removes a component from focusable list.
func (fm *FocusManager) Unregister(id string) {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	
	// Remove from focusable list
	for i, existingID := range fm.focusableIDs {
		if existingID == id {
			fm.focusableIDs = append(fm.focusableIDs[:i], fm.focusableIDs[i+1:]...)
			break
		}
	}
	
	// If focused component was removed, focus next
	if fm.focusedID == id {
		fm.notifyListeners(id, false)
		fm.focusedID = ""
		if len(fm.focusableIDs) > 0 {
			fm.focusedID = fm.focusableIDs[0]
			fm.notifyListeners(fm.focusedID, true)
		}
	}
	
	// Remove listeners
	delete(fm.focusListeners, id)
}

// Focus sets focus to a specific component.
func (fm *FocusManager) Focus(id string) bool {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	
	// Check if component is focusable
	found := false
	for _, focusableID := range fm.focusableIDs {
		if focusableID == id {
			found = true
			break
		}
	}
	
	if !found {
		return false
	}
	
	// Unfocus previous
	if fm.focusedID != "" && fm.focusedID != id {
		fm.notifyListeners(fm.focusedID, false)
	}
	
	// Focus new
	fm.focusedID = id
	fm.notifyListeners(id, true)
	
	return true
}

// FocusNext focuses the next component in the list (Tab navigation).
func (fm *FocusManager) FocusNext() {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	
	if len(fm.focusableIDs) == 0 {
		return
	}
	
	// Find current index
	currentIdx := -1
	for i, id := range fm.focusableIDs {
		if id == fm.focusedID {
			currentIdx = i
			break
		}
	}
	
	// Calculate next index
	nextIdx := (currentIdx + 1) % len(fm.focusableIDs)
	
	// Unfocus current
	if fm.focusedID != "" {
		fm.notifyListeners(fm.focusedID, false)
	}
	
	// Focus next
	fm.focusedID = fm.focusableIDs[nextIdx]
	fm.notifyListeners(fm.focusedID, true)
}

// FocusPrevious focuses the previous component (Shift+Tab navigation).
func (fm *FocusManager) FocusPrevious() {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	
	if len(fm.focusableIDs) == 0 {
		return
	}
	
	// Find current index
	currentIdx := -1
	for i, id := range fm.focusableIDs {
		if id == fm.focusedID {
			currentIdx = i
			break
		}
	}
	
	// Calculate previous index
	prevIdx := currentIdx - 1
	if prevIdx < 0 {
		prevIdx = len(fm.focusableIDs) - 1
	}
	
	// Unfocus current
	if fm.focusedID != "" {
		fm.notifyListeners(fm.focusedID, false)
	}
	
	// Focus previous
	fm.focusedID = fm.focusableIDs[prevIdx]
	fm.notifyListeners(fm.focusedID, true)
}

// IsFocused returns true if the component with the given ID is focused.
func (fm *FocusManager) IsFocused(id string) bool {
	fm.mu.RLock()
	defer fm.mu.RUnlock()
	return fm.focusedID == id
}

// GetFocused returns the ID of the currently focused component.
func (fm *FocusManager) GetFocused() string {
	fm.mu.RLock()
	defer fm.mu.RUnlock()
	return fm.focusedID
}

// OnFocus registers a listener for focus changes on a component.
func (fm *FocusManager) OnFocus(id string, listener FocusListener) func() {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	
	if fm.focusListeners[id] == nil {
		fm.focusListeners[id] = make([]FocusListener, 0)
	}
	
	fm.focusListeners[id] = append(fm.focusListeners[id], listener)
	
	// Return unsubscribe function
	return func() {
		fm.mu.Lock()
		defer fm.mu.Unlock()
		
		listeners := fm.focusListeners[id]
		for i, l := range listeners {
			// Compare function pointers (not perfect but works for most cases)
			if &l == &listener {
				fm.focusListeners[id] = append(listeners[:i], listeners[i+1:]...)
				break
			}
		}
	}
}

// notifyListeners notifies all listeners for a component (must be called with lock held).
func (fm *FocusManager) notifyListeners(id string, focused bool) {
	listeners := fm.focusListeners[id]
	for _, listener := range listeners {
		// Call listener outside of lock to avoid deadlock
		go listener(focused)
	}
}

// Clear clears all focus state.
func (fm *FocusManager) Clear() {
	fm.mu.Lock()
	defer fm.mu.Unlock()
	
	// Notify all focused components they're losing focus
	if fm.focusedID != "" {
		fm.notifyListeners(fm.focusedID, false)
	}
	
	fm.focusedID = ""
	fm.focusableIDs = make([]string, 0)
	fm.focusListeners = make(map[string][]FocusListener)
}

// GetFocusableCount returns the number of focusable components.
func (fm *FocusManager) GetFocusableCount() int {
	fm.mu.RLock()
	defer fm.mu.RUnlock()
	return len(fm.focusableIDs)
}
