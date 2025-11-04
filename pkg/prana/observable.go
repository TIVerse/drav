package prana

import (
	"sync"
	"sync/atomic"
)

// Observable represents an observable value that notifies watchers on changes.
type Observable[T any] struct {
	mu       sync.RWMutex
	value    T
	watchers map[uint64]Watcher[T]
	nextID   atomic.Uint64
}

// Watcher is a function called when the observable value changes.
type Watcher[T any] func(oldValue, newValue T)

// NewObservable creates a new observable with an initial value.
func NewObservable[T any](initial T) *Observable[T] {
	return &Observable[T]{
		value:    initial,
		watchers: make(map[uint64]Watcher[T]),
	}
}

// Get returns the current value (read-only).
func (o *Observable[T]) Get() T {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.value
}

// Set updates the value and notifies watchers.
func (o *Observable[T]) Set(newValue T) {
	o.mu.Lock()
	oldValue := o.value
	o.value = newValue
	watchers := make([]Watcher[T], 0, len(o.watchers))
	for _, w := range o.watchers {
		watchers = append(watchers, w)
	}
	o.mu.Unlock()

	// Notify watchers outside the lock
	// Recover from panics in watchers to prevent crashes
	for _, watcher := range watchers {
		func(w Watcher[T]) {
			defer func() {
				if r := recover(); r != nil {
					// Log panic but continue with other watchers
					// In production, this should use proper logging
				}
			}()
			w(oldValue, newValue)
		}(watcher)
	}

	// Request a re-render after notifying watchers
	requestRender()
}

// Update applies a function to update the value.
func (o *Observable[T]) Update(fn func(T) T) {
	o.mu.Lock()
	oldValue := o.value
	o.value = fn(o.value)
	newValue := o.value
	watchers := make([]Watcher[T], 0, len(o.watchers))
	for _, w := range o.watchers {
		watchers = append(watchers, w)
	}
	o.mu.Unlock()

	// Notify watchers outside the lock
	// Recover from panics in watchers to prevent crashes
	for _, watcher := range watchers {
		func(w Watcher[T]) {
			defer func() {
				if r := recover(); r != nil {
					// Log panic but continue with other watchers
					// In production, this should use proper logging
				}
			}()
			w(oldValue, newValue)
		}(watcher)
	}

	// Request a re-render after notifying watchers
	requestRender()
}

// Watch registers a watcher and returns an unwatch function.
func (o *Observable[T]) Watch(watcher Watcher[T]) func() {
	id := o.nextID.Add(1)
	
	o.mu.Lock()
	o.watchers[id] = watcher
	o.mu.Unlock()

	// Return unwatch function
	return func() {
		o.Unwatch(id)
	}
}

// Unwatch removes a watcher by ID.
func (o *Observable[T]) Unwatch(id uint64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	delete(o.watchers, id)
}

// WatcherCount returns the number of active watchers.
func (o *Observable[T]) WatcherCount() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.watchers)
}

// ClearWatchers removes all watchers.
func (o *Observable[T]) ClearWatchers() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.watchers = make(map[uint64]Watcher[T])
}

// Derive creates a new observable derived from this one.
func (o *Observable[T]) Derive(fn func(T) T) *Observable[T] {
	derived := NewObservable(fn(o.Get()))
	o.Watch(func(_, newValue T) {
		derived.Set(fn(newValue))
	})
	return derived
}
