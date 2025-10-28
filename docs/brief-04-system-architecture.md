# Section 4: System Architecture

[← Back to Index](brief-index.md) | [Previous: Market Analysis](brief-03-market-analysis.md) | [Next: Core Modules →](brief-05-core-modules.md)

---

## 4.1 Architectural Overview

### High-Level Architecture

DRAV follows a **modular, event-driven architecture** with seven core modules organized in three layers:

```
┌─────────────────────────────────────────────────────────────────┐
│                     APPLICATION LAYER                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │ User Apps   │  │  Widgets    │  │  Plugins    │             │
│  │ & Commands  │  │ & Themes    │  │ & Extensions│             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
└─────────────────────────┬───────────────────────────────────────┘
                          │
┌─────────────────────────▼───────────────────────────────────────┐
│                     FRAMEWORK LAYER                              │
│                                                                   │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │              UI & PRESENTATION                            │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │  │
│  │  │  Renderer    │  │    Theme     │  │   Command    │   │  │
│  │  │   (Māyā)     │  │   Engine     │  │   Engine     │   │  │
│  │  │              │  │    (Śrī)     │  │    (Vāk)     │   │  │
│  │  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘   │  │
│  └─────────┼──────────────────┼──────────────────┼───────────┘  │
│            │                  │                  │               │
│  ┌─────────▼──────────────────▼──────────────────▼───────────┐  │
│  │             STATE & CONTROL                               │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │  │
│  │  │  Reactive    │  │  Event Hub   │  │   Plugin     │   │  │
│  │  │    State     │  │   (Agni)     │  │   System     │   │  │
│  │  │  (Prāṇa)     │  │              │  │   (Vāyu)     │   │  │
│  │  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘   │  │
│  └─────────┼──────────────────┼──────────────────┼───────────┘  │
│            │                  │                  │               │
│            └──────────────────┴──────────────────┘               │
│                               │                                  │
│  ┌────────────────────────────▼────────────────────────────┐    │
│  │                    RUNTIME CORE                          │    │
│  │                     (Dravya)                             │    │
│  │  • Lifecycle Management  • Resource Allocation           │    │
│  │  • Error Recovery       • Concurrency Coordination       │    │
│  └──────────────────────────┬────────────────────────────┘    │
└─────────────────────────────┼───────────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────────┐
│                     PLATFORM LAYER                               │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐        │
│  │  tcell   │  │ terminfo │  │   pty    │  │  termios │        │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘        │
│  ┌──────────────────────────────────────────────────────┐       │
│  │         OS (Windows / macOS / Linux)                  │       │
│  └──────────────────────────────────────────────────────┘       │
└──────────────────────────────────────────────────────────────────┘
```

### Architectural Principles

#### 1. Separation of Concerns

Each module has a **single, well-defined responsibility**:

- **Māyā** (Renderer): Translates virtual UI to terminal output
- **Vāk** (Commands): Processes user commands
- **Agni** (Events): Dispatches events to handlers
- **Prāṇa** (State): Manages observable state
- **Vāyu** (Plugins): Loads and manages extensions
- **Śrī** (Themes): Applies visual styling
- **Dravya** (Runtime): Coordinates everything

#### 2. Loose Coupling via Interfaces

Modules communicate through **well-defined interfaces**, not concrete types:

```go
// Modules depend on interfaces, not implementations
type Renderer interface {
    Render(view View) error
    Size() (width, height int)
}

type EventHub interface {
    Subscribe(eventType EventType, handler Handler)
    Publish(event Event)
}

// Easy to mock, test, and replace implementations
```

#### 3. Inversion of Control

The **framework controls the flow**, not user code:

```go
// Framework calls user code
app := drav.NewApp()
app.OnInit(func(ctx Context) {
    // Framework calls this when ready
})
app.OnEvent(KeyPress, func(e KeyEvent) {
    // Framework calls this on key press
})
app.Run()  // Framework takes control
```

#### 4. Dependency Injection

Components receive dependencies via **constructors**, enabling testing:

