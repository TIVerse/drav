package sandbox

import (
	"context"
	"fmt"
	"time"
)

// Runner executes functions with resource limits.
type Runner struct {
	timeout       time.Duration
	maxGoroutines int
}

// NewRunner creates a new sandbox runner.
func NewRunner(timeout time.Duration, maxGoroutines int) *Runner {
	return &Runner{
		timeout:       timeout,
		maxGoroutines: maxGoroutines,
	}
}

// Run executes a function with timeout and resource limits.
func (r *Runner) Run(ctx context.Context, fn func(context.Context) error) error {
	if r.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
		defer cancel()
	}

	errCh := make(chan error, 1)

	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				errCh <- fmt.Errorf("panic in sandboxed function: %v", rec)
			}
		}()

		errCh <- fn(ctx)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
