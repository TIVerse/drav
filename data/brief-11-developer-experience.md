# Section 11: Developer Experience

[← Back to Index](brief-index.md) | [Previous: Testing Strategy](brief-10-testing-strategy.md) | [Next: Use Cases →](brief-12-use-cases.md)

---

## 11.1 Getting Started

### Installation

```bash
# Install via go get
go get github.com/abhineesh/drav

# Or install CLI tool
go install github.com/abhineesh/drav/cmd/drav@latest
```

### First Application (5 minutes)

```go
package main

import "github.com/abhineesh/drav"

func main() {
    app := drav.NewApp()
    
    app.SetRoot(drav.Text("Hello, DRAV!"))
    
    app.Run()
}
```

**Output**: Terminal displays "Hello, DRAV!"

### Interactive Counter (10 minutes)

```go
package main

import (
    "fmt"
    "github.com/abhineesh/drav"
)

func main() {
    app := drav.NewApp()
    
    counter := drav.NewObservable(0)
    
    app.SetRoot(drav.Column(
        drav.Text(counter.Format("Count: %d")),
        drav.Button("Increment", func() {
            counter.Set(counter.Get() + 1)
        }),
        drav.Button("Decrement", func() {
            counter.Set(counter.Get() - 1)
        }),
    ))
    
    app.Run()
}
```

### Dashboard (30 minutes)

```go
package main

import "github.com/abhineesh/drav"

func main() {
    app := drav.NewApp(drav.WithTheme(drav.DarkTheme))
    
    cpuUsage := drav.NewObservable(0.0)
    memUsage := drav.NewObservable(0.0)
    
    // Update metrics every second
    app.OnTick(time.Second, func() {
        cpuUsage.Set(getCPUUsage())
        memUsage.Set(getMemUsage())
    })
    
    app.SetRoot(drav.Row(
        drav.Panel("CPU", drav.ProgressBar(cpuUsage)),
        drav.Panel("Memory", drav.ProgressBar(memUsage)),
    ))
    
    app.Run()
}
```

---

## 11.2 Documentation Strategy

### Documentation Tiers

#### Tier 1: Getting Started (5-30 minutes)
- Installation guide
- Hello World tutorial
- Core concepts overview
- Interactive examples

#### Tier 2: Guides (1-2 hours each)
- Building Components
- State Management
- Command Development
- Plugin Creation
- Theming & Styling
- Performance Optimization

#### Tier 3: API Reference (Reference material)
- Complete API documentation
- Every public type, function, method
- Code examples for each API
- Auto-generated from godoc

#### Tier 4: Advanced Topics (2-4 hours each)
- Architecture deep dives
- Performance tuning
- Security best practices
- Contributing guide

### Documentation Site Structure

```
docs/
├── index.md                    # Landing page
├── getting-started/
│   ├── installation.md
│   ├── first-app.md
│   ├── concepts.md
│   └── examples.md
├── guides/
│   ├── components.md
│   ├── state.md
│   ├── commands.md
│   ├── plugins.md
│   ├── themes.md
│   └── performance.md
├── api/
│   ├── app.md
│   ├── component.md
│   ├── view.md
│   ├── state.md
│   ├── command.md
│   └── ...
├── advanced/
│   ├── architecture.md
│   ├── internals.md
│   ├── security.md
│   └── contributing.md
└── examples/
    ├── counter.md
    ├── dashboard.md
    ├── editor.md
    └── ...
```

---

## 11.3 Error Messages

### Principles

1. **Clear**: Explain what went wrong
2. **Actionable**: Suggest how to fix it
3. **Contextual**: Show where the error occurred
4. **Helpful**: Link to documentation

### Examples

#### Bad Error Messages
```
Error: invalid input
Error: failed
Error: nil pointer
```

#### Good Error Messages
```
Error: Observable value cannot be nil
  → Hint: Initialize with NewObservable(defaultValue)
  → Docs: https://drav.dev/docs/state

Error: Command "save" not found
  → Available commands: load, quit, help
  → Hint: Use :help to list all commands
  → Docs: https://drav.dev/docs/commands

Error: Plugin "git-integration" failed to load
  → Reason: Missing dependency "libgit2"
  → Fix: Install with: sudo apt install libgit2-dev
  → Docs: https://drav.dev/docs/plugins/git-integration
```

