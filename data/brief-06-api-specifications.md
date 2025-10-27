# Section 6: API Specifications

[← Back to Index](brief-index.md) | [Previous: Core Modules](brief-05-core-modules.md) | [Next: Implementation Strategy →](brief-07-implementation-strategy.md)

---

## 6.1 Core Application API

### Application Initialization

```go
// Create new application
app := drav.NewApp(options ...Option)

// Configuration options
type Option func(*App)

func WithTheme(theme Theme) Option
func WithLogLevel(level LogLevel) Option
func WithPluginDir(dir string) Option
func WithConfig(config Config) Option

// Example
app := drav.NewApp(
    drav.WithTheme(drav.DarkTheme),
    drav.WithLogLevel(drav.Debug),
)
```

### Application Lifecycle

```go
// Initialize application
func (app *App) Init() error

// Set root component
func (app *App) SetRoot(component Component) error

// Run application (blocking)
func (app *App) Run() error

// Run with context
func (app *App) RunContext(ctx context.Context) error

// Shutdown gracefully
func (app *App) Shutdown() error

// Example
app.SetRoot(myRootComponent)
if err := app.Run(); err != nil {
    log.Fatal(err)
}
```

---

## 6.2 Component API

### Component Interface

```go
type Component interface {
    // Render returns view representation
    Render() View
    
    // Update handles messages (optional)
    Update(msg Message) Command
    
    // Init runs on mount (optional)
    Init() Command
    
    // Lifecycle hooks
    OnMount()
    OnUnmount()
}

// Functional component
func MyComponent(props Props) Component {
    return drav.Func(func() View {
        return drav.Text("Hello, DRAV!")
    })
}
```

### Stateful Components

```go
type Counter struct {
    count drav.Observable[int]
}

func NewCounter() *Counter {
    return &Counter{
        count: drav.NewObservable(0),
    }
}

func (c *Counter) Render() View {
    return drav.Column(
        drav.Text(fmt.Sprintf("Count: %d", c.count.Get())),
        drav.Button("Increment", func() {
            c.count.Set(c.count.Get() + 1)
        }),
        drav.Button("Decrement", func() {
            c.count.Set(c.count.Get() - 1)
        }),
    )
}
```

---

## 6.3 View API

### Layout Primitives

```go
// Row - horizontal layout
func Row(children ...View) View

// Column - vertical layout
func Column(children ...View) View

// Stack - layered layout
func Stack(children ...View) View

// Flex - flexible box layout
func Flex(children ...FlexChild) View

// Grid - grid layout
func Grid(rows, cols int, children ...View) View

// Example
drav.Row(
    drav.Panel("Left", leftContent),
    drav.Column(
        drav.Panel("Top Right", topContent),
        drav.Panel("Bottom Right", bottomContent),
    ),
)
```

### Basic Widgets

```go
// Text display
func Text(content string, style ...Style) View
func Textf(format string, args ...interface{}) View

// Panel with border and title
func Panel(title string, content View, style ...Style) View

// Button
func Button(label string, onClick func(), style ...Style) View

// Input field
func Input(value Observable[string], placeholder string) View

// Checkbox
func Checkbox(label string, checked Observable[bool]) View

// Select dropdown
func Select(options []string, selected Observable[int]) View

// List
func List(items []string, selected Observable[int]) View

// Table
func Table(headers []string, rows [][]string) View

// Progress bar
func ProgressBar(value Observable[float64]) View

// Spinner
func Spinner(message string) View
```

### Advanced Widgets

```go
// Chart widgets
func LineChart(data Observable[[]Point]) View
func BarChart(data Observable[[]Bar]) View
func PieChart(data Observable[[]Slice]) View

// Tree view
func Tree(root TreeNode, expanded Observable[map[string]bool]) View

// Tabs
func Tabs(tabs []Tab, active Observable[int]) View

// Modal
func Modal(title string, content View, onClose func()) View

// Notification
func Notification(message string, level Level) View
```

---

## 6.4 State Management API

### Observable

