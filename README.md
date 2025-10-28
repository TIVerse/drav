# DRAV (द्रव) — Dynamic Reactive Application View

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/TIVerse/drav)](https://goreportcard.com/report/github.com/TIVerse/drav)

A next-generation terminal UI framework for Go featuring reactive state management, first-class command engine, secure plugins, and modern rendering.

## 🌊 What is DRAV?

DRAV (Sanskrit: द्रव, meaning "fluid" or "dynamic") is a comprehensive framework for building sophisticated terminal user interfaces in Go. Unlike traditional TUI libraries that require manual state synchronization and event wiring, DRAV provides:

- **🔄 Reactive State**: Observable values that automatically trigger UI updates
- **⌨️ Command Engine**: Built-in command palette with completion, history, and undo/redo
- **🧩 Plugin System**: Hot-loadable, capability-secured plugins (WASM/out-of-process)
- **🎨 Modern Renderer**: Diff-based rendering with flexbox-inspired layouts
- **🚀 Performance**: 60 FPS targeting with <50MB baseline memory
- **🛡️ Security**: Capability-based plugin isolation and input sanitization

## Quick Start

### Installation

```bash
go get github.com/TIVerse/drav
```

### Hello DRAV (5 minutes)

```go
package main

import (
    "context"
    "github.com/TIVerse/drav/pkg/dravya"
    "github.com/TIVerse/drav/pkg/maya"
)

func main() {
    app := dravya.NewApp()
    
    // Set root component
    app.SetRoot(maya.Text("Hello, DRAV! 🌊"))
    
    // Run application
    if err := app.Run(context.Background()); err != nil {
        panic(err)
    }
}
```

### Reactive Counter (10 minutes)

```go
package main

import (
    "context"
    "fmt"
    "github.com/TIVerse/drav/pkg/dravya"
    "github.com/TIVerse/drav/pkg/prana"
    "github.com/TIVerse/drav/pkg/maya"
    "github.com/TIVerse/drav/pkg/agni"
)

type Counter struct {
    count *prana.Observable[int]
}

func (c *Counter) Render(ctx maya.RenderContext) maya.View {
    return maya.Column(
        maya.Text(fmt.Sprintf("Count: %d", c.count.Get())),
        maya.Text("Press [+] to increment, [-] to decrement"),
    )
}

func main() {
    app := dravya.NewApp()
    
    counter := &Counter{
        count: prana.NewObservable(0),
    }
    
    // Handle keyboard input
    // (event hub integration to be added)
    
    app.SetRoot(counter)
    
    if err := app.Run(context.Background()); err != nil {
        panic(err)
    }
}
```

## 📦 Architecture

DRAV consists of seven core modules, each named after Sanskrit concepts:

### Core Modules

- **Dravya (द्रव्य)**: Runtime core — lifecycle, concurrency, main loop
- **Agni (अग्नि)**: Event hub — dispatch, priorities, timers
- **Māyā (माया)**: Renderer — virtual UI, diff, layout engine
- **Prāṇa (प्राण)**: Reactive state — observables, stores, effects
- **Vāk (वाक्)**: Command engine — parse, execute, history, undo/redo
- **Vāyu (वायु)**: Plugin system — WASM/process loaders, capabilities
- **Śrī (श्री)**: Theme engine — palettes, styles, animations

### Design Principles

1. **Framework over Library**: DRAV owns the main loop and coordinates all subsystems
2. **Reactive by Default**: State changes automatically propagate to the UI
3. **Security First**: Capability-based isolation for plugins and untrusted code
4. **Performance Engineering**: Profiled, benchmarked, and optimized for 60 FPS
5. **Developer Experience**: Quick starts, debugging tools, comprehensive docs

## 🎯 Features

### Reactive State Management

```go
// Observable values
count := prana.NewObservable(0)
count.Watch(func(old, new int) {
    fmt.Printf("Count changed: %d -> %d\n", old, new)
})
count.Set(42)

// Computed values with dependency tracking
doubled := prana.ComputedFromObservable(count, func(n int) int {
    return n * 2
})

// Stores with actions and effects
type AppState struct {
    User string
    LoggedIn bool
}

store := prana.NewStore(AppState{})
store.RegisterReducer("LOGIN", func(state AppState, payload any) AppState {
    state.User = payload.(string)
    state.LoggedIn = true
    return state
})
```

### Event System

```go
// Priority-based event dispatch
dispatcher := agni.NewDispatcher(1000, 10)

// Register handlers
unsubscribe := dispatcher.On(agni.EventTypeKey, func(ctx context.Context, event agni.Event) error {
    keyEvent := event.(*agni.KeyEvent)
    fmt.Printf("Key pressed: %c\n", keyEvent.Rune)
    return nil
}, agni.WithPriority(agni.PriorityHigh))

// Emit events
dispatcher.Emit(ctx, agni.NewKeyEvent(agni.KeyRune, 'a', agni.ModNone, false))

// Timers
dispatcher.After(ctx, "timeout", 5*time.Second, func(ctx context.Context) {
    fmt.Println("Timeout!")
})
```

### Command Engine

```go
// Register commands
registry.Register(vak.Command{
    Name:    "theme",
    Summary: "Switch theme",
    Flags:   []vak.Flag{{Name: "name", Type: "string"}},
    Execute: func(ctx context.Context, args []string) (vak.Result, error) {
        // Switch theme logic
        return vak.Success("Theme changed"), nil
    },
    Undo: func(ctx context.Context) error {
        // Restore previous theme
        return nil
    },
})

// Execute with autocomplete
result, err := registry.Execute(ctx, "theme --name dark")
```

### Plugin System

```go
// Load a WASM plugin with capabilities
plugin, err := pluginMgr.Load("./plugins/myplugin.wasm", vayu.Capabilities{
    Filesystem: vayu.FSCapability{
        Read:  []string{"/tmp"},
        Write: []string{"/tmp/output"},
    },
    Network: vayu.NetworkCapability{
        AllowedDomains: []string{"api.example.com"},
    },
})
```

## 🧪 Testing

DRAV includes comprehensive testing infrastructure:

```bash
# Run tests
make test

# With coverage
make cover

# Benchmarks
make bench

# E2E tests
go test ./tests/e2e/...

# Fuzz tests
go test -fuzz=FuzzDiff ./tests/fuzz/
```

## 📊 Performance

- **60 FPS target**: ~16ms per frame budget
- **Baseline memory**: <50MB for typical applications
- **Startup time**: <500ms
- **Plugin load**: <200ms (WASM)

Performance is continuously monitored with:
- pprof integration
- Frame timing metrics
- Memory profiling
- CI performance gates

## 🔒 Security

- **Input sanitization**: ANSI escape sequence filtering
- **Path traversal protection**: Safe filesystem operations
- **Plugin isolation**: Capability tokens, resource limits
- **Dependency scanning**: Automated vulnerability checks

See [SECURITY.md](SECURITY.md) for details.

## 📚 Documentation

- [Getting Started](docs/src/getting-started.md)
- [Core Concepts](docs/src/concepts.md)
- [API Reference](docs/src/api/reference.md)
- [Module Guides](docs/src/modules/)
- [Testing Guide](docs/src/guides/testing.md)
- [Security Model](docs/src/guides/security.md)

## 🗺️ Roadmap

- **v0.1**: Foundation (Dravya, Agni, basic Māyā)
- **v0.2**: Reactivity (Prāṇa, auto re-render)
- **v0.3**: Commands (Vāk, palette basics)
- **v0.4**: Plugins (Vāyu with WASM)
- **v0.5**: Polish (Śrī themes/animations) → Beta
- **v1.0**: Stable APIs, production-ready

See [ROADMAP.md](ROADMAP.md) for detailed milestones.

## 🤝 Contributing

We welcome contributions! Please see:

- [CONTRIBUTING.md](CONTRIBUTING.md) — Development guidelines
- [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) — Community standards
- [GitHub Issues](https://github.com/TIVerse/drav/issues) — Bug reports and features

## 📄 License

MIT License — see [LICENSE](LICENSE) for details.

## 👨‍💻 Author

Created by [Abhineesh Priyam](https://github.com/abhineeshpriyam) ([@abhineeshpriyam](https://github.com/abhineeshpriyam))

Part of the [TIVerse](https://github.com/TIVerse) organization.

## 🙏 Acknowledgments

DRAV is inspired by:
- **Elm Architecture** — Unidirectional data flow
- **React** — Virtual DOM and diffing
- **Bubbletea** — Go TUI patterns
- **Textual (Python)** — Modern TUI framework design

---

**Built with 💙 by [Abhineesh Priyam](https://github.com/abhineeshpriyam)**
