# API Reference

Complete API documentation for DRAV modules.

## Module Overview

DRAV consists of seven core modules:

| Module | Package | Description |
|--------|---------|-------------|
| Dravya | `pkg/dravya` | Runtime core, application lifecycle |
| Agni | `pkg/agni` | Event system, dispatch, timers |
| Māyā | `pkg/maya` | Renderer, layouts, components |
| Prāṇa | `pkg/prana` | Reactive state, observables, stores |
| Vāk | `pkg/vak` | Command engine, parser, history |
| Vāyu | `pkg/vayu` | Plugin system, capabilities, loaders |
| Śrī | `pkg/sri` | Themes, styles, animations |

## Quick Links

- [Dravya API](../modules/dravya.md) - Application and lifecycle
- [Agni API](../modules/agni.md) - Events and dispatch
- [Māyā API](../modules/maya.md) - Rendering and layout
- [Prāṇa API](../modules/prana.md) - Reactive state
- [Vāk API](../modules/vak.md) - Commands
- [Vāyu API](../modules/vayu.md) - Plugins
- [Śrī API](../modules/sri.md) - Themes

## Common Patterns

### Creating an Application

```go
import "github.com/TIVerse/drav/pkg/dravya"

app := dravya.NewApp(
    dravya.WithLogLevel(slog.LevelInfo),
    dravya.WithFPSCap(60),
    dravya.WithTheme("default-dark"),
)
```

### Reactive State

```go
import "github.com/TIVerse/drav/pkg/prana"

// Observable
value := prana.NewObservable(42)

// Store
store := prana.NewStore(MyState{})
store.RegisterReducer("ACTION", reducer)
store.Dispatch(ctx, action)
```

### Components

```go
import "github.com/TIVerse/drav/pkg/maya"

type MyComponent struct {}

func (m *MyComponent) Render(ctx maya.RenderContext) maya.View {
    return maya.Text("Hello!")
}
```

### Events

```go
import "github.com/TIVerse/drav/pkg/agni"

dispatcher := agni.NewDispatcher(1000, 10)
dispatcher.On(agni.EventTypeKey, handler)
```

### Commands

```go
import "github.com/TIVerse/drav/pkg/vak"

registry := vak.NewRegistry()
registry.Register(command)
result, err := registry.Execute(ctx, "command args")
```

## Type Reference

### Application Types

```go
// App is the main DRAV application
type App struct { /* ... */ }

// AppOption configures the application
type AppOption func(*appConfig)

// LifecycleState represents app state
type LifecycleState int
```

### Component Types

```go
// Component renders UI
type Component interface {
    Render(ctx RenderContext) View
}

// View is a virtual UI node
type View interface {
    Type() string
    Children() []View
    Attrs() Attributes
}
```

### State Types

```go
// Observable is a reactive value
type Observable[T any] struct { /* ... */ }

// Computed is a derived value
type Computed[T any] struct { /* ... */ }

// Store manages app state
type Store[S any] struct { /* ... */ }
```

### Event Types

```go
// Event is the base event interface
type Event interface {
    Type() string
    Time() time.Time
}

// KeyEvent represents keyboard input
type KeyEvent struct {
    Key       Key
    Rune      rune
    Modifiers ModMask
    Repeat    bool
}

// MouseEvent represents mouse input
type MouseEvent struct {
    X, Y      int
    Button    MouseButton
    Action    MouseAction
}
```

### Command Types

```go
// Command is an executable command
type Command struct {
    Name     string
    Summary  string
    Execute  ExecuteFunc
    Undo     UndoFunc
    Complete CompleteFunc
}

// Result is command execution result
type Result interface {
    Success() bool
    Message() string
    Data() any
}
```

### Plugin Types

```go
// Plugin interface
type Plugin interface {
    Metadata() PluginMetadata
    Init(ctx context.Context) error
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Capabilities() Capabilities
}

// Capabilities define plugin permissions
type Capabilities struct {
    Filesystem FSCapability
    Network    NetworkCapability
    State      StateCapability
    UI         UICapability
}
```

### Theme Types

```go
// Theme represents a UI theme
type Theme struct {
    Name    string
    Palette Palette
    Styles  map[string]Style
}

// Palette defines colors
type Palette struct {
    Primary, Secondary Color
    Success, Warning, Error Color
    Background, Surface, Text Color
}

// Style defines text styling
type Style struct {
    Foreground, Background Color
    Bold, Italic, Underline bool
}
```

## Error Handling

DRAV uses standard Go error handling:

```go
// Check errors explicitly
result, err := registry.Execute(ctx, input)
if err != nil {
    return fmt.Errorf("execution failed: %w", err)
}

// Context cancellation
if ctx.Err() == context.Canceled {
    return ctx.Err()
}
```

## Concurrency

DRAV is thread-safe where documented:

```go
// Safe: Observable operations
observable.Set(value)  // Thread-safe

// Safe: Event dispatch
dispatcher.Emit(ctx, event)  // Thread-safe

// Safe: Store operations
store.Dispatch(ctx, action)  // Thread-safe

// Unsafe: Direct component access
// Don't access component fields from multiple goroutines
```

## Examples

See [examples directory](https://github.com/TIVerse/drav/tree/main/examples) for complete working examples.

## Generated Documentation

For auto-generated API docs from source code:

```bash
godoc -http=:6060
```

Then visit http://localhost:6060/pkg/github.com/TIVerse/drav/

## Version Compatibility

This documentation is for DRAV v0.1+. See [CHANGELOG](https://github.com/TIVerse/drav/blob/main/CHANGELOG.md) for version-specific changes.