```go
// Create observable
func NewObservable[T any](initial T) *Observable[T]

// Get current value
func (o *Observable[T]) Get() T

// Set new value (triggers observers)
func (o *Observable[T]) Set(value T)

// Watch for changes
func (o *Observable[T]) Watch(fn func(T)) Unsubscribe

// Map to derived observable
func (o *Observable[T]) Map(fn func(T) T) *Observable[T]

// Example
counter := drav.NewObservable(0)
doubled := counter.Map(func(v int) int { return v * 2 })

counter.Watch(func(v int) {
    fmt.Printf("Counter changed to: %d\n", v)
})

counter.Set(5)  // Prints: "Counter changed to: 5"
// doubled is now 10
```

### Computed State

```go
// Compute from multiple observables
func Computed[T any](fn func() T, deps ...Observable[any]) *Observable[T]

// Example
firstName := drav.NewObservable("John")
lastName := drav.NewObservable("Doe")

fullName := drav.Computed(func() string {
    return firstName.Get() + " " + lastName.Get()
}, firstName, lastName)

firstName.Set("Jane")
// fullName automatically becomes "Jane Doe"
```

### Store Pattern

```go
type Store[T any] struct {
    state   *Observable[T]
    actions map[string]Action[T]
}

type Action[T any] func(*T, ...interface{})

func NewStore[T any](initial T) *Store[T]

func (s *Store[T]) Dispatch(action string, args ...interface{})

// Example
type AppState struct {
    Counter int
    User    *User
}

store := drav.NewStore(AppState{Counter: 0})
store.RegisterAction("increment", func(state *AppState) {
    state.Counter++
})
store.RegisterAction("setUser", func(state *AppState, user *User) {
    state.User = user
})

store.Dispatch("increment")
store.Dispatch("setUser", currentUser)
```

---

## 6.5 Command API

### Command Registration

```go
// Register command
func RegisterCommand(name string, cmd Command) error

// Command structure
type Command struct {
    Name        string
    Description string
    Usage       string
    Examples    []string
    Handler     func(ctx Context, args []string) error
    Complete    Completer
    Flags       []Flag
}

// Example
drav.RegisterCommand("save", drav.Command{
    Description: "Save current document",
    Usage:       "save [filename]",
    Examples: []string{
        "save",
        "save document.txt",
        "save --force",
    },
    Flags: []drav.Flag{
        {Name: "force", Short: "f", Type: drav.Bool},
    },
    Handler: func(ctx drav.Context, args []string) error {
        filename := "untitled.txt"
        if len(args) > 0 {
            filename = args[0]
        }
        force := ctx.Flag("force").Bool()
        return saveDocument(filename, force)
    },
})
```

### Command Execution

```go
// Execute command programmatically
func Execute(command string) error

// Execute with context
func ExecuteContext(ctx context.Context, command string) error

// Example
drav.Execute("save document.txt --force")
```

### Autocompletion

```go
// Built-in completers
func FileCompleter() Completer
func DirectoryCompleter() Completer
func CommandCompleter() Completer

// Custom completer
type Completer func(partial string, pos int) []Completion

type Completion struct {
    Text        string
    Display     string
    Description string
}

// Example
func envCompleter(partial string, pos int) []Completion {
    completions := []Completion{}
    for _, env := range os.Environ() {
        parts := strings.SplitN(env, "=", 2)
        if strings.HasPrefix(parts[0], partial) {
            completions = append(completions, Completion{
                Text:        parts[0],
                Display:     parts[0],
                Description: parts[1],
            })
        }
    }
    return completions
}
```

---

## 6.6 Event API

### Event Handlers

```go
// Register global event handler
func OnKey(key Key, handler func(KeyEvent)) Unsubscribe
func OnMouse(button MouseButton, handler func(MouseEvent)) Unsubscribe
func OnResize(handler func(ResizeEvent)) Unsubscribe
func OnTick(interval time.Duration, handler func(TickEvent)) Unsubscribe

// Component-level events
component.OnKey(drav.KeyEnter, func(e drav.KeyEvent) {
    // Handle enter key
})

component.OnMouse(drav.MouseLeft, func(e drav.MouseEvent) {
    // Handle left click at e.X, e.Y
})
```

