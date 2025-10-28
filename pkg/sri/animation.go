package sri

import (
	"sync"
	"time"
)

// Animation represents an animation.
type Animation struct {
	mu         sync.RWMutex
	duration   time.Duration
	easing     EasingFunc
	startTime  time.Time
	running    bool
	onUpdate   func(progress float64)
	onComplete func()
}

// NewAnimation creates a new animation.
func NewAnimation(duration time.Duration, easing EasingFunc) *Animation {
	return &Animation{
		duration: duration,
		easing:   easing,
	}
}

// OnUpdate sets the update callback.
func (a *Animation) OnUpdate(callback func(progress float64)) *Animation {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.onUpdate = callback
	return a
}

// OnComplete sets the completion callback.
func (a *Animation) OnComplete(callback func()) *Animation {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.onComplete = callback
	return a
}

// Start starts the animation.
func (a *Animation) Start() {
	a.mu.Lock()
	a.running = true
	a.startTime = time.Now()
	a.mu.Unlock()

	go a.run()
}

// Stop stops the animation.
func (a *Animation) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.running = false
}

// IsRunning returns whether the animation is running.
func (a *Animation) IsRunning() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.running
}

// run executes the animation loop.
func (a *Animation) run() {
	ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
	defer ticker.Stop()

	for range ticker.C {
		a.mu.RLock()
		if !a.running {
			a.mu.RUnlock()
			return
		}

		elapsed := time.Since(a.startTime)
		progress := float64(elapsed) / float64(a.duration)

		if progress >= 1.0 {
			progress = 1.0
			a.mu.RUnlock()

			// Final update
			a.mu.RLock()
			onUpdate := a.onUpdate
			onComplete := a.onComplete
			a.mu.RUnlock()

			if onUpdate != nil {
				onUpdate(1.0)
			}

			a.mu.Lock()
			a.running = false
			a.mu.Unlock()

			if onComplete != nil {
				onComplete()
			}

			return
		}

		easing := a.easing
		onUpdate := a.onUpdate
		a.mu.RUnlock()

		// Apply easing
		easedProgress := progress
		if easing != nil {
			easedProgress = easing(progress)
		}

		if onUpdate != nil {
			onUpdate(easedProgress)
		}
	}
}

// Animator manages multiple animations.
type Animator struct {
	mu         sync.RWMutex
	animations map[string]*Animation
}

// NewAnimator creates a new animator.
func NewAnimator() *Animator {
	return &Animator{
		animations: make(map[string]*Animation),
	}
}

// Add adds an animation with a key.
func (am *Animator) Add(key string, animation *Animation) {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.animations[key] = animation
}

// Start starts an animation by key.
func (am *Animator) Start(key string) {
	am.mu.RLock()
	animation, exists := am.animations[key]
	am.mu.RUnlock()

	if exists {
		animation.Start()
	}
}

// Stop stops an animation by key.
func (am *Animator) Stop(key string) {
	am.mu.RLock()
	animation, exists := am.animations[key]
	am.mu.RUnlock()

	if exists {
		animation.Stop()
	}
}

// StopAll stops all animations.
func (am *Animator) StopAll() {
	am.mu.RLock()
	animations := make([]*Animation, 0, len(am.animations))
	for _, anim := range am.animations {
		animations = append(animations, anim)
	}
	am.mu.RUnlock()

	for _, anim := range animations {
		anim.Stop()
	}
}

// Remove removes an animation by key.
func (am *Animator) Remove(key string) {
	am.mu.Lock()
	defer am.mu.Unlock()
	if anim, exists := am.animations[key]; exists {
		anim.Stop()
		delete(am.animations, key)
	}
}

// FadeIn creates a fade-in animation.
func FadeIn(duration time.Duration, onUpdate func(opacity float64)) *Animation {
	return NewAnimation(duration, EaseInOut).OnUpdate(onUpdate)
}

// FadeOut creates a fade-out animation.
func FadeOut(duration time.Duration, onUpdate func(opacity float64)) *Animation {
	return NewAnimation(duration, EaseInOut).OnUpdate(func(progress float64) {
		onUpdate(1.0 - progress)
	})
}

// SlideIn creates a slide-in animation.
func SlideIn(duration time.Duration, distance int, onUpdate func(offset int)) *Animation {
	return NewAnimation(duration, EaseOut).OnUpdate(func(progress float64) {
		offset := distance - int(float64(distance)*progress)
		onUpdate(offset)
	})
}

// Pulse creates a pulse animation.
func Pulse(duration time.Duration, onUpdate func(scale float64)) *Animation {
	return NewAnimation(duration, EaseInOut).OnUpdate(func(progress float64) {
		// Sine wave for pulsing effect
		scale := 1.0 + 0.1*sinApprox(progress*2*3.14159)
		onUpdate(scale)
	})
}

// sinApprox approximates sine function.
func sinApprox(x float64) float64 {
	// Simple sine approximation using Taylor series (first 3 terms)
	x2 := x * x
	return x - (x*x2)/6.0 + (x*x2*x2)/120.0
}
