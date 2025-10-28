package agni

// Priority represents the priority level of an event handler.
type Priority int

// Priority constants.
const (
	// PriorityLow is for low-priority handlers.
	PriorityLow Priority = iota
	// PriorityNormal is the default priority.
	PriorityNormal
	// PriorityHigh is for high-priority handlers.
	PriorityHigh
	// PriorityCritical is for critical system handlers.
	PriorityCritical
)

// String returns the string representation of the priority.
func (p Priority) String() string {
	switch p {
	case PriorityLow:
		return "low"
	case PriorityNormal:
		return "normal"
	case PriorityHigh:
		return "high"
	case PriorityCritical:
		return "critical"
	default:
		return "unknown"
	}
}

// HandlerOptions configures an event handler.
type HandlerOptions struct {
	Priority Priority
	OneShot  bool
	Filter   EventFilter
}

// EventFilter filters events before passing to the handler.
type EventFilter func(Event) bool

// DefaultHandlerOptions returns default handler options.
func DefaultHandlerOptions() HandlerOptions {
	return HandlerOptions{
		Priority: PriorityNormal,
		OneShot:  false,
		Filter:   nil,
	}
}

// WithPriority sets the handler priority.
func WithPriority(p Priority) func(*HandlerOptions) {
	return func(opts *HandlerOptions) {
		opts.Priority = p
	}
}

// WithOneShot makes the handler execute only once.
func WithOneShot() func(*HandlerOptions) {
	return func(opts *HandlerOptions) {
		opts.OneShot = true
	}
}

// WithFilter sets an event filter.
func WithFilter(filter EventFilter) func(*HandlerOptions) {
	return func(opts *HandlerOptions) {
		opts.Filter = filter
	}
}
