package prana

import (
	"sync"
)

// Computed represents a computed value with dependency tracking.
type Computed[T any] struct {
	mu           sync.RWMutex
	compute      func() T
	value        T
	dirty        bool
	dependencies []Dependency
	watchers     map[uint64]Watcher[T]
	nextID       uint64
}

// Dependency represents a dependency for computed values.
type Dependency interface {
	OnChange(callback func())
}

// NewComputed creates a new computed value.
func NewComputed[T any](compute func() T) *Computed[T] {
	c := &Computed[T]{
		compute:  compute,
		dirty:    true,
		watchers: make(map[uint64]Watcher[T]),
	}
	return c
}

// Get returns the computed value, recomputing if dirty.
func (c *Computed[T]) Get() T {
	c.mu.RLock()
	if !c.dirty {
		value := c.value
		c.mu.RUnlock()
		return value
	}
	c.mu.RUnlock()

	// Need to recompute
	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if !c.dirty {
		return c.value
	}

	oldValue := c.value
	c.value = c.compute()
	c.dirty = false
	newValue := c.value

	// Collect watchers
	watchers := make([]Watcher[T], 0, len(c.watchers))
	for _, w := range c.watchers {
		watchers = append(watchers, w)
	}

	// Notify watchers outside the main lock
	go func() {
		for _, watcher := range watchers {
			watcher(oldValue, newValue)
		}
	}()

	return c.value
}

// Invalidate marks the computed value as dirty.
func (c *Computed[T]) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dirty = true
}

// Watch registers a watcher for changes.
func (c *Computed[T]) Watch(watcher Watcher[T]) func() {
	c.mu.Lock()
	id := c.nextID
	c.nextID++
	c.watchers[id] = watcher
	c.mu.Unlock()

	return func() {
		c.mu.Lock()
		defer c.mu.Unlock()
		delete(c.watchers, id)
	}
}

// AddDependency adds a dependency that triggers recomputation.
func (c *Computed[T]) AddDependency(dep Dependency) {
	c.mu.Lock()
	c.dependencies = append(c.dependencies, dep)
	c.mu.Unlock()

	dep.OnChange(func() {
		c.Invalidate()
	})
}

// ComputedFromObservable creates a computed value from an observable.
func ComputedFromObservable[T, U any](obs *Observable[T], transform func(T) U) *Computed[U] {
	comp := NewComputed(func() U {
		return transform(obs.Get())
	})

	obs.Watch(func(_, _ T) {
		comp.Invalidate()
	})

	return comp
}
