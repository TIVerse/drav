package prana

import (
	"context"
	"fmt"
	"sync"
)

// Store represents a state store with actions and effects.
type Store[S any] struct {
	mu        sync.RWMutex
	state     S
	reducers  map[string]Reducer[S]
	effects   map[string]Effect[S]
	watchers  map[uint64]StoreWatcher[S]
	nextID    uint64
	history   []S
	maxHistory int
}

// Reducer is a pure function that updates state.
type Reducer[S any] func(state S, payload any) S

// Effect is a side-effect function.
type Effect[S any] func(ctx context.Context, state S, payload any) error

// StoreWatcher is called when store state changes.
type StoreWatcher[S any] func(oldState, newState S)

// Action represents a dispatched action.
type Action struct {
	Type    string
	Payload any
}

// NewStore creates a new store with initial state.
func NewStore[S any](initial S) *Store[S] {
	return &Store[S]{
		state:      initial,
		reducers:   make(map[string]Reducer[S]),
		effects:    make(map[string]Effect[S]),
		watchers:   make(map[uint64]StoreWatcher[S]),
		history:    make([]S, 0),
		maxHistory: 100,
	}
}

// GetState returns the current state (read-only).
func (s *Store[S]) GetState() S {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

// SetState directly sets the state (use with caution).
func (s *Store[S]) SetState(newState S) {
	s.mu.Lock()
	oldState := s.state
	s.state = newState
	s.addToHistory(oldState)
	watchers := s.collectWatchers()
	s.mu.Unlock()

	s.notifyWatchers(watchers, oldState, newState)
}

// RegisterReducer registers a reducer for an action type.
func (s *Store[S]) RegisterReducer(actionType string, reducer Reducer[S]) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.reducers[actionType] = reducer
}

// RegisterEffect registers an effect for an action type.
func (s *Store[S]) RegisterEffect(actionType string, effect Effect[S]) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.effects[actionType] = effect
}

// Dispatch dispatches an action to update state.
func (s *Store[S]) Dispatch(ctx context.Context, action Action) error {
	// Apply reducer
	s.mu.Lock()
	reducer, hasReducer := s.reducers[action.Type]
	if hasReducer {
		oldState := s.state
		s.state = reducer(s.state, action.Payload)
		newState := s.state
		s.addToHistory(oldState)
		watchers := s.collectWatchers()
		s.mu.Unlock()
		s.notifyWatchers(watchers, oldState, newState)
	} else {
		s.mu.Unlock()
	}

	// Run effect
	s.mu.RLock()
	effect, hasEffect := s.effects[action.Type]
	currentState := s.state
	s.mu.RUnlock()

	if hasEffect {
		return effect(ctx, currentState, action.Payload)
	}

	if !hasReducer && !hasEffect {
		return fmt.Errorf("no reducer or effect registered for action type: %s", action.Type)
	}

	return nil
}

// Watch registers a state watcher.
func (s *Store[S]) Watch(watcher StoreWatcher[S]) func() {
	s.mu.Lock()
	id := s.nextID
	s.nextID++
	s.watchers[id] = watcher
	s.mu.Unlock()

	return func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.watchers, id)
	}
}

// Select extracts a slice of state.
func (s *Store[S]) Select(selector func(S) any) any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return selector(s.state)
}

// History returns the state history.
func (s *Store[S]) History() []S {
	s.mu.RLock()
	defer s.mu.RUnlock()
	historyCopy := make([]S, len(s.history))
	copy(historyCopy, s.history)
	return historyCopy
}

// addToHistory adds state to history.
func (s *Store[S]) addToHistory(state S) {
	if len(s.history) >= s.maxHistory {
		s.history = s.history[1:]
	}
	s.history = append(s.history, state)
}

// collectWatchers collects all watchers.
func (s *Store[S]) collectWatchers() []StoreWatcher[S] {
	watchers := make([]StoreWatcher[S], 0, len(s.watchers))
	for _, w := range s.watchers {
		watchers = append(watchers, w)
	}
	return watchers
}

// notifyWatchers notifies all watchers.
func (s *Store[S]) notifyWatchers(watchers []StoreWatcher[S], oldState, newState S) {
	for _, watcher := range watchers {
		watcher(oldState, newState)
	}
}