```go
type Renderer struct {
    state  StateManager
    events EventHub
    theme  ThemeEngine
}

func NewRenderer(state StateManager, events EventHub, theme ThemeEngine) *Renderer {
    // Dependencies injected, easy to mock for testing
    return &Renderer{state, events, theme}
}
```

---

## 4.2 Data Flow Architecture

### Event-Driven Flow

DRAV uses **unidirectional data flow** for predictability:

```
┌─────────────────────────────────────────────────────────┐
│                    DATA FLOW CYCLE                       │
└─────────────────────────────────────────────────────────┘

   User Input (Keyboard/Mouse)
          │
          ▼
   ┌─────────────┐
   │  Event Hub  │  (Agni captures and dispatches)
   │   (Agni)    │
   └──────┬──────┘
          │
          ▼
   ┌─────────────┐
   │  Handlers   │  (Commands or event handlers)
   │   (Vāk)     │
   └──────┬──────┘
          │
          ▼
   ┌─────────────┐
   │    State    │  (Observable state updated)
   │  (Prāṇa)    │
   └──────┬──────┘
          │ (Automatic notification)
          ▼
   ┌─────────────┐
   │ Components  │  (Re-render triggered)
   └──────┬──────┘
          │
          ▼
   ┌─────────────┐
   │   Renderer  │  (Diff and draw)
   │   (Māyā)    │
   └──────┬──────┘
          │
          ▼
   Terminal Output
          │
          └──────────── Loop ─────────┐
                                       │
          User sees update, provides input
```

### State Flow Patterns

#### Pattern 1: Direct State Update

```go
// Simple: Component directly updates state
button := drav.Button("Click me", func() {
    counter.Set(counter.Get() + 1)  // Direct update
})
// Observer pattern automatically triggers re-render
```

#### Pattern 2: Command-Driven Update

```go
// User types :increment
// → Vāk parses command
// → Handler executes
// → State updated
// → UI automatically refreshes
drav.Command("increment", func(ctx Context, args []string) {
    amount := parseInt(args[0])
    counter.Set(counter.Get() + amount)
})
```

#### Pattern 3: Event-Driven Update

```go
// System event triggers update
events.On(ResizeEvent, func(e ResizeEvent) {
    windowSize.Set(Size{e.Width, e.Height})
    // All components observing windowSize re-render
})
```

---

## 4.3 Concurrency Model

### Goroutine Architecture

DRAV uses **structured concurrency** with clear goroutine ownership:

```
Main Goroutine (Application)
│
├─ Event Loop Goroutine
│  ├─ Input Reader
│  └─ Event Dispatcher
│
├─ Renderer Goroutine
│  ├─ Diff Calculator
│  └─ Terminal Writer
│
├─ Command Executor Pool
│  ├─ Worker 1
│  ├─ Worker 2
│  ├─ ...
│  └─ Worker N
│
├─ State Observer Pool
│  ├─ Observer 1
│  ├─ Observer 2
│  └─ ...
│
└─ Plugin Manager
   ├─ Plugin 1 Goroutines
   └─ Plugin 2 Goroutines
```

### Synchronization Primitives

#### 1. Channels for Communication

```go
// Event dispatch via channels
type EventHub struct {
    events chan Event
    stop   chan struct{}
}

func (h *EventHub) Run() {
    for {
        select {
        case event := <-h.events:
            h.dispatch(event)
        case <-h.stop:
            return
        }
    }
}
```

#### 2. Mutexes for Shared State

```go
// State protected by mutex
type Observable[T any] struct {
    value     T
    observers []Observer[T]
    mu        sync.RWMutex
}

func (o *Observable[T]) Get() T {
    o.mu.RLock()
    defer o.mu.RUnlock()
    return o.value
}

func (o *Observable[T]) Set(value T) {
    o.mu.Lock()
    o.value = value
    observers := o.observers
    o.mu.Unlock()
    
    // Notify outside lock to prevent deadlocks
    for _, obs := range observers {
        obs.Notify(value)
    }
}
```

#### 3. Context for Cancellation

