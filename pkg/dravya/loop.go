package dravya

import (
	"context"
	"sync"
	"time"
)

// Loop represents the main event loop coordinator.
type Loop struct {
	mu             sync.Mutex
	frameDuration  time.Duration
	lastFrameTime  time.Time
	frameCount     uint64
	droppedFrames  uint64
	running        bool
	ticker         *time.Ticker
	stopCh         chan struct{}
	frameCallbacks []FrameCallback
}

// FrameCallback is called on each frame.
type FrameCallback func(ctx context.Context, frameTime time.Time, delta time.Duration) error

// NewLoop creates a new event loop.
func NewLoop(fps int) *Loop {
	return &Loop{
		frameDuration: time.Second / time.Duration(fps),
		stopCh:        make(chan struct{}),
	}
}

// OnFrame registers a callback to be called each frame.
func (l *Loop) OnFrame(callback FrameCallback) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.frameCallbacks = append(l.frameCallbacks, callback)
}

// Start begins the event loop.
func (l *Loop) Start(ctx context.Context) error {
	l.mu.Lock()
	if l.running {
		l.mu.Unlock()
		return nil
	}
	l.running = true
	l.ticker = time.NewTicker(l.frameDuration)
	l.lastFrameTime = time.Now()
	l.mu.Unlock()

	defer l.ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-l.stopCh:
			return nil
		case frameTime := <-l.ticker.C:
			if err := l.processFrame(ctx, frameTime); err != nil {
				return err
			}
		}
	}
}

// Stop stops the event loop.
func (l *Loop) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.running {
		close(l.stopCh)
		l.running = false
	}
}

// processFrame processes a single frame.
func (l *Loop) processFrame(ctx context.Context, frameTime time.Time) error {
	l.mu.Lock()
	delta := frameTime.Sub(l.lastFrameTime)
	l.lastFrameTime = frameTime
	l.frameCount++

	// Check for dropped frames
	if delta > l.frameDuration*2 {
		dropped := uint64(delta / l.frameDuration)
		l.droppedFrames += dropped
	}

	callbacks := l.frameCallbacks
	l.mu.Unlock()

	// Execute frame callbacks
	for _, callback := range callbacks {
		if err := callback(ctx, frameTime, delta); err != nil {
			return err
		}
	}

	return nil
}

// Stats returns loop statistics.
func (l *Loop) Stats() LoopStats {
	l.mu.Lock()
	defer l.mu.Unlock()
	return LoopStats{
		FrameCount:    l.frameCount,
		DroppedFrames: l.droppedFrames,
		FrameDuration: l.frameDuration,
		Running:       l.running,
	}
}

// LoopStats contains event loop statistics.
type LoopStats struct {
	FrameCount    uint64
	DroppedFrames uint64
	FrameDuration time.Duration
	Running       bool
}

// SetFPS updates the target FPS at runtime.
func (l *Loop) SetFPS(fps int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.frameDuration = time.Second / time.Duration(fps)
	if l.ticker != nil {
		l.ticker.Reset(l.frameDuration)
	}
}
