package dravya

import (
	"context"
	"errors"
	"sync"
	"time"
)

// LifecycleState represents the current state of the application.
type LifecycleState int

const (
	// StateUninitialized is the initial state before initialization.
	StateUninitialized LifecycleState = iota
	// StateInitializing is when the app is being initialized.
	StateInitializing
	// StateRunning is when the app is running normally.
	StateRunning
	// StateShuttingDown is when the app is gracefully shutting down.
	StateShuttingDown
	// StateTerminated is when the app has fully shut down.
	StateTerminated
)

// String returns the string representation of the lifecycle state.
func (s LifecycleState) String() string {
	switch s {
	case StateUninitialized:
		return "uninitialized"
	case StateInitializing:
		return "initializing"
	case StateRunning:
		return "running"
	case StateShuttingDown:
		return "shutting_down"
	case StateTerminated:
		return "terminated"
	default:
		return "unknown"
	}
}

// Lifecycle manages the application lifecycle and graceful shutdown.
type Lifecycle struct {
	mu             sync.RWMutex
	state          LifecycleState
	hooks          map[LifecycleState][]Hook
	shutdownSignal chan struct{}
	shutdownOnce   sync.Once
	shutdownErr    error
}

// Hook is a function called during lifecycle transitions.
type Hook func(ctx context.Context) error

// NewLifecycle creates a new lifecycle manager.
func NewLifecycle() *Lifecycle {
	return &Lifecycle{
		state:          StateUninitialized,
		hooks:          make(map[LifecycleState][]Hook),
		shutdownSignal: make(chan struct{}),
	}
}

// State returns the current lifecycle state.
func (lc *Lifecycle) State() LifecycleState {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	return lc.state
}

// SetState transitions to a new state and runs hooks.
func (lc *Lifecycle) SetState(ctx context.Context, newState LifecycleState) error {
	lc.mu.Lock()
	oldState := lc.state
	lc.state = newState
	hooks := lc.hooks[newState]
	lc.mu.Unlock()

	if oldState == newState {
		return nil
	}

	// Run hooks for this state
	for _, hook := range hooks {
		if err := hook(ctx); err != nil {
			return err
		}
	}

	return nil
}

// OnState registers a hook to be called when entering a specific state.
func (lc *Lifecycle) OnState(state LifecycleState, hook Hook) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.hooks[state] = append(lc.hooks[state], hook)
}

// Shutdown initiates graceful shutdown.
func (lc *Lifecycle) Shutdown(err error) {
	lc.shutdownOnce.Do(func() {
		lc.shutdownErr = err
		close(lc.shutdownSignal)
	})
}

// ShutdownSignal returns a channel that is closed when shutdown is initiated.
func (lc *Lifecycle) ShutdownSignal() <-chan struct{} {
	return lc.shutdownSignal
}

// ShutdownError returns the error that caused shutdown, if any.
func (lc *Lifecycle) ShutdownError() error {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	return lc.shutdownErr
}

// WaitForShutdown blocks until shutdown is initiated, with optional timeout.
func (lc *Lifecycle) WaitForShutdown(timeout time.Duration) error {
	if timeout == 0 {
		<-lc.shutdownSignal
		return lc.ShutdownError()
	}

	select {
	case <-lc.shutdownSignal:
		return lc.ShutdownError()
	case <-time.After(timeout):
		return errors.New("shutdown timeout exceeded")
	}
}

// IsRunning returns true if the app is in the running state.
func (lc *Lifecycle) IsRunning() bool {
	return lc.State() == StateRunning
}

// IsShuttingDown returns true if the app is shutting down or terminated.
func (lc *Lifecycle) IsShuttingDown() bool {
	state := lc.State()
	return state == StateShuttingDown || state == StateTerminated
}
