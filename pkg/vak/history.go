package vak

import (
	"sync"
	"time"
)

// History manages command execution history.
type History struct {
	mu      sync.RWMutex
	entries []HistoryEntry
	maxSize int
}

// HistoryEntry represents a command in history.
type HistoryEntry struct {
	Input     string
	Command   string
	Success   bool
	Timestamp time.Time
}

// NewHistory creates a new history with a maximum size.
func NewHistory(maxSize int) *History {
	return &History{
		entries: make([]HistoryEntry, 0, maxSize),
		maxSize: maxSize,
	}
}

// Add adds an entry to history.
func (h *History) Add(entry HistoryEntry) {
	h.mu.Lock()
	defer h.mu.Unlock()

	entry.Timestamp = time.Now()

	if len(h.entries) >= h.maxSize {
		// Remove oldest entry
		h.entries = h.entries[1:]
	}

	h.entries = append(h.entries, entry)
}

// Get retrieves a history entry by index (0 is oldest).
func (h *History) Get(index int) (HistoryEntry, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if index < 0 || index >= len(h.entries) {
		return HistoryEntry{}, false
	}

	return h.entries[index], true
}

// Last retrieves the most recent entry.
func (h *History) Last() (HistoryEntry, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.entries) == 0 {
		return HistoryEntry{}, false
	}

	return h.entries[len(h.entries)-1], true
}

// All returns all history entries (oldest to newest).
func (h *History) All() []HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()

	entries := make([]HistoryEntry, len(h.entries))
	copy(entries, h.entries)
	return entries
}

// Clear clears all history.
func (h *History) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = make([]HistoryEntry, 0, h.maxSize)
}

// Size returns the number of entries.
func (h *History) Size() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.entries)
}

// Search searches history for entries matching a query.
func (h *History) Search(query string) []HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()

	matches := make([]HistoryEntry, 0)
	for _, entry := range h.entries {
		if contains(entry.Input, query) || contains(entry.Command, query) {
			matches = append(matches, entry)
		}
	}
	return matches
}

// contains checks if s contains substr (case-insensitive).
func contains(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if toLower(s[i+j]) != toLower(substr[j]) {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
