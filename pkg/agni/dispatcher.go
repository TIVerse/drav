package agni

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Dispatcher manages event distribution to handlers.
type Dispatcher struct {
	mu           sync.RWMutex
	handlers     map[string][]*handlerRegistration // eventType -> handlers
	allHandlers  []*handlerRegistration            // handlers for all events
	nextID       atomic.Uint64
	timerMgr     *TimerManager
	eventQueue   chan Event
	workerPool   chan struct{}
	running      atomic.Bool
	ctx          context.Context
	cancel       context.CancelFunc
	wg           sync.WaitGroup
}

// NewDispatcher creates a new event dispatcher.
func NewDispatcher(queueSize, workerPoolSize int) *Dispatcher {
	ctx, cancel := context.WithCancel(context.Background())
	return &Dispatcher{
		handlers:   make(map[string][]*handlerRegistration),
		timerMgr:   NewTimerManager(),
		eventQueue: make(chan Event, queueSize),
		workerPool: make(chan struct{}, workerPoolSize),
		ctx:        ctx,
		cancel:     cancel,
	}
}

// On registers an event handler.
func (d *Dispatcher) On(eventType string, handler Handler, optFuncs ...func(*HandlerOptions)) func() {
	opts := DefaultHandlerOptions()
	for _, fn := range optFuncs {
		fn(&opts)
	}

	id := d.nextID.Add(1)
	reg := newHandlerRegistration(id, eventType, handler, opts)

	d.mu.Lock()
	defer d.mu.Unlock()

	if eventType == "" {
		// Handler for all events
		d.allHandlers = append(d.allHandlers, reg)
	} else {
		d.handlers[eventType] = append(d.handlers[eventType], reg)
		// Sort by priority
		d.sortHandlers(eventType)
	}

	// Return unsubscribe function
	return func() {
		d.unregister(id, eventType)
	}
}

// Emit emits an event to all registered handlers.
func (d *Dispatcher) Emit(ctx context.Context, event Event) error {
	if !d.running.Load() {
		return fmt.Errorf("dispatcher not running")
	}

	select {
	case d.eventQueue <- event:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	default:
		return fmt.Errorf("event queue full")
	}
}

// Start starts the event dispatcher.
func (d *Dispatcher) Start(ctx context.Context) error {
	if !d.running.CompareAndSwap(false, true) {
		return nil // Already running
	}

	d.wg.Add(1)
	go d.processEvents(ctx)

	return nil
}

// Stop stops the event dispatcher.
func (d *Dispatcher) Stop() error {
	if !d.running.CompareAndSwap(true, false) {
		return nil // Already stopped
	}

	d.cancel()
	d.timerMgr.CancelAll()
	close(d.eventQueue)
	d.wg.Wait()

	return nil
}

// After schedules a one-shot timer.
func (d *Dispatcher) After(ctx context.Context, id string, duration interface{}, callback func(ctx context.Context)) *Timer {
	return d.timerMgr.After(ctx, id, duration.(time.Duration), callback)
}

// Every schedules a repeating timer.
func (d *Dispatcher) Every(ctx context.Context, id string, interval interface{}, callback func(ctx context.Context)) *Timer {
	return d.timerMgr.Every(ctx, id, interval.(time.Duration), callback)
}

// CancelTimer cancels a timer by ID.
func (d *Dispatcher) CancelTimer(id string) bool {
	return d.timerMgr.Cancel(id)
}

// processEvents processes events from the queue.
func (d *Dispatcher) processEvents(ctx context.Context) {
	defer d.wg.Done()

	for {
		select {
		case event, ok := <-d.eventQueue:
			if !ok {
				return
			}
			d.dispatchEvent(ctx, event)
		case <-ctx.Done():
			return
		}
	}
}

// dispatchEvent dispatches an event to all matching handlers.
func (d *Dispatcher) dispatchEvent(ctx context.Context, event Event) {
	d.mu.RLock()
	
	// Get handlers for this specific event type
	handlers := make([]*handlerRegistration, 0)
	if typeHandlers, exists := d.handlers[event.Type()]; exists {
		handlers = append(handlers, typeHandlers...)
	}
	
	// Add handlers that listen to all events
	handlers = append(handlers, d.allHandlers...)
	
	d.mu.RUnlock()

	// Execute handlers
	for _, reg := range handlers {
		if !reg.isActive() {
			continue
		}

		// Acquire worker from pool
		select {
		case d.workerPool <- struct{}{}:
			d.wg.Add(1)
			go func(r *handlerRegistration) {
				defer d.wg.Done()
				defer func() { <-d.workerPool }()
				
				if err := r.handle(ctx, event); err != nil {
					// Log error but continue processing
					// In a real implementation, we'd use a logger here
				}
			}(reg)
		default:
			// Worker pool full, execute synchronously
			if err := reg.handle(ctx, event); err != nil {
				// Log error
			}
		}
	}
}

// unregister removes a handler registration.
func (d *Dispatcher) unregister(id uint64, eventType string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if eventType == "" {
		// Remove from allHandlers
		for i, reg := range d.allHandlers {
			if reg.id == id {
				reg.deactivate()
				d.allHandlers = append(d.allHandlers[:i], d.allHandlers[i+1:]...)
				return
			}
		}
	} else {
		// Remove from specific event type handlers
		if handlers, exists := d.handlers[eventType]; exists {
			for i, reg := range handlers {
				if reg.id == id {
					reg.deactivate()
					d.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
					return
				}
			}
		}
	}
}

// sortHandlers sorts handlers by priority.
func (d *Dispatcher) sortHandlers(eventType string) {
	handlers := d.handlers[eventType]
	// Simple bubble sort by priority (good enough for small lists)
	for i := 0; i < len(handlers); i++ {
		for j := i + 1; j < len(handlers); j++ {
			if handlers[i].options.Priority < handlers[j].options.Priority {
				handlers[i], handlers[j] = handlers[j], handlers[i]
			}
		}
	}
}

// Stats returns dispatcher statistics.
func (d *Dispatcher) Stats() DispatcherStats {
	d.mu.RLock()
	defer d.mu.RUnlock()

	totalHandlers := len(d.allHandlers)
	for _, handlers := range d.handlers {
		totalHandlers += len(handlers)
	}

	return DispatcherStats{
		Running:       d.running.Load(),
		TotalHandlers: totalHandlers,
		ActiveTimers:  d.timerMgr.Count(),
		QueuedEvents:  len(d.eventQueue),
	}
}

// DispatcherStats contains dispatcher statistics.
type DispatcherStats struct {
	Running       bool
	TotalHandlers int
	ActiveTimers  int
	QueuedEvents  int
}
