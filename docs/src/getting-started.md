# Getting Started

Get up and running with DRAV in minutes!

## Installation

### Prerequisites

- Go 1.22 or later
- A terminal that supports ANSI escape sequences
- Windows Terminal (on Windows), iTerm2/Terminal.app (on macOS), or any modern Linux terminal

### Install via Go

```bash
go get github.com/TIVerse/drav
```

### Verify Installation

```bash
go version
# Should output: go version go1.22 or later
```

## Your First DRAV Application

Let's build a simple "Hello DRAV" application.

### Step 1: Create a New Project

```bash
mkdir hello-drav
cd hello-drav
go mod init hello-drav
```

### Step 2: Install DRAV

```bash
go get github.com/TIVerse/drav
```

### Step 3: Write the Code

Create `main.go`:

```go
package main

import (
    "context"
    "log"

    "github.com/TIVerse/drav/pkg/dravya"
    "github.com/TIVerse/drav/pkg/maya"
)

func main() {
    // Create the application
    app := dravya.NewApp()

    // Create a simple text view
    root := maya.Text("Hello, DRAV! 🌊\n\nPress Ctrl+C to exit.")

    // Set the root component
    app.SetRoot(maya.Stateless(root))

    // Run the application
    if err := app.Run(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

### Step 4: Run the Application

```bash
go run main.go
```

You should see "Hello, DRAV! 🌊" in your terminal!

## Building a Reactive Counter

Now let's build something more interesting with reactive state.

### The Counter Component

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/TIVerse/drav/pkg/dravya"
    "github.com/TIVerse/drav/pkg/maya"
    "github.com/TIVerse/drav/pkg/prana"
)

type Counter struct {
    count *prana.Observable[int]
}

func NewCounter() *Counter {
    return &Counter{
        count: prana.NewObservable(0),
    }
}

func (c *Counter) Render(ctx maya.RenderContext) maya.View {
    return maya.Column(
        maya.Text("=== Reactive Counter ==="),
        maya.Text(""),
        maya.Text(fmt.Sprintf("Count: %d", c.count.Get())),
        maya.Text(""),
        maya.Text("Press [+] to increment"),
        maya.Text("Press [-] to decrement"),
        maya.Text("Press [r] to reset"),
        maya.Text("Press [Ctrl+C] to exit"),
    )
}

func (c *Counter) Increment() {
    c.count.Update(func(n int) int { return n + 1 })
}

func (c *Counter) Decrement() {
    c.count.Update(func(n int) int { return n - 1 })
}

func (c *Counter) Reset() {
    c.count.Set(0)
}

func main() {
    app := dravya.NewApp()
    counter := NewCounter()

    // Watch for changes (optional)
    counter.count.Watch(func(old, new int) {
        log.Printf("Count changed: %d → %d", old, new)
    })

    app.SetRoot(counter)

    if err := app.Run(context.Background()); err != nil {
        log.Fatal(err)
    }
}
```

## Key Concepts

### 1. Application Lifecycle

```go
app := dravya.NewApp(
    dravya.WithLogLevel(slog.LevelInfo),
    dravya.WithFPSCap(60),
    dravya.WithMetrics(true),
)

// Set root component
app.SetRoot(myComponent)

// Run application
ctx := context.Background()
if err := app.Run(ctx); err != nil {
    log.Fatal(err)
}
```

### 2. Components

Components implement the `Component` interface:

```go
type Component interface {
    Render(ctx RenderContext) View
}
```

Create components as structs:

```go
type MyComponent struct {
    state *prana.Observable[string]
}

func (m *MyComponent) Render(ctx maya.RenderContext) maya.View {
    return maya.Text(m.state.Get())
}
```

### 3. Reactive State

Use observables for reactive values:

```go
// Create observable
name := prana.NewObservable("Alice")

// Get value
fmt.Println(name.Get()) // "Alice"

// Set value (triggers watchers)
name.Set("Bob")

// Watch for changes
unwatch := name.Watch(func(old, new string) {
    fmt.Printf("%s → %s\n", old, new)
})
defer unwatch()
```

### 4. Layout System

Compose UIs with layout primitives:

```go
view := maya.Column(
    maya.Text("Header"),
    maya.Row(
        maya.Text("Left"),
        maya.Text("Right"),
    ),
    maya.Text("Footer"),
)
```

## Next Steps

Now that you've built your first DRAV application, explore:

- [Core Concepts](concepts.md) - Deep dive into reactive state, events, and layouts
- [Module Guides](modules/dravya.md) - Learn about each DRAV module
- [Examples](https://github.com/TIVerse/drav/tree/main/examples) - Browse complete example applications
- [API Reference](api/reference.md) - Complete API documentation

## Troubleshooting

### Application doesn't start

Make sure you're using Go 1.22 or later:
```bash
go version
```

### Colors not showing

Check your terminal supports ANSI colors. On Windows, use Windows Terminal instead of cmd.exe.

### Import errors

Ensure dependencies are installed:
```bash
go mod tidy
```

## Getting Help

- **Documentation**: You're reading it!
- **GitHub Issues**: [Report bugs or request features](https://github.com/TIVerse/drav/issues)
- **Discussions**: [Ask questions](https://github.com/TIVerse/drav/discussions)
- **Examples**: [Browse working code](https://github.com/TIVerse/drav/tree/main/examples)

---

Ready to build something amazing? Let's go! 🌊