### Error Helper Functions

```go
func NewError(code ErrorCode, message string) *DravError {
    return &DravError{
        Code:    code,
        Message: message,
        Hint:    errorHints[code],
        DocsURL: fmt.Sprintf("https://drav.dev/errors/%s", code),
    }
}

func (e *DravError) Error() string {
    var sb strings.Builder
    
    sb.WriteString(fmt.Sprintf("Error: %s\n", e.Message))
    
    if e.Hint != "" {
        sb.WriteString(fmt.Sprintf("  → Hint: %s\n", e.Hint))
    }
    
    if e.DocsURL != "" {
        sb.WriteString(fmt.Sprintf("  → Docs: %s\n", e.DocsURL))
    }
    
    return sb.String()
}
```

---

## 11.4 IDE Integration

### VS Code Extension (Future)

**Features**:
- Syntax highlighting for DRAV DSL
- Component snippets
- API autocomplete
- Live preview of TUI
- Debugging support

**Snippets**:
```json
{
  "DRAV Component": {
    "prefix": "drav-component",
    "body": [
      "type ${1:ComponentName} struct {",
      "\t${2:state} drav.Observable[${3:type}]",
      "}",
      "",
      "func (c *${1:ComponentName}) Render() drav.View {",
      "\treturn drav.${4:Text}(${5:content})",
      "}"
    ]
  }
}
```

### gopls Integration

**Features**:
- Hover documentation for DRAV APIs
- Go to definition
- Find references
- Rename refactoring

---

## 11.5 Debugging Tools

### Built-in Debugger

```go
// Enable debug mode
app := drav.NewApp(drav.WithDebug(true))

// Debug overlay shows:
// - Current FPS
// - Memory usage
// - Goroutine count
// - Event queue size
// - Render time
```

### Component Inspector

```go
// Press F12 to open inspector
app.EnableInspector()

// Inspector shows:
// - Component tree
// - State values
// - Event handlers
// - Performance metrics
```

### Logging

```go
// Structured logging
drav.Log(drav.Debug, "Rendering component",
    "component", component.Name(),
    "duration", duration,
)

// Log levels:
// Trace, Debug, Info, Warn, Error, Fatal

// Configure log output
app.ConfigureLogger(drav.LoggerConfig{
    Level:  drav.Debug,
    Format: drav.JSON,  // or drav.Text
    Output: os.Stderr,
})
```

---

## 11.6 Development Workflow

### Hot Reload (Future)

```bash
# Watch mode - recompiles and restarts on change
drav dev ./cmd/myapp

# Output:
# Watching for changes...
# [12:34:56] Building...
# [12:34:57] Running...
# [12:35:10] Change detected: component.go
# [12:35:11] Restarting...
```

### Testing Workflow

```bash
# Run tests in watch mode
drav test --watch

# Run specific test
drav test TestCounter

# Run with coverage
drav test --coverage

# Run benchmarks
drav bench
```

---

## 11.7 Common Patterns

### Pattern: Loading State

```go
type DataComponent struct {
    loading drav.Observable[bool]
    data    drav.Observable[*Data]
    error   drav.Observable[error]
}

func (c *DataComponent) Render() drav.View {
    if c.loading.Get() {
        return drav.Spinner("Loading...")
    }
    
    if err := c.error.Get(); err != nil {
        return drav.Text(fmt.Sprintf("Error: %v", err))
    }
    
    return drav.Text(c.data.Get().String())
}

func (c *DataComponent) FetchData() {
    c.loading.Set(true)
    c.error.Set(nil)
    
    go func() {
        data, err := fetchFromAPI()
        c.loading.Set(false)
        
        if err != nil {
            c.error.Set(err)
        } else {
            c.data.Set(data)
        }
    }()
}
```

### Pattern: Form Validation

```go
type FormComponent struct {
    name  drav.Observable[string]
    email drav.Observable[string]
    valid drav.Observable[bool]
}

func (c *FormComponent) Render() drav.View {
    // Validate on change
    c.valid.Set(c.isValid())
    
    return drav.Column(
        drav.Input(c.name, "Name"),
        drav.Input(c.email, "Email"),
        drav.Button("Submit", c.submit).
            Enabled(c.valid.Get()),
    )
}

func (c *FormComponent) isValid() bool {
    return c.name.Get() != "" && 
           strings.Contains(c.email.Get(), "@")
}
```

