package agni

import (
	"time"
)

// Event is the base interface for all events.
type Event interface {
	Type() string
	Time() time.Time
}

// BaseEvent provides common event fields.
type BaseEvent struct {
	EventType string
	Timestamp time.Time
}

// Type returns the event type.
func (e *BaseEvent) Type() string {
	return e.EventType
}

// Time returns the event timestamp.
func (e *BaseEvent) Time() time.Time {
	return e.Timestamp
}

// EventType constants.
const (
	EventTypeKey     = "key"
	EventTypeMouse   = "mouse"
	EventTypeResize  = "resize"
	EventTypeTick    = "tick"
	EventTypeCustom  = "custom"
	EventTypeQuit    = "quit"
	EventTypeFocus   = "focus"
	EventTypeBlur    = "blur"
)

// KeyEvent represents a keyboard event.
type KeyEvent struct {
	BaseEvent
	Key       Key
	Rune      rune
	Modifiers ModMask
	Repeat    bool
}

// NewKeyEvent creates a new key event.
func NewKeyEvent(key Key, r rune, mods ModMask, repeat bool) *KeyEvent {
	return &KeyEvent{
		BaseEvent: BaseEvent{
			EventType: EventTypeKey,
			Timestamp: time.Now(),
		},
		Key:       key,
		Rune:      r,
		Modifiers: mods,
		Repeat:    repeat,
	}
}

// Key represents a keyboard key.
type Key int

// Key constants.
const (
	KeyRune Key = iota
	KeyEnter
	KeyBackspace
	KeyTab
	KeyEscape
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyHome
	KeyEnd
	KeyPageUp
	KeyPageDown
	KeyDelete
	KeyInsert
	KeyF1
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyCtrlC
	KeyCtrlD
	KeyCtrlZ
)

// ModMask represents modifier keys.
type ModMask int

// Modifier constants.
const (
	ModNone ModMask = 0
	ModShift ModMask = 1 << iota
	ModCtrl
	ModAlt
	ModMeta
)

// MouseEvent represents a mouse event.
type MouseEvent struct {
	BaseEvent
	X       int
	Y       int
	Button  MouseButton
	Action  MouseAction
	Wheel   int
	Modifiers ModMask
}

// NewMouseEvent creates a new mouse event.
func NewMouseEvent(x, y int, button MouseButton, action MouseAction, wheel int, mods ModMask) *MouseEvent {
	return &MouseEvent{
		BaseEvent: BaseEvent{
			EventType: EventTypeMouse,
			Timestamp: time.Now(),
		},
		X:         x,
		Y:         y,
		Button:    button,
		Action:    action,
		Wheel:     wheel,
		Modifiers: mods,
	}
}

// MouseButton represents a mouse button.
type MouseButton int

// Mouse button constants.
const (
	MouseNone MouseButton = iota
	MouseLeft
	MouseMiddle
	MouseRight
)

// MouseAction represents a mouse action.
type MouseAction int

// Mouse action constants.
const (
	MousePress MouseAction = iota
	MouseRelease
	MouseMove
	MouseDrag
	MouseScroll
)

// ResizeEvent represents a terminal resize event.
type ResizeEvent struct {
	BaseEvent
	Width  int
	Height int
}

// NewResizeEvent creates a new resize event.
func NewResizeEvent(width, height int) *ResizeEvent {
	return &ResizeEvent{
		BaseEvent: BaseEvent{
			EventType: EventTypeResize,
			Timestamp: time.Now(),
		},
		Width:  width,
		Height: height,
	}
}

// TickEvent represents a periodic tick event.
type TickEvent struct {
	BaseEvent
	TickID   string
	Interval time.Duration
}

// NewTickEvent creates a new tick event.
func NewTickEvent(id string, interval time.Duration) *TickEvent {
	return &TickEvent{
		BaseEvent: BaseEvent{
			EventType: EventTypeTick,
			Timestamp: time.Now(),
		},
		TickID:   id,
		Interval: interval,
	}
}

// CustomEvent represents a custom application event.
type CustomEvent struct {
	BaseEvent
	Name    string
	Payload any
}

// NewCustomEvent creates a new custom event.
func NewCustomEvent(name string, payload any) *CustomEvent {
	return &CustomEvent{
		BaseEvent: BaseEvent{
			EventType: EventTypeCustom,
			Timestamp: time.Now(),
		},
		Name:    name,
		Payload: payload,
	}
}

// QuitEvent represents a quit/shutdown event.
type QuitEvent struct {
	BaseEvent
	Reason string
}

// NewQuitEvent creates a new quit event.
func NewQuitEvent(reason string) *QuitEvent {
	return &QuitEvent{
		BaseEvent: BaseEvent{
			EventType: EventTypeQuit,
			Timestamp: time.Now(),
		},
		Reason: reason,
	}
}

// FocusEvent represents a component focus event.
type FocusEvent struct {
	BaseEvent
	ComponentID string
	Gained      bool
}

// NewFocusEvent creates a new focus event.
func NewFocusEvent(componentID string, gained bool) *FocusEvent {
	eventType := EventTypeFocus
	if !gained {
		eventType = EventTypeBlur
	}
	return &FocusEvent{
		BaseEvent: BaseEvent{
			EventType: eventType,
			Timestamp: time.Now(),
		},
		ComponentID: componentID,
		Gained:      gained,
	}
}