```go
// Use context for cancellation propagation
func (app *App) Run(ctx context.Context) error {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()
    
    // Pass context to all components
    go app.events.Run(ctx)
    go app.renderer.Run(ctx)
    go app.commands.Run(ctx)
    
    // Wait for ctx.Done()
    <-ctx.Done()
    return ctx.Err()
}
```

### Deadlock Prevention

**Strategy 1: Lock Ordering**
```go
// Always acquire locks in consistent order
// Rule: State lock before UI lock
state.mu.Lock()
ui.mu.Lock()
// Never reverse!
ui.mu.Unlock()
state.mu.Unlock()
```

**Strategy 2: Release Before Notify**
```go
// Don't call observers while holding lock
o.mu.Lock()
observers := append([]Observer[T]{}, o.observers...)  // Copy
o.mu.Unlock()

// Notify without lock
for _, obs := range observers {
    obs.Notify(value)
}
```

**Strategy 3: Timeout Guards**
```go
// Use select with timeout to detect deadlocks
select {
case result := <-ch:
    return result
case <-time.After(5 * time.Second):
    return fmt.Errorf("operation timed out - possible deadlock")
}
```

---

## 4.4 Module Interaction Patterns

### Pattern 1: Publish-Subscribe (Event Hub)

```go
// Agni (Event Hub) implements pub-sub
type EventHub struct {
    subscribers map[EventType][]chan Event
}

// Subscribe to event type
func (h *EventHub) Subscribe(eventType EventType) <-chan Event {
    ch := make(chan Event, 10)
    h.subscribers[eventType] = append(h.subscribers[eventType], ch)
    return ch
}

// Publish to all subscribers
func (h *EventHub) Publish(event Event) {
    for _, ch := range h.subscribers[event.Type()] {
        select {
        case ch <- event:
        default:
            // Subscriber slow, drop event
        }
    }
}
```

### Pattern 2: Observer (Reactive State)

```go
// Prāṇa (State) implements observer pattern
type Observable[T any] struct {
    value     T
    observers []func(T)
}

func (o *Observable[T]) Watch(fn func(T)) Unsubscribe {
    o.observers = append(o.observers, fn)
    return func() { /* remove fn */ }
}

func (o *Observable[T]) Set(value T) {
    o.value = value
    for _, obs := range o.observers {
        obs(value)  // Notify all observers
    }
}
```

### Pattern 3: Command Pattern (Command Engine)

```go
// Vāk (Commands) implements command pattern
type Command struct {
    execute func(Context, []string) error
    undo    func(Context) error
}

type CommandHistory struct {
    executed []Command
}

func (h *CommandHistory) Execute(cmd Command, ctx Context, args []string) error {
    err := cmd.execute(ctx, args)
    if err == nil {
        h.executed = append(h.executed, cmd)
    }
    return err
}

func (h *CommandHistory) Undo() error {
    if len(h.executed) == 0 {
        return errors.New("nothing to undo")
    }
    cmd := h.executed[len(h.executed)-1]
    h.executed = h.executed[:len(h.executed)-1]
    return cmd.undo(ctx)
}
```

### Pattern 4: Plugin (Extension Points)

```go
// Vāyu (Plugins) provides extension points via hooks
type HookRegistry struct {
    beforeRender []func(View) View
    afterRender  []func()
    onCommand    []func(string) bool
}

func (r *HookRegistry) RegisterBeforeRender(fn func(View) View) {
    r.beforeRender = append(r.beforeRender, fn)
}

func (r *HookRegistry) TriggerBeforeRender(view View) View {
    for _, hook := range r.beforeRender {
        view = hook(view)  // Chain transformations
    }
    return view
}
```

---

## 4.5 Rendering Architecture

### Virtual UI Tree

DRAV maintains a **virtual representation** of the UI:

