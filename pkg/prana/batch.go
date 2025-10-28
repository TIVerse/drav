package prana

import (
	"context"
	"sync"
	"time"
)

// Batch manages batched updates to prevent excessive notifications.
type Batch struct {
	mu        sync.Mutex
	pending   []func()
	scheduled bool
	delay     time.Duration
}

// NewBatch creates a new batch updater.
func NewBatch(delay time.Duration) *Batch {
	return &Batch{
		pending: make([]func(), 0),
		delay:   delay,
	}
}

// Queue queues an update function.
func (b *Batch) Queue(fn func()) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.pending = append(b.pending, fn)

	if !b.scheduled {
		b.scheduled = true
		go b.flush()
	}
}

// flush executes all pending updates.
func (b *Batch) flush() {
	time.Sleep(b.delay)

	b.mu.Lock()
	pending := b.pending
	b.pending = make([]func(), 0)
	b.scheduled = false
	b.mu.Unlock()

	for _, fn := range pending {
		fn()
	}
}

// BatchUpdater provides batched updates for observables.
type BatchUpdater[T any] struct {
	observable *Observable[T]
	batch      *Batch
}

// NewBatchUpdater creates a new batch updater for an observable.
func NewBatchUpdater[T any](obs *Observable[T], delay time.Duration) *BatchUpdater[T] {
	return &BatchUpdater[T]{
		observable: obs,
		batch:      NewBatch(delay),
	}
}

// Set queues a batched update.
func (bu *BatchUpdater[T]) Set(value T) {
	bu.batch.Queue(func() {
		bu.observable.Set(value)
	})
}

// Update queues a batched update function.
func (bu *BatchUpdater[T]) Update(fn func(T) T) {
	bu.batch.Queue(func() {
		bu.observable.Update(fn)
	})
}

// Transaction manages transactional updates across multiple stores.
type Transaction struct {
	mu      sync.Mutex
	updates []func(context.Context) error
	commit  chan struct{}
}

// NewTransaction creates a new transaction.
func NewTransaction() *Transaction {
	return &Transaction{
		updates: make([]func(context.Context) error, 0),
		commit:  make(chan struct{}),
	}
}

// Add adds an update to the transaction.
func (t *Transaction) Add(fn func(context.Context) error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.updates = append(t.updates, fn)
}

// Commit commits all updates atomically.
func (t *Transaction) Commit(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, update := range t.updates {
		if err := update(ctx); err != nil {
			return err
		}
	}

	close(t.commit)
	return nil
}

// Wait waits for the transaction to commit.
func (t *Transaction) Wait() {
	<-t.commit
}
