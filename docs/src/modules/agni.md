# Agni - Event Hub

**Package:** `github.com/TIVerse/drav/pkg/agni`

## Overview

Agni (Sanskrit: अग्नि, "fire, transformative force") is DRAV's event system. It provides priority-based event dispatching, keyboard/mouse event handling, and timer management.

## Key Concepts

### Event Types

Built-in event types:

- `EventTypeKey`: Keyboard input
- `EventTypeMouse`: Mouse input
- `EventTypeResize`: Terminal resize
- `EventTypeTick`: Timer tick
- `EventTypeQuit`: Application quit
- `EventTypeFocus`: Component focus gained
- `EventTypeBlur`: Component focus lost
- `EventTypeCustom`: User-defined events

### Event Flow

1. Terminal events captured by tcell
2. Converted to Agni events
3. Dispatched through priority queue
4. Handlers execute concurrently (worker pool)

### Priority Levels

Events can be prioritized:

```go
type Priority int

const (
    PriorityLow    Priority = 0
    PriorityNormal Priority = 50
    PriorityHigh   Priority = 100
)
```

## Core API

### Event Hub

Access through the app:

```go
eventHub := app.EventHub()
```

### Registering Handlers

```go
unsubscribe := eventHub.On(agni.EventTypeKey, func(ctx context.Context, event dravya.Event) error {
    keyEvent := event.(*agni.KeyEvent)
    // Handle key event
    return nil
})

// Clean up
defer unsubscribe()
```

### Emitting Events

```go
customEvent := agni.NewCustomEvent("my-event", payload)
eventHub.Emit(ctx, customEvent)
```

## Event Types

### KeyEvent

Keyboard input event:

```go
type KeyEvent struct {
    Key       Key       // Special key (Enter, Tab, etc.)
    Rune      rune      // Character for regular keys
    Modifiers ModMask   // Ctrl, Alt, Shift, Meta
    Repeat    bool      // Is this a repeat event?
}
```

Key constants:

```go
agni.KeyRune        // Regular character
agni.KeyEnter
agni.KeyBackspace
agni.KeyTab
agni.KeyEscape
agni.KeyUp, agni.KeyDown, agni.KeyLeft, agni.KeyRight
agni.KeyHome, agni.KeyEnd
agni.KeyPageUp, agni.KeyPageDown
agni.KeyDelete, agni.KeyInsert
agni.KeyF1, agni.KeyF2, ... agni.KeyF12
agni.KeyCtrlC, agni.KeyCtrlD, agni.KeyCtrlZ
```

Modifiers:

```go
agni.ModShift
agni.ModCtrl
agni.ModAlt
agni.ModMeta
```

### MouseEvent

Mouse input event:

```go
type MouseEvent struct {
    X         int          // Column position
    Y         int          // Row position
    Button    MouseButton  // Which button
    Action    MouseAction  // Press, release, move, etc.
    Wheel     int          // Scroll amount (1=up, -1=down)
    Modifiers ModMask      // Modifier keys held
}
```

Buttons:

```go
agni.MouseNone
agni.MouseLeft
agni.MouseMiddle
agni.MouseRight
```

Actions:

```go
agni.MousePress
agni.MouseRelease
agni.MouseMove
agni.MouseDrag
agni.MouseScroll
```

### ResizeEvent

Terminal resize event:

```go
type ResizeEvent struct {
    Width  int  // New width in columns
    Height int  // New height in rows
}
```

### CustomEvent

User-defined events:

```go
type CustomEvent struct {
    Name    string
    Payload any
}

event := agni.NewCustomEvent("data-loaded", data)
```

## Patterns

### Keyboard Handler

```go
app.OnReady(func() {
    eventHub := app.EventHub()
    eventHub.On(agni.EventTypeKey, func(ctx context.Context, event dravya.Event) error {
        keyEvent, ok := event.(*agni.KeyEvent)
        if !ok {
            return nil
        }
        
        switch keyEvent.Key {
        case agni.KeyEnter:
            return handleEnter()
        case agni.KeyRune:
            if keyEvent.Rune == 'q' {
                return handleQuit()
            }
        }
        
        return nil
    })
})
```

### Modifier Keys

```go
eventHub.On(agni.EventTypeKey, func(ctx context.Context, event dravya.Event) error {
    keyEvent := event.(*agni.KeyEvent)
    
    // Ctrl+S
    if keyEvent.Rune == 's' && keyEvent.Modifiers&agni.ModCtrl != 0 {
        return handleSave()
    }
    
    // Shift+Tab
    if keyEvent.Key == agni.KeyTab && keyEvent.Modifiers&agni.ModShift != 0 {
        return handleShiftTab()
    }
    
    return nil
})
```

### Mouse Handler

```go
eventHub.On(agni.EventTypeMouse, func(ctx context.Context, event dravya.Event) error {
    mouseEvent := event.(*agni.MouseEvent)
    
    switch mouseEvent.Action {
    case agni.MousePress:
        if mouseEvent.Button == agni.MouseLeft {
            return handleClick(mouseEvent.X, mouseEvent.Y)
        }
    case agni.MouseScroll:
        if mouseEvent.Wheel > 0 {
            return handleScrollUp()
        } else {
            return handleScrollDown()
        }
    }
    
    return nil
})
```

### Custom Events

```go
// Component A emits event
func (c *ComponentA) loadData() {
    data := fetchData()
    event := agni.NewCustomEvent("data-loaded", data)
    c.eventHub.Emit(context.Background(), event)
}

// Component B handles event
func (c *ComponentB) init() {
    c.eventHub.On(agni.EventTypeCustom, func(ctx context.Context, event dravya.Event) error {
        customEvent := event.(*agni.CustomEvent)
        if customEvent.Name == "data-loaded" {
            data := customEvent.Payload.(MyData)
            c.updateWithData(data)
        }
        return nil
    })
}
```