```go
type View interface {
    Type() string
    Children() []View
    Props() map[string]interface{}
    
    // Layout computes positions
    Layout(constraints Constraints) Dimensions
    
    // Draw renders to buffer
    Draw(buf *Buffer, region Region)
}

// Example view tree
Root
├─ Row
│  ├─ Panel ("CPU")
│  │  └─ Text ("85%")
│  └─ Panel ("Memory")
│     └─ Text ("2.4GB")
└─ Footer
   └─ Text ("Press ? for help")
```

### Diff Algorithm

**Two-Pass Differential Rendering**:

#### Pass 1: Line-Level Hashing

```go
func hashLine(cells []Cell) uint64 {
    h := fnv.New64a()
    for _, cell := range cells {
        binary.Write(h, binary.LittleEndian, cell.Rune)
        binary.Write(h, binary.LittleEndian, cell.Style)
    }
    return h.Sum64()
}

// Compare line hashes
oldHashes := computeLineHashes(oldBuffer)
newHashes := computeLineHashes(newBuffer)

dirtyLines := []int{}
for i := 0; i < height; i++ {
    if oldHashes[i] != newHashes[i] {
        dirtyLines = append(dirtyLines, i)
    }
}
```

**Complexity**: O(height) — fast for typical terminal sizes (50-100 lines)

#### Pass 2: Cell-Level Diff (Only Dirty Lines)

```go
for _, lineNum := range dirtyLines {
    oldLine := oldBuffer.Line(lineNum)
    newLine := newBuffer.Line(lineNum)
    
    // Find contiguous changed regions
    regions := findChangedRegions(oldLine, newLine)
    
    // Update only changed regions
    for _, region := range regions {
        terminal.MoveCursor(region.X, lineNum)
        terminal.Write(newLine[region.Start:region.End])
    }
}
```

**Complexity**: O(width × dirty_lines)

### Layout System

DRAV uses a **flexbox-inspired** layout algorithm:

```go
type LayoutNode struct {
    Type       NodeType  // Row, Column, Stack, Grid
    Children   []*LayoutNode
    Flex       FlexProps
    Style      Style
    Computed   *ComputedLayout
}

type FlexProps struct {
    Grow       float64  // Flex grow factor
    Shrink     float64  // Flex shrink factor
    Basis      Size     // Initial size
    Align      Alignment
    Justify    Justification
}
```

**Layout Algorithm**:

```
function Layout(node, constraints):
    if node is leaf:
        return node.measure(constraints)
    
    if node.type == Row:
        return layoutRow(node, constraints)
    
    if node.type == Column:
        return layoutColumn(node, constraints)
    
    // ... other types

function layoutRow(node, constraints):
    // 1. Measure intrinsic sizes
    intrinsicSizes := []
    totalFlex := 0
    fixedWidth := 0
    
    for each child:
        size := child.measure(unbounded)
        intrinsicSizes.append(size)
        if child.flex.grow > 0:
            totalFlex += child.flex.grow
        else:
            fixedWidth += size.width
    
    // 2. Distribute remaining space
    remainingWidth := constraints.width - fixedWidth
    
    for each child:
        if child.flex.grow > 0:
            flexWidth := remainingWidth * (child.flex.grow / totalFlex)
            child.layout(Constraints{width: flexWidth, ...})
        else:
            child.layout(Constraints{width: intrinsicSizes[i].width, ...})
    
    // 3. Position children
    x := 0
    for each child:
        child.computed.x = x
        child.computed.y = alignChild(child, node.flex.align)
        x += child.computed.width
```

---

## 4.6 Event System Architecture

### Event Types Hierarchy

```
Event (interface)
├─ InputEvent
│  ├─ KeyEvent
│  ├─ MouseEvent
│  └─ PasteEvent
├─ SystemEvent
│  ├─ ResizeEvent
│  ├─ FocusEvent
│  └─ SignalEvent
├─ TimerEvent
│  ├─ TickEvent
│  └─ TimeoutEvent
└─ CustomEvent
   └─ UserDefinedEvent
```

### Event Priority Queue

Events are processed by **priority**:

```go
type PriorityQueue struct {
    critical chan Event  // 0-10ms (Signals, Resize)
    high     chan Event  // 10-50ms (Keyboard, Mouse)
    normal   chan Event  // 50-100ms (Ticks, Custom)
    low      chan Event  // 100ms+ (Background)
}

func (pq *PriorityQueue) Dispatch() {
    for {
        select {
        case e := <-pq.critical:
            handle(e)
        case e := <-pq.high:
            if len(pq.critical) == 0 {  // Only if no critical
                handle(e)
            }
        case e := <-pq.normal:
            if len(pq.critical) == 0 && len(pq.high) == 0 {
                handle(e)
            }
        case e := <-pq.low:
            if len(pq.critical) == 0 && len(pq.high) == 0 && len(pq.normal) == 0 {
                handle(e)
            }
        }
    }
}
```

### Event Handler Registration

```go
// Fluent API for event handlers
component.
    OnKey(KeyEnter, handleEnter).
    OnKey(KeyEsc, handleEscape).
    OnMouse(MouseLeft, handleClick).
    OnResize(handleResize).
    OnTick(time.Second, handleTick)
```

---

## 4.7 Plugin Architecture

### Plugin Interface

```go
type Plugin interface {
    // Metadata
    Name() string
    Version() string
    Author() string
    Dependencies() []Dependency
    
    // Lifecycle hooks
    Init(runtime Runtime) error
    Start(ctx Context) error
    Stop() error
    
    // Registration hooks
    RegisterCommands(registry CommandRegistry)
    RegisterWidgets(registry WidgetRegistry)
    RegisterThemes(registry ThemeRegistry)
    RegisterEvents(registry EventRegistry)
    
    // Optional hooks
    OnBeforeRender(view View) View
    OnAfterRender()
    OnCommand(cmd string) bool
}
```

### Plugin Loading Mechanism

```
┌────────────────────────────────────────┐
│         Plugin Loading Process          │
└────────────────────────────────────────┘

1. Discovery
   ├─ Scan ~/.drav/plugins/
   ├─ Check manifest.json
   └─ Validate signature

2. Dependency Resolution
   ├─ Build dependency graph
   ├─ Check conflicts
   └─ Determine load order

3. Loading
   ├─ Load .so file (Go plugin)
   ├─ Or load .wasm (WebAssembly)
   └─ Invoke Init()

4. Registration
   ├─ Call RegisterCommands()
   ├─ Call RegisterWidgets()
   └─ Call RegisterThemes()

5. Activation
   ├─ Call Start()
   └─ Plugin now active

6. Hot Reload (optional)
   ├─ Watch plugin file
   ├─ On change: Stop() → Unload → Load → Start()
   └─ State migration
```

### Plugin Security

**Sandboxing Strategies**:

1. **Resource Limits**:
```go
type PluginContext struct {
    maxMemory  int64
    maxCPU     time.Duration
    maxFiles   int
    maxThreads int
}
```

2. **Capability-Based Access**:
```go
// Plugins only get capabilities they request
type Capabilities struct {
    FileSystem   bool
    Network      bool
    StateAccess  bool
    UIAccess     bool
}
```

3. **API Isolation**:
```go
// Plugins access framework via restricted interface
type PluginAPI interface {
    Log(level, message string)
    GetState(key string) (interface{}, error)
    SetState(key string, value interface{}) error
    // No direct access to internals
}
```

---

## Summary

DRAV's architecture prioritizes:

1. **Modularity**: Seven independent, cohesive modules
2. **Reactivity**: Automatic state-to-UI synchronization
3. **Performance**: Diff-based rendering, efficient layout
4. **Concurrency**: Structured goroutines, deadlock prevention
5. **Extensibility**: Plugin system with security

**Key Innovations**:
- Terminal-specific diff algorithm
- Observable state with automatic UI updates
- Integrated command engine
- Hot-reloadable plugins

**Next**: [Core Modules](brief-05-core-modules.md) provides detailed implementation specs for each module.

---

[← Back to Index](brief-index.md) | [Previous: Market Analysis](brief-03-market-analysis.md) | [Next: Core Modules →](brief-05-core-modules.md)
