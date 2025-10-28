package prana

import (
	"context"
	"time"
)

// EffectRunner manages async effects.
type EffectRunner struct {
	pending chan effectTask
	workers int
}

// effectTask represents a pending effect.
type effectTask struct {
	ctx    context.Context
	effect func(ctx context.Context) error
	done   chan error
}

// NewEffectRunner creates a new effect runner.
func NewEffectRunner(workers int) *EffectRunner {
	er := &EffectRunner{
		pending: make(chan effectTask, 100),
		workers: workers,
	}
	er.start()
	return er
}

// start starts the worker pool.
func (er *EffectRunner) start() {
	for i := 0; i < er.workers; i++ {
		go er.worker()
	}
}

// worker processes effects from the queue.
func (er *EffectRunner) worker() {
	for task := range er.pending {
		err := task.effect(task.ctx)
		select {
		case task.done <- err:
		case <-task.ctx.Done():
		}
	}
}

// Run runs an effect asynchronously.
func (er *EffectRunner) Run(ctx context.Context, effect func(ctx context.Context) error) <-chan error {
	done := make(chan error, 1)
	task := effectTask{
		ctx:    ctx,
		effect: effect,
		done:   done,
	}

	select {
	case er.pending <- task:
	case <-ctx.Done():
		done <- ctx.Err()
	}

	return done
}

// RunWithTimeout runs an effect with a timeout.
func (er *EffectRunner) RunWithTimeout(ctx context.Context, timeout time.Duration, effect func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	errCh := er.Run(ctx, effect)
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