## Timers

### One-time Timer

```go
dispatcher := // get dispatcher
dispatcher.After(2*time.Second, func() {
    // Runs once after 2 seconds
})
```

### Recurring Timer

```go
stop := dispatcher.Every(1*time.Second, func() {
    // Runs every second
})

// Stop timer when done
defer stop()
```

## Dispatcher

Direct access to the dispatcher for advanced use:

```go
dispatcher := agni.NewDispatcher(
    1000,  // Queue size
    10,    // Worker pool size
)

// Register with priority
dispatcher.OnWithPriority(
    agni.EventTypeKey,
    agni.PriorityHigh,
    handler,
)

// Start processing
dispatcher.Start(ctx)
```

## Best Practices

### 1. Clean Up Subscriptions

Always unsubscribe when done:

```go
type MyComponent struct {
    cleanup []func()
}

func (c *MyComponent) init(eventHub dravya.EventHub) {
    unsub := eventHub.On(agni.EventTypeKey, c.handleKey)
    c.cleanup = append(c.cleanup, unsub)
}

func (c *MyComponent) dispose() {
    for _, fn := range c.cleanup {
        fn()
    }
}
```

### 2. Use OnReady Hook

Register handlers after app initialization:

```go
app.OnReady(func() {
    eventHub := app.EventHub()
    eventHub.On(agni.EventTypeKey, handler)
})
```

### 3. Handle Errors

Return errors from handlers for logging:

```go
eventHub.On(agni.EventTypeKey, func(ctx context.Context, event dravya.Event) error {
    if err := processKey(event); err != nil {
        return fmt.Errorf("failed to process key: %w", err)
    }
    return nil
})
```

### 4. Avoid Blocking

Handlers execute in worker pool - don't block:

```go
// Bad - blocks event processing
eventHub.On(agni.EventTypeKey, func(ctx context.Context, event dravya.Event) error {
    time.Sleep(5 * time.Second)  // Blocks worker!
    return nil
})

// Good - dispatch to goroutine
eventHub.On(agni.EventTypeKey, func(ctx context.Context, event dravya.Event) error {
    go processExpensiveTask(event)
    return nil
})
```

### 5. Type Assertions

Always check type assertions:

```go
eventHub.On(agni.EventTypeKey, func(ctx context.Context, event dravya.Event) error {
    keyEvent, ok := event.(*agni.KeyEvent)
    if !ok {
        return nil  // Wrong type, ignore
    }
    // Use keyEvent safely
    return nil
})
```

## Examples

### Keyboard Shortcuts

```go
type ShortcutHandler struct {
    shortcuts map[string]func() error
}

func (s *ShortcutHandler) register(eventHub dravya.EventHub) {
    s.shortcuts = map[string]func() error{
        "ctrl+s": s.save,
        "ctrl+q": s.quit,
        "ctrl+n": s.new,
    }
    
    eventHub.On(agni.EventTypeKey, func(ctx context.Context, event dravya.Event) error {
        keyEvent := event.(*agni.KeyEvent)
        key := s.keyString(keyEvent)
        
        if handler, ok := s.shortcuts[key]; ok {
            return handler()
        }
        return nil
    })
}

func (s *ShortcutHandler) keyString(e *agni.KeyEvent) string {
    var parts []string
    if e.Modifiers&agni.ModCtrl != 0 {
        parts = append(parts, "ctrl")
    }
    if e.Modifiers&agni.ModShift != 0 {
        parts = append(parts, "shift")
    }
    if e.Key == agni.KeyRune {
        parts = append(parts, string(e.Rune))
    }
    return strings.Join(parts, "+")
}
```

## Performance Considerations

### Queue Size

Larger queue handles burst events:

```go
dispatcher := agni.NewDispatcher(
    5000,  // Large queue for high event rate
    20,    // More workers
)
```

### Worker Pool

Match workers to CPU cores:

```go
workers := runtime.NumCPU()
dispatcher := agni.NewDispatcher(1000, workers)
```

### Event Filtering

Filter early to reduce processing:

```go
eventHub.On(agni.EventTypeKey, func(ctx context.Context, event dravya.Event) error {
    keyEvent := event.(*agni.KeyEvent)
    
    // Ignore repeats early
    if keyEvent.Repeat {
        return nil
    }
    
    return processKey(keyEvent)
})
```

## Troubleshooting

### Events Not Firing

Check handler registration timing:

```go
// Bad - too early
eventHub := app.EventHub()  // Might be nil
eventHub.On(...)

// Good - use OnReady
app.OnReady(func() {
    eventHub := app.EventHub()
    eventHub.On(...)
})
```

### Event Loss

Increase queue size:

```go
dispatcher := agni.NewDispatcher(10000, 10)  // Bigger queue
```

### Handler Deadlock

Don't wait for render in handler:

```go
// Bad - deadlock potential
eventHub.On(agni.EventTypeKey, func(ctx context.Context, event dravya.Event) error {
    app.Render()  // May deadlock
    return nil
})

// Good - request render
eventHub.On(agni.EventTypeKey, func(ctx context.Context, event dravya.Event) error {
    app.RequestRender()  // Non-blocking
    return nil
})
```

## Related Modules

- **[Dravya](dravya.md)**: Provides EventHub interface
- **[Maya](maya.md)**: Focus events

## See Also

- [Counter Example](../../examples/02-counter/)
- [Event Handling Guide](../concepts.md#events)
