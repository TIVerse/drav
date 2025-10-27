# Section 5: Core Modules — Technical Deep Dive

[← Back to Index](brief-index.md) | [Previous: System Architecture](brief-04-system-architecture.md) | [Next: API Specifications →](brief-06-api-specifications.md)

---

This document provides detailed technical specifications for each of DRAV's seven core modules. Each module section includes purpose, architecture, implementation details, and API examples.

## 5.1 Module 1: Renderer (Māyā)

### Purpose & Philosophy

**Māyā** (माया) transforms abstract state into visual terminal output — the projection of reality into perceivable form.

### Core Components

**Virtual Buffer**:
```go
type Buffer struct {
    cells   [][]Cell
    width   int
    height  int
}

type Cell struct {
    Rune  rune
    Style Style
}
```

**Diff Algorithm**: Two-pass rendering
1. Line hashing (O(h)) — fast comparison
2. Cell-level diff for changed lines only
3. Minimal terminal updates

**Layout Engine**: Flexbox-inspired constraint-based layout supporting Row, Column, Stack, and Flex containers.

### Key Features
- Diff-based rendering (only draw changes)
- Double buffering (flicker-free)
- Terminal capability detection
- Style optimization (minimize escape sequences)
- Layout caching

---

## 5.2 Module 2: Command Engine (Vāk)

### Purpose & Philosophy

**Vāk** (वाक्) — Commands as speech acts that transform reality. Every action is accessible via command.

### Architecture

```go
type Command struct {
    Name        string
    Description string
    Handler     CommandHandler
    Completer   Completer
    Validator   Validator
}

type CommandRegistry struct {
    commands map[string]*Command
    history  *CommandHistory
    aliases  map[string]string
}
```

### Key Features
- Command registration with metadata
- Tab completion engine
- Command history (up/down arrows)
- Undo/redo support
- Flag parsing (--flag=value)
- Alias support
- Command palette (fuzzy search)

### Parser Features
- Quoted string support
- Escape sequence handling
- Long and short flags
- Positional arguments

---

## 5.3 Module 3: Event Hub (Agni)

### Purpose & Philosophy

**Agni** (अग्नि) — The messenger fire carrying events between components.

### Event Priority System

```go
type EventHub struct {
    critical chan Event  // Signals, Resize
    high     chan Event  // Keyboard, Mouse
    normal   chan Event  // Ticks, Custom
    low      chan Event  // Background
}
```

### Event Types
- **InputEvent**: Key, Mouse, Paste
- **SystemEvent**: Resize, Focus, Signal
- **TimerEvent**: Tick, Timeout
- **CustomEvent**: User-defined

### Features
- Priority-based dispatch
- Non-blocking handlers
- Timer management
- Event filtering
- Subscription management

---

## 5.4 Module 4: Reactive State (Prāṇa)

### Purpose & Philosophy

**Prāṇa** (प्राण) — State as life force. When state changes, the UI automatically "breathes" (updates).

### Observable Pattern

```go
type Observable[T any] struct {
    value     T
    observers []Observer[T]
}

func (o *Observable[T]) Set(value T) {
    o.value = value
    // Automatically notify all observers
    for _, obs := range o.observers {
        obs(value)
    }
}
```

### Key Features
- Generic observables for any type
- Automatic observer notification
- Derived state (computed observables)
- State transactions (batch updates)
- Dependency tracking
- Memory leak prevention

### Patterns
- **Local State**: Component-specific observables
- **Global State**: Application-wide stores
- **Derived State**: Computed from other observables
- **Async State**: Loading, error, data patterns

---

## 5.5 Module 5: Plugin System (Vāyu)

### Purpose & Philosophy

**Vāyu** (वायु) — Like wind, plugins permeate the system, extending it without being confined.

### Plugin Interface

```go
type Plugin interface {
    Name() string
    Version() string
    Init(runtime Runtime) error
    RegisterCommands(registry CommandRegistry)
    RegisterWidgets(registry WidgetRegistry)
    RegisterThemes(registry ThemeRegistry)
}
```

### Loading Mechanisms
1. **Go Plugins** (.so/.dll): Native Go plugin system
2. **WebAssembly** (.wasm): Sandboxed execution
3. **Embedded**: Compiled into binary

### Features
- Hot reload support
- Dependency resolution
- Version compatibility checks
- Sandboxing (resource limits)
- Plugin marketplace integration

### Security
- Capability-based access control
- Resource quotas (CPU, memory, files)
- API isolation
- Signature verification

---

