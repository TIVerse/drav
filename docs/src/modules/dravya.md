# Dravya - Runtime Core

**Package:** `github.com/TIVerse/drav/pkg/dravya`

## Overview

Dravya (Sanskrit: द्रव्य, "substance, essence") is the runtime core of DRAV. It manages the application lifecycle, orchestrates the event loop, and coordinates all other modules.

## Key Concepts

### Application Instance

The `App` struct is the central coordinator:

```go
app := dravya.NewApp(
    dravya.WithLogLevel(slog.LevelInfo),
    dravya.WithFPS(60),
)
```

### Lifecycle Management

Dravya manages five application states:

- **Initializing**: Setting up modules and resources
- **Running**: Main event loop active
- **Paused**: Suspended but can resume
- **ShuttingDown**: Cleanup in progress
- **Terminated**: Application ended

### Main Loop

The FPS-capped render loop ensures smooth performance:

```go
app.SetRoot(myComponent)
if err := app.Run(context.Background()); err != nil {
    log.Fatal(err)
}
```

## Core Components

### 1. Lifecycle Hooks

Register callbacks for state transitions:

```go
app.OnReady(func() {
    // Called after initialization
    setupEventHandlers()
})
```

### 2. Event Hub Integration

Connects keyboard, mouse, and system events to your components:

```go
eventHub := app.EventHub()
eventHub.On(agni.EventTypeKey, handleKeyPress)
```

### 3. Focus Management

Built-in focus system for keyboard navigation:

```go
focusMgr := app.FocusManager()
focusMgr.Register("my-input-field")
```

### 4. Automatic Re-rendering

Observable state changes trigger UI updates automatically:

```go
counter := prana.NewObservable(0)
counter.Set(42)  // UI re-renders automatically
```

## Configuration Options

### WithLogLevel

Set logging verbosity:

```go
dravya.WithLogLevel(slog.LevelDebug)
```

### WithFPS

Target frame rate (default: 60):

```go
dravya.WithFPS(30)  // Lower for reduced CPU usage
```

### WithPprof

Enable profiling server:

```go
dravya.WithPprof(":6060")
```

## API Reference

### App Methods

#### `SetRoot(component Component)`

Sets the root UI component to render.

#### `Run(ctx context.Context) error`

Starts the application main loop. Blocks until shutdown.

#### `Shutdown(err error) error`

Gracefully shuts down the application.

#### `Logger() *Logger`

Returns the structured logger instance.

#### `EventHub() EventHub`

Returns the event hub for registering handlers.

#### `FocusManager() *FocusManager`

Returns the focus manager for keyboard navigation.

#### `Stats() AppStats`

Returns runtime statistics (uptime, frame count, etc.).

## Best Practices

### 1. Graceful Shutdown

Always use context for cancellation:

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// Shutdown on signal
go func() {
    <-sigChan
    cancel()
}()

app.Run(ctx)
```

### 2. Resource Management

Register cleanup in shutdown hooks:

```go
app.Lifecycle().OnState(dravya.StateShuttingDown, func(ctx context.Context) error {
    // Clean up resources
    return nil
})
```

### 3. Performance Monitoring

Use pprof for optimization:

```go
app := dravya.NewApp(dravya.WithPprof(":6060"))
// View: http://localhost:6060/debug/pprof/
```

## Examples

### Basic Application

```go
package main

import (
    "context"
    "github.com/TIVerse/drav/pkg/dravya"
    "github.com/TIVerse/drav/pkg/maya"
)

type HelloWorld struct{}

func (h *HelloWorld) Render(ctx maya.RenderContext) maya.View {
    return maya.Text("Hello, DRAV!")
}

func main() {
    app := dravya.NewApp()
    app.SetRoot(&HelloWorld{})
    app.Run(context.Background())
}
```

### With Event Handling

```go
app.OnReady(func() {
    eventHub := app.EventHub()
    eventHub.On(agni.EventTypeKey, func(ctx context.Context, event dravya.Event) error {
        if keyEvent, ok := event.(*agni.KeyEvent); ok {
            app.Logger().Info("Key pressed", "key", keyEvent.Key)
        }
        return nil
    })
})
```

## Related Modules

- **[Agni](agni.md)**: Event hub and dispatcher
- **[Maya](maya.md)**: Rendering and layout
- **[Prana](prana.md)**: Reactive state management

## Performance Considerations

- **Frame Budget**: 60 FPS = 16.67ms per frame
- **Render Optimization**: Only dirty regions are redrawn
- **Event Batching**: Multiple events processed per frame
- **Memory**: ~50MB baseline, scales with UI complexity

## Troubleshooting

### Application Hangs on Startup

Check if tcell can initialize the terminal:

```go
app := dravya.NewApp()
if err := app.Run(ctx); err != nil {
    log.Printf("Failed to start: %v", err)
}
```

### Events Not Working

Ensure handlers are registered in `OnReady`:

```go
app.OnReady(func() {
    // Register handlers here
})
```

### Memory Leaks

Use `RequestRender()` sparingly and unsubscribe from observables:

```go
unwatch := observable.Watch(handler)
defer unwatch()
```

## See Also

- [Getting Started Guide](../getting-started.md)
- [Concepts](../concepts.md)
- [Examples](../../examples/)