### Custom Events

```go
// Define custom event
type DataLoadedEvent struct {
    Data interface{}
}

func (e DataLoadedEvent) Type() drav.EventType {
    return "data_loaded"
}

func (e DataLoadedEvent) Timestamp() time.Time {
    return time.Now()
}

// Publish custom event
drav.Publish(DataLoadedEvent{Data: myData})

// Subscribe to custom event
drav.On("data_loaded", func(event drav.Event) {
    e := event.(DataLoadedEvent)
    fmt.Printf("Data loaded: %v\n", e.Data)
})
```

---

## 6.7 Plugin API

### Plugin Development

```go
// Plugin interface
type Plugin interface {
    Name() string
    Version() string
    Init(runtime Runtime) error
    RegisterCommands(registry CommandRegistry)
    RegisterWidgets(registry WidgetRegistry)
}

// Example plugin
type GitPlugin struct{}

func (p *GitPlugin) Name() string {
    return "git-integration"
}

func (p *GitPlugin) Version() string {
    return "1.0.0"
}

func (p *GitPlugin) Init(runtime drav.Runtime) error {
    // Initialize plugin
    return nil
}

func (p *GitPlugin) RegisterCommands(registry drav.CommandRegistry) {
    registry.Register("git", drav.Command{
        Description: "Git operations",
        Handler:     handleGitCommand,
    })
}

func (p *GitPlugin) RegisterWidgets(registry drav.WidgetRegistry) {
    registry.Register("GitStatus", NewGitStatusWidget)
}

// Load plugin
drav.LoadPlugin(&GitPlugin{})
```

### Plugin Loading

```go
// Load plugin from file
func LoadPluginFile(path string) error

// Load plugin from bytes
func LoadPluginBytes(data []byte) error

// Unload plugin
func UnloadPlugin(name string) error

// List loaded plugins
func ListPlugins() []PluginInfo
```

---

## 6.8 Theme API

### Theme Definition

```go
type Theme struct {
    Name   string
    Colors ColorPalette
    Styles StyleMap
}

type ColorPalette struct {
    Background Color
    Foreground Color
    Primary    Color
    Secondary  Color
    Accent     Color
    Success    Color
    Warning    Color
    Error      Color
    Info       Color
}

// Built-in themes
var (
    DarkTheme  = Theme{ /* ... */ }
    LightTheme = Theme{ /* ... */ }
    SolarizedDark = Theme{ /* ... */ }
    Dracula = Theme{ /* ... */ }
)

// Apply theme
drav.SetTheme(drav.DarkTheme)
```

### Custom Themes

```go
// Create custom theme
myTheme := drav.Theme{
    Name: "Custom Theme",
    Colors: drav.ColorPalette{
        Background: drav.RGB(30, 30, 30),
        Foreground: drav.RGB(200, 200, 200),
        Primary:    drav.RGB(100, 150, 255),
        Accent:     drav.RGB(255, 100, 150),
    },
}

drav.RegisterTheme(myTheme)
drav.SetTheme(myTheme)
```

### Style API

```go
// Style structure
type Style struct {
    Fg     Color
    Bg     Color
    Bold   bool
    Italic bool
    Underline bool
}

// Style builder
style := drav.NewStyle().
    Foreground(drav.Blue).
    Background(drav.Black).
    Bold(true).
    Underline(true)

// Apply style to view
drav.Text("Hello", style)
```

---

## 6.9 Animation API

### Basic Animations

```go
// Fade in/out
view.FadeIn(duration time.Duration)
view.FadeOut(duration time.Duration)

// Slide
view.SlideIn(direction Direction, duration time.Duration)
view.SlideOut(direction Direction, duration time.Duration)

// Scale
view.ScaleTo(scale float64, duration time.Duration)

// Pulse
view.Pulse(duration time.Duration)

// Example
panel := drav.Panel("Welcome", content)
panel.FadeIn(300 * time.Millisecond)
```