## 5.6 Module 6: Theme Engine (Śrī)

### Purpose & Philosophy

**Śrī** (श्री) — Beauty and aesthetic coherence. Visual delight is not superficial but essential.

### Theme Structure

```go
type Theme struct {
    Name       string
    Colors     ColorPalette
    Styles     StyleMap
    Animations AnimationPresets
}

type ColorPalette struct {
    Background Color
    Foreground Color
    Primary    Color
    Secondary  Color
    Accent     Color
    // ... semantic colors
}
```

### Features
- Semantic color system
- Style inheritance
- Dark/light mode support
- Animation presets
- Gradient support
- Terminal capability adaptation

### Animation System
- Easing curves (linear, ease-in, ease-out, etc.)
- Transitions (fade, slide, scale)
- Duration control
- Frame interpolation

---

## 5.7 Module 7: Runtime Core (Dravya)

### Purpose & Philosophy

**Dravya** (द्रव्य) — The substance, the foundation upon which all modules exist.

### Responsibilities

```go
type Runtime struct {
    lifecycle  *Lifecycle
    resources  *ResourceManager
    errors     *ErrorRecovery
    metrics    *Metrics
}
```

### Lifecycle Management
- Initialization sequence
- Graceful shutdown
- Resource cleanup
- State persistence

### Resource Management
- Memory pools
- Goroutine tracking
- File handle management
- Connection pooling

### Error Recovery
- Panic recovery
- Graceful degradation
- Error boundaries
- Automatic retry logic

### Concurrency Coordination
- Goroutine lifecycle management
- Context propagation
- Deadlock detection
- Race condition prevention

---

## 5.8 Module Interaction Matrix

| Module | Depends On | Used By | Key Interaction |
|--------|-----------|---------|-----------------|
| **Māyā** | Prāṇa, Śrī, Agni | Application | Observes state, applies theme |
| **Vāk** | Agni, Prāṇa | Application, Plugins | Receives events, updates state |
| **Agni** | Dravya | All modules | Distributes events |
| **Prāṇa** | Dravya | Māyā, Vāk, Application | Notifies observers |
| **Vāyu** | Dravya | Application | Loads/manages plugins |
| **Śrī** | Dravya | Māyā | Provides styling |
| **Dravya** | None | All modules | Foundation |

---

## 5.9 Performance Characteristics

| Module | Critical Path | Typical Overhead | Optimization Strategy |
|--------|--------------|------------------|----------------------|
| **Māyā** | Yes (every frame) | 2-5ms | Diff algorithm, caching |
| **Vāk** | No (on command) | 1-2ms | Parser optimization |
| **Agni** | Yes (all events) | <1ms | Priority queue, batching |
| **Prāṇa** | Yes (state changes) | <1ms | Lock-free reads |
| **Vāyu** | No (load time) | 50-100ms | Lazy loading |
| **Śrī** | Yes (rendering) | <1ms | Style caching |
| **Dravya** | Yes (always) | <1ms | Minimal overhead |

---

## 5.10 Testing Strategies

### Māyā (Renderer)
- Snapshot testing of rendered output
- Diff algorithm correctness tests
- Layout constraint tests
- Performance benchmarks

### Vāk (Commands)
- Command parsing tests
- Completion engine tests
- Undo/redo correctness
- Integration tests with state

### Agni (Events)
- Event dispatch ordering tests
- Priority queue tests
- Timer accuracy tests
- Concurrency stress tests

### Prāṇa (State)
- Observer notification tests
- Derived state correctness
- Transaction isolation tests
- Memory leak tests

### Vāyu (Plugins)
- Plugin loading tests
- Dependency resolution tests
- Sandboxing tests
- Hot reload tests

### Śrī (Themes)
- Color conversion tests
- Animation interpolation tests
- Terminal capability tests

### Dravya (Runtime)
- Lifecycle tests
- Resource cleanup tests
- Error recovery tests
- Graceful shutdown tests

---

## Summary

DRAV's seven modules work in harmony:
- **Māyā** projects state to terminal
- **Vāk** processes commands
- **Agni** distributes events
- **Prāṇa** animates via reactive state
- **Vāyu** extends via plugins
- **Śrī** beautifies the interface
- **Dravya** provides the foundation

Each module is independently testable, loosely coupled, and highly cohesive.

**Next**: [API Specifications](brief-06-api-specifications.md) provides detailed public API reference.

---

[← Back to Index](brief-index.md) | [Previous: System Architecture](brief-04-system-architecture.md) | [Next: API Specifications →](brief-06-api-specifications.md)