### Pattern: List with Selection

```go
type ListComponent struct {
    items    []string
    selected drav.Observable[int]
}

func (c *ListComponent) Render() drav.View {
    return drav.List(c.items, c.selected).
        OnKey(drav.KeyUp, func() {
            sel := c.selected.Get()
            if sel > 0 {
                c.selected.Set(sel - 1)
            }
        }).
        OnKey(drav.KeyDown, func() {
            sel := c.selected.Get()
            if sel < len(c.items)-1 {
                c.selected.Set(sel + 1)
            }
        })
}
```

---

## 11.8 Learning Resources

### Official Resources

1. **Documentation Site**: https://drav.dev
2. **API Reference**: https://pkg.go.dev/github.com/abhineesh/drav
3. **Examples Repository**: https://github.com/abhineesh/drav-examples
4. **Video Tutorials**: YouTube channel (planned)
5. **Interactive Playground**: https://play.drav.dev (planned)

### Community Resources

1. **Discord Server**: Community chat
2. **GitHub Discussions**: Q&A, ideas, showcases
3. **Stack Overflow**: Tag `drav`
4. **Blog Posts**: Community tutorials
5. **Conference Talks**: GopherCon, etc.

### Learning Path

#### Beginner (1-2 weeks)
- Complete getting started guide
- Build 3-5 simple apps
- Read core concepts
- Join Discord

#### Intermediate (1-2 months)
- Build medium-sized app
- Learn state management patterns
- Create custom commands
- Contribute to examples

#### Advanced (3-6 months)
- Build complex app with plugins
- Optimize performance
- Contribute to core
- Write plugin

---

## 11.9 Migration Guides

### From BubbleTea

```go
// BubbleTea
type model struct {
    counter int
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "+" {
            m.counter++
        }
    }
    return m, nil
}

func (m model) View() string {
    return fmt.Sprintf("Counter: %d", m.counter)
}

// DRAV equivalent
type Counter struct {
    count drav.Observable[int]
}

func (c *Counter) Render() drav.View {
    return drav.Column(
        drav.Text(c.count.Format("Counter: %d")),
    ).OnKey(drav.Key('+'), func() {
        c.count.Set(c.count.Get() + 1)
    })
}
```

### From tview

```go
// tview
box := tview.NewBox().
    SetBorder(true).
    SetTitle("Hello")
app.SetRoot(box, true).Run()

// DRAV equivalent
app := drav.NewApp()
app.SetRoot(drav.Panel("Hello", drav.Text("Content")))
app.Run()
```

---

## 11.10 Troubleshooting

### Common Issues

#### Issue: Observable not triggering re-render

**Problem**:
```go
type Component struct {
    data []int  // Not observable!
}
```

**Solution**:
```go
type Component struct {
    data drav.Observable[[]int]  // Observable
}
```

#### Issue: Memory leak

**Problem**:
```go
// Forgotten subscription
obs.Watch(func(v int) {
    // Never unsubscribed
})
```

**Solution**:
```go
// Save unsubscribe function
unsub := obs.Watch(func(v int) {
    // ...
})

// Call on cleanup
defer unsub()
```

#### Issue: Slow rendering

**Check**:
1. Enable profiler: `drav --profile`
2. Check component count
3. Verify no expensive operations in Render()
4. Use caching for computed values

---

## Summary

**Developer Experience Goals**:
- **5-minute setup**: From zero to first app
- **Excellent documentation**: Multiple tiers, comprehensive
- **Helpful errors**: Clear, actionable messages
- **Great tooling**: IDE support, debugging, hot reload
- **Active community**: Discord, forums, examples

**Resources**:
- Comprehensive docs site
- Interactive tutorials
- Example applications
- Migration guides
- Troubleshooting help

**Philosophy**: Make the easy things easy, and the hard things possible.

---

[← Back to Index](brief-index.md) | [Previous: Testing Strategy](brief-10-testing-strategy.md) | [Next: Use Cases →](brief-12-use-cases.md)