### Easing Functions

```go
// Built-in easing
type Easing int

const (
    Linear Easing = iota
    EaseIn
    EaseOut
    EaseInOut
    EaseInQuad
    EaseOutQuad
    EaseInCubic
    EaseOutCubic
)

// Custom animation
view.Animate(drav.Animation{
    Property: "opacity",
    From:     0.0,
    To:       1.0,
    Duration: 500 * time.Millisecond,
    Easing:   drav.EaseInOut,
})
```

### Transition Groups

```go
// Animate multiple properties
view.Transition(drav.Transition{
    Properties: []drav.AnimatedProperty{
        {Property: "x", From: 0, To: 100},
        {Property: "y", From: 0, To: 50},
        {Property: "opacity", From: 0, To: 1},
    },
    Duration: 300 * time.Millisecond,
    Easing:   drav.EaseOut,
})
```

---

## 6.10 Utility API

### Context

```go
type Context interface {
    // Application state
    State() StateManager
    
    // Command execution
    Execute(command string) error
    
    // Event publishing
    Publish(event Event)
    
    // Logging
    Log(level LogLevel, message string)
    Logf(level LogLevel, format string, args ...interface{})
    
    // Flag access (in command handlers)
    Flag(name string) FlagValue
}
```

### Logging

```go
// Log levels
type LogLevel int

const (
    Trace LogLevel = iota
    Debug
    Info
    Warn
    Error
    Fatal
)

// Logging functions
drav.Log(drav.Info, "Application started")
drav.Logf(drav.Debug, "User ID: %d", userID)
drav.Error(err)
```

### Configuration

```go
// Load configuration
config, err := drav.LoadConfig("config.yaml")

// Configuration structure
type Config struct {
    Theme      string
    LogLevel   string
    PluginDir  string
    KeyBindings map[string]string
    Custom     map[string]interface{}
}

// Access config
value := config.Get("custom.key")
config.Set("custom.key", "value")
config.Save()
```

---

## 6.11 Testing API

### Component Testing

```go
// Create test environment
env := drav.NewTestEnv()

// Render component
component := NewCounter()
view := env.Render(component)

// Assertions
env.AssertText(view, "Count: 0")

// Simulate events
env.SendKey(drav.KeyEnter)
env.SendMouse(drav.MouseLeft, 10, 5)

// Assert state changes
env.AssertText(view, "Count: 1")
```

### Snapshot Testing

```go
// Create snapshot
snapshot := env.Snapshot(component)

// Compare with saved snapshot
env.AssertMatchesSnapshot(snapshot, "counter-initial")

// Update snapshot
env.UpdateSnapshot("counter-initial", snapshot)
```

---

## 6.12 Error Handling

### Error Types

```go
type DravError struct {
    Code    ErrorCode
    Message string
    Cause   error
}

type ErrorCode int

const (
    ErrInvalidComponent ErrorCode = iota
    ErrCommandNotFound
    ErrPluginLoadFailed
    ErrStateCorrupted
    ErrRenderFailed
)
```

### Error Recovery

```go
// Error boundary
drav.ErrorBoundary(component, func(err error) View {
    return drav.Text(fmt.Sprintf("Error: %v", err))
})

// Global error handler
drav.OnError(func(err error) {
    drav.Logf(drav.Error, "Unhandled error: %v", err)
})
```

---

## Summary

DRAV provides comprehensive APIs for:
- **Application**: Lifecycle management
- **Components**: Functional and stateful
- **Views**: Layout and widgets
- **State**: Observables and stores
- **Commands**: Registration and execution
- **Events**: Handlers and custom events
- **Plugins**: Extension and loading
- **Themes**: Styling and animation
- **Testing**: Component and snapshot testing

All APIs follow Go idioms and prioritize simplicity and type safety.

---

[← Back to Index](brief-index.md) | [Previous: Core Modules](brief-05-core-modules.md) | [Next: Implementation Strategy →](brief-07-implementation-strategy.md)
