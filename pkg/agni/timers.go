package agni

import (
	"context"
	"sync"
	"time"
)

// Timer represents a scheduled timer.
type Timer struct {
	id       string
	interval time.Duration
	repeat   bool
	callback func(ctx context.Context)
	ticker   *time.Ticker
	cancel   context.CancelFunc
	mu       sync.Mutex
}

// NewTimer creates a new timer.
func NewTimer(id string, interval time.Duration, repeat bool, callback func(ctx context.Context)) *Timer {
	return &Timer{
		id:       id,
		interval: interval,
		repeat:   repeat,
		callback: callback,
	}
}

// Start starts the timer.
func (t *Timer) Start(ctx context.Context) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.cancel != nil {
		return // Already started
	}

	timerCtx, cancel := context.WithCancel(ctx)
	t.cancel = cancel

	if t.repeat {
		t.ticker = time.NewTicker(t.interval)
		go t.runRepeating(timerCtx)
	} else {
		go t.runOnce(timerCtx)
	}
}

// Stop stops the timer.
func (t *Timer) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.cancel != nil {
		t.cancel()
		t.cancel = nil
	}

	if t.ticker != nil {
		t.ticker.Stop()
		t.ticker = nil
	}
}

// runOnce runs the timer callback once.
func (t *Timer) runOnce(ctx context.Context) {
	select {
	case <-time.After(t.interval):
		t.callback(ctx)
	case <-ctx.Done():
		return
	}
}

// runRepeating runs the timer callback repeatedly.
func (t *Timer) runRepeating(ctx context.Context) {
	for {
		select {
		case <-t.ticker.C:
			t.callback(ctx)
		case <-ctx.Done():
			return
		}
	}
}

// ID returns the timer ID.
func (t *Timer) ID() string {
	return t.id
}

// TimerManager manages timers.
type TimerManager struct {
	mu     sync.RWMutex
	timers map[string]*Timer
}

// NewTimerManager creates a new timer manager.
func NewTimerManager() *TimerManager {
	return &TimerManager{
		timers: make(map[string]*Timer),
	}
}

// After schedules a one-shot timer.
func (tm *TimerManager) After(ctx context.Context, id string, duration time.Duration, callback func(ctx context.Context)) *Timer {
	timer := NewTimer(id, duration, false, callback)
	
	tm.mu.Lock()
	tm.timers[id] = timer
	tm.mu.Unlock()
	
	timer.Start(ctx)
	return timer
}

// Every schedules a repeating timer.
func (tm *TimerManager) Every(ctx context.Context, id string, interval time.Duration, callback func(ctx context.Context)) *Timer {
	timer := NewTimer(id, interval, true, callback)
	
	tm.mu.Lock()
	tm.timers[id] = timer
	tm.mu.Unlock()
	
	timer.Start(ctx)
	return timer
}

// Cancel cancels a timer by ID.
func (tm *TimerManager) Cancel(id string) bool {
	tm.mu.Lock()
	timer, exists := tm.timers[id]
	if exists {
		delete(tm.timers, id)
	}
	tm.mu.Unlock()

	if exists {
		timer.Stop()
		return true
	}
	return false
}

// CancelAll cancels all timers.
func (tm *TimerManager) CancelAll() {
	tm.mu.Lock()
	timers := make([]*Timer, 0, len(tm.timers))
	for _, timer := range tm.timers {
		timers = append(timers, timer)
	}
	tm.timers = make(map[string]*Timer)
	tm.mu.Unlock()

	for _, timer := range timers {
		timer.Stop()
	}
}

// Count returns the number of active timers.
func (tm *TimerManager) Count() int {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return len(tm.timers)
}
