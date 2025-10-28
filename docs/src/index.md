# DRAV (द्रव) - Dynamic Reactive Application View

Welcome to DRAV, a next-generation terminal UI framework for Go that brings modern development patterns to the terminal.

## What is DRAV?

DRAV (Sanskrit: द्रव, meaning "fluid" or "dynamic") is a comprehensive framework for building sophisticated terminal user interfaces in Go. Unlike traditional TUI libraries, DRAV provides:

🔄 **Reactive State Management**  
Observable values and stores that automatically trigger UI updates

⌨️ **First-Class Commands**  
Built-in command palette with completion, history, and undo/redo

🧩 **Secure Plugin System**  
Hot-loadable plugins with capability-based security

🎨 **Modern Rendering**  
Diff-based rendering with flexbox-inspired layouts

🚀 **Performance Engineered**  
60 FPS targeting with <50MB baseline memory

🛡️ **Security First**  
Input sanitization, plugin isolation, and secure-by-default APIs

## Quick Example

```go
package main

import (
    "context"
    "fmt"
    "github.com/TIVerse/drav/pkg/dravya"
    "github.com/TIVerse/drav/pkg/maya"
    "github.com/TIVerse/drav/pkg/prana"
)

type Counter struct {
    count *prana.Observable[int]
}

func (c *Counter) Render(ctx maya.RenderContext) maya.View {
    return maya.Column(
        maya.Text(fmt.Sprintf("Count: %d", c.count.Get())),
        maya.Text("Press [+] to increment"),
    )
}

func main() {
    app := dravya.NewApp()
    counter := &Counter{count: prana.NewObservable(0)}
    
    app.SetRoot(counter)
    app.Run(context.Background())
}
```

## Why DRAV?

### Problem

Existing TUI frameworks require manual state synchronization, lack unified command systems, and don't provide secure plugin architectures. Building complex terminal applications feels like building web apps from 2005.

### Solution

DRAV brings modern patterns from web development to the terminal:

- **Reactive by Default**: State changes automatically propagate to the UI, just like React
- **Framework over Library**: DRAV owns the main loop and coordinates all subsystems
- **Security Built-In**: Capability-based isolation for plugins and untrusted code
- **Developer Experience**: Quick starts, debugging tools, comprehensive docs

## Architecture

DRAV consists of seven core modules, each named after Sanskrit concepts:

- **Dravya (द्रव्य)**: Runtime core — lifecycle, concurrency, main loop
- **Agni (अग्नि)**: Event hub — dispatch, priorities, timers
- **Māyā (माया)**: Renderer — virtual UI, diff, layout engine
- **Prāṇa (प्राण)**: Reactive state — observables, stores, effects
- **Vāk (वाक्)**: Command engine — parse, execute, history, undo/redo
- **Vāyu (वायु)**: Plugin system — WASM/process loaders, capabilities
- **Śrī (श्री)**: Theme engine — palettes, styles, animations

## Getting Started

Choose your path:

<div class="grid cards" markdown>

-   :material-clock-fast:{ .lg .middle } __5-Minute Quick Start__

    ---

    Get up and running with DRAV in just 5 minutes

    [:octicons-arrow-right-24: Quick Start](getting-started.md)

-   :material-book-open-variant:{ .lg .middle } __Learn the Concepts__

    ---

    Understand reactive state, events, and component composition

    [:octicons-arrow-right-24: Core Concepts](concepts.md)

-   :material-code-braces:{ .lg .middle } __API Reference__

    ---

    Complete API documentation for all modules

    [:octicons-arrow-right-24: API Docs](api/reference.md)

-   :material-lightbulb:{ .lg .middle } __Examples__

    ---

    Browse example applications demonstrating DRAV features

    [:octicons-arrow-right-24: Examples](https://github.com/TIVerse/drav/tree/main/examples)

</div>

## Features

### Reactive State

```go
// Observable values
count := prana.NewObservable(0)
count.Watch(func(old, new int) {
    fmt.Printf("Count: %d → %d\n", old, new)
})
count.Set(42) // Triggers watchers

// Computed values
doubled := prana.ComputedFromObservable(count, func(n int) int {
    return n * 2
})
```

### Event System

```go
// Priority-based dispatch
dispatcher.On(agni.EventTypeKey, handler, 
    agni.WithPriority(agni.PriorityHigh))

// Timers
dispatcher.After(ctx, "timeout", 5*time.Second, func(ctx context.Context) {
    fmt.Println("Timeout!")
})
```

### Command Engine

```go
registry.Register(vak.Command{
    Name:    "theme",
    Summary: "Switch theme",
    Execute: func(ctx context.Context, args []string) (vak.Result, error) {
        return vak.SuccessResult("Theme changed"), nil
    },
    Undo: func(ctx context.Context) error {
        return nil
    },
})
```

## Community

- **GitHub**: [TIVerse/drav](https://github.com/TIVerse/drav)
- **Discussions**: [GitHub Discussions](https://github.com/TIVerse/drav/discussions)
- **Issues**: [Bug Reports & Features](https://github.com/TIVerse/drav/issues)
- **Twitter**: [@TIVerse](https://twitter.com/TIVerse)

## Contributing

We welcome contributions! See our [Contributing Guide](https://github.com/TIVerse/drav/blob/main/CONTRIBUTING.md) to get started.

## License

DRAV is released under the [MIT License](https://github.com/TIVerse/drav/blob/main/LICENSE).

---

**Built with 💙 by the TIVerse team**
