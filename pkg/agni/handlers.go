package agni

import (
	"context"
	"sync/atomic"
)

// Handler is a function that handles events.
type Handler func(ctx context.Context, event Event) error

// handlerRegistration wraps a handler with metadata.
type handlerRegistration struct {
	id       uint64
	handler  Handler
	options  HandlerOptions
	eventType string
	active   atomic.Bool
}

// newHandlerRegistration creates a new handler registration.
func newHandlerRegistration(id uint64, eventType string, handler Handler, opts HandlerOptions) *handlerRegistration {
	reg := &handlerRegistration{
		id:       id,
		handler:  handler,
		options:  opts,
		eventType: eventType,
	}
	reg.active.Store(true)
	return reg
}

// shouldHandle checks if this handler should process the event.
func (h *handlerRegistration) shouldHandle(event Event) bool {
	if !h.active.Load() {
		return false
	}

	// Check event type match (empty string means handle all events)
	if h.eventType != "" && h.eventType != event.Type() {
		return false
	}

	// Apply filter if set
	if h.options.Filter != nil && !h.options.Filter(event) {
		return false
	}

	return true
}

// handle executes the handler.
func (h *handlerRegistration) handle(ctx context.Context, event Event) error {
	if !h.shouldHandle(event) {
		return nil
	}

	err := h.handler(ctx, event)

	// Deactivate if one-shot
	if h.options.OneShot {
		h.active.Store(false)
	}

	return err
}

// deactivate marks the handler as inactive.
func (h *handlerRegistration) deactivate() {
	h.active.Store(false)
}

// isActive returns true if the handler is active.
func (h *handlerRegistration) isActive() bool {
	return h.active.Load()
}
