package vayu

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Sandbox provides resource limits and isolation for plugins.
type Sandbox struct {
	mu              sync.RWMutex
	maxMemory       uint64
	maxGoroutines   int
	timeout         time.Duration
	capChecker      *CapabilityChecker
	goroutineCount  int
	memoryUsed      uint64
}

// SandboxOptions configures a sandbox.
type SandboxOptions struct {
	MaxMemory     uint64
	MaxGoroutines int
	Timeout       time.Duration
	Capabilities  Capabilities
}

// NewSandbox creates a new sandbox.
func NewSandbox(opts SandboxOptions) *Sandbox {
	return &Sandbox{
		maxMemory:     opts.MaxMemory,
		maxGoroutines: opts.MaxGoroutines,
		timeout:       opts.Timeout,
		capChecker:    NewCapabilityChecker(opts.Capabilities),
	}
}

// Run runs a function in the sandbox.
func (s *Sandbox) Run(ctx context.Context, fn func(context.Context) error) error {
	// Check goroutine limit
	s.mu.Lock()
	if s.goroutineCount >= s.maxGoroutines {
		s.mu.Unlock()
		return fmt.Errorf("goroutine limit exceeded")
	}
	s.goroutineCount++
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		s.goroutineCount--
		s.mu.Unlock()
	}()

	// Apply timeout
	if s.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.timeout)
		defer cancel()
	}

	// Run function
	return fn(ctx)
}

// CheckCapability checks if an operation is allowed.
func (s *Sandbox) CheckCapability(op string, target string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	switch op {
	case "fs:read":
		if !s.capChecker.CanReadPath(target) {
			return fmt.Errorf("filesystem read access denied: %s", target)
		}
	case "fs:write":
		if !s.capChecker.CanWritePath(target) {
			return fmt.Errorf("filesystem write access denied: %s", target)
		}
	case "net:access":
		if !s.capChecker.CanAccessDomain(target) {
			return fmt.Errorf("network access denied: %s", target)
		}
	default:
		return fmt.Errorf("unknown operation: %s", op)
	}

	return nil
}

// Stats returns sandbox statistics.
func (s *Sandbox) Stats() SandboxStats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return SandboxStats{
		GoroutineCount: s.goroutineCount,
		MemoryUsed:     s.memoryUsed,
		MaxMemory:      s.maxMemory,
		MaxGoroutines:  s.maxGoroutines,
	}
}

// SandboxStats contains sandbox statistics.
type SandboxStats struct {
	GoroutineCount int
	MemoryUsed     uint64
	MaxMemory      uint64
	MaxGoroutines  int
}
