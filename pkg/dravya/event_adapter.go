package dravya

import (
	"context"
	"time"

	"github.com/TIVerse/drav/pkg/agni"
)

// agniEventAdapter adapts agni.Dispatcher to the EventHub interface.
type agniEventAdapter struct {
	dispatcher *agni.Dispatcher
}

// NewAgniEventHub creates an EventHub from an agni.Dispatcher.
func NewAgniEventHub(dispatcher *agni.Dispatcher) EventHub {
	return &agniEventAdapter{
		dispatcher: dispatcher,
	}
}

// NewDefaultEventHub creates a default EventHub with sensible defaults.
func NewDefaultEventHub() EventHub {
	// Queue size: 1000 events
	// Worker pool: 10 concurrent workers
	dispatcher := agni.NewDispatcher(1000, 10)
	return NewAgniEventHub(dispatcher)
}

// Emit emits an event to all registered handlers.
func (a *agniEventAdapter) Emit(ctx context.Context, event Event) error {
	// Convert dravya.Event to agni.Event
	agniEvent := wrapEvent(event)
	return a.dispatcher.Emit(ctx, agniEvent)
}

// On registers an event handler.
func (a *agniEventAdapter) On(eventType string, handler EventHandler) (unsubscribe func()) {
	// Convert dravya.EventHandler to agni.Handler
	agniHandler := func(ctx context.Context, event agni.Event) error {
		// Wrap agni.Event back to dravya.Event
		dravyaEvent := unwrapEvent(event)
		return handler(ctx, dravyaEvent)
	}
	
	return a.dispatcher.On(eventType, agniHandler)
}

// Start starts the event hub.
func (a *agniEventAdapter) Start(ctx context.Context) error {
	return a.dispatcher.Start(ctx)
}

// Stop stops the event hub.
func (a *agniEventAdapter) Stop() error {
	return a.dispatcher.Stop()
}

// eventWrapper wraps a dravya.Event to implement agni.Event.
type eventWrapper struct {
	inner Event
}

func (e *eventWrapper) Type() string {
	return e.inner.Type()
}

func (e *eventWrapper) Time() time.Time {
	return e.inner.Time()
}

func wrapEvent(event Event) agni.Event {
	// If it's already an agni.Event, return as-is
	if agniEvent, ok := event.(agni.Event); ok {
		return agniEvent
	}
	return &eventWrapper{inner: event}
}

func unwrapEvent(event agni.Event) Event {
	// If it's wrapped, unwrap it
	if wrapper, ok := event.(*eventWrapper); ok {
		return wrapper.inner
	}
	// Otherwise, create a simple adapter
	return &simpleEvent{
		typ:  event.Type(),
		time: event.Time(),
	}
}

// simpleEvent is a basic Event implementation.
type simpleEvent struct {
	typ  string
	time time.Time
}

func (e *simpleEvent) Type() string {
	return e.typ
}

func (e *simpleEvent) Time() time.Time {
	return e.time
}

// NewSimpleEvent creates a simple event with just type and time.
func NewSimpleEvent(typ string) Event {
	return &simpleEvent{
		typ:  typ,
		time: time.Now(),
	}
}
