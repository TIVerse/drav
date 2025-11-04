# Prana - Reactive State

**Package:** `github.com/TIVerse/drav/pkg/prana`

## Overview

Prana (Sanskrit: प्राण, "life force, vital energy") provides reactive state management for DRAV applications. Observable values automatically trigger UI re-renders when changed.

## Key Concepts

### Observables

Observables are generic containers that notify watchers on change:

```go
counter := prana.NewObservable(0)
counter.Set(42)  // UI re-renders automatically
value := counter.Get()  // Returns 42
```

### Automatic Re-rendering

When an observable changes, DRAV automatically triggers a re-render:

```go
type Counter struct {
    count *prana.Observable[int]
}

func (c *Counter) increment() {
    c.count.Update(func(n int) int { return n + 1 })
    // UI updates automatically - no manual refresh needed!
}
```

### Type Safety

Observables are fully generic and type-safe:

```go
stringObs := prana.NewObservable[string]("hello")
intObs := prana.NewObservable[int](42)
structObs := prana.NewObservable[User](User{Name: "Alice"})
```

## Core API

### Creating Observables

```go
obs := prana.NewObservable(initialValue)
```

### Getting Values

```go
value := obs.Get()  // Thread-safe read
```

### Setting Values

```go
obs.Set(newValue)  // Triggers watchers and re-render
```

### Updating Values

Apply a function to transform the value:

```go
obs.Update(func(current int) int {
    return current + 1
})
```

### Watching for Changes

Register callbacks for change notifications:

```go
unwatch := obs.Watch(func(oldVal, newVal int) {
    fmt.Printf("Changed from %d to %d\n", oldVal, newVal)
})

// Clean up when done
defer unwatch()
```

## Advanced Features

### Computed Values

Derive observables from others:

```go
celsius := prana.NewObservable(0.0)
fahrenheit := celsius.Derive(func(c float64) float64 {
    return c*9/5 + 32
})

celsius.Set(100.0)
fmt.Println(fahrenheit.Get())  // 212.0
```

### Batch Updates

Update multiple observables without triggering multiple renders:

```go
batch := prana.NewBatch()
batch.Add(func() {
    name.Set("Alice")
    age.Set(30)
    city.Set("NYC")
})
batch.Execute()  // Single re-render for all changes
```

### Effects

Execute side effects when observables change:

```go
effect := prana.NewEffect(func() {
    // This runs when dependencies change
    fmt.Printf("Count is now: %d\n", counter.Get())
}, counter)

defer effect.Dispose()
```

### Stores

Organize related state with actions and reducers:

```go
type AppState struct {
    Count    int
    Username string
}

store := prana.NewStore(AppState{Count: 0})

// Register reducer
store.AddReducer("increment", func(state AppState, payload any) AppState {
    state.Count++
    return state
})

// Dispatch action
store.Dispatch("increment", nil)
```

## Patterns

### Component State

```go
type TodoList struct {
    items  *prana.Observable[[]Todo]
    filter *prana.Observable[FilterType]
}

func NewTodoList() *TodoList {
    return &TodoList{
        items:  prana.NewObservable([]Todo{}),
        filter: prana.NewObservable(FilterAll),
    }
}

func (t *TodoList) addTodo(text string) {
    t.items.Update(func(items []Todo) []Todo {
        return append(items, Todo{Text: text, Done: false})
    })
}
```

### Computed Properties

```go
type Dashboard struct {
    data      *prana.Observable[[]DataPoint]
    filtered  *prana.Observable[[]DataPoint]
    sortOrder *prana.Observable[SortOrder]
}

func NewDashboard() *Dashboard {
    d := &Dashboard{
        data:      prana.NewObservable([]DataPoint{}),
        sortOrder: prana.NewObservable(SortAsc),
    }
    
    // filtered updates when data or sortOrder changes
    d.filtered = prana.NewComputed(func() []DataPoint {
        data := d.data.Get()
        order := d.sortOrder.Get()
        return sortData(data, order)
    }, d.data, d.sortOrder)
    
    return d
}
```

### Form State

```go
type LoginForm struct {
    username *prana.Observable[string]
    password *prana.Observable[string]
    valid    *prana.Observable[bool]
}

func NewLoginForm() *LoginForm {
    f := &LoginForm{
        username: prana.NewObservable(""),
        password: prana.NewObservable(""),
    }
    
    // Auto-validate
    f.valid = prana.NewComputed(func() bool {
        return len(f.username.Get()) > 0 && len(f.password.Get()) >= 8
    }, f.username, f.password)
    
    return f
}
```

## Best Practices

### 1. Initialize in Constructor

Create observables during component construction:

```go
func NewCounter() *Counter {
    return &Counter{
        count: prana.NewObservable(0),  // Good
    }
}
```

### 2. Unsubscribe from Watchers

Always clean up watchers to prevent memory leaks:

```go
func (c *Component) init() {
    unwatch := observable.Watch(handler)
    c.cleanup = append(c.cleanup, unwatch)
}

func (c *Component) dispose() {
    for _, fn := range c.cleanup {
        fn()
    }
}
```

### 3. Use Computed for Derived State

Don't duplicate state - compute it:

```go
// Bad
type BadComponent struct {
    celsius    *prana.Observable[float64]
    fahrenheit *prana.Observable[float64]  // Duplicate state!
}

// Good
type GoodComponent struct {
    celsius    *prana.Observable[float64]
    fahrenheit *prana.Observable[float64]  // Computed from celsius
}
```

### 4. Batch Related Updates

```go
// Bad - multiple re-renders
user.name.Set("Alice")
user.age.Set(30)
user.city.Set("NYC")

// Good - single re-render
batch := prana.NewBatch()
batch.Add(func() {
    user.name.Set("Alice")
    user.age.Set(30)
    user.city.Set("NYC")
})
batch.Execute()
```

### 5. Keep Observables Private

Expose methods instead of raw observables:

```go
type Counter struct {
    count *prana.Observable[int]  // Private
}

func (c *Counter) Get() int {
    return c.count.Get()
}

func (c *Counter) Increment() {
    c.count.Update(func(n int) int { return n + 1 })
}
```

## Examples

### Counter Example

```go
type Counter struct {
    count *prana.Observable[int]
}

func NewCounter() *Counter {
    return &Counter{
        count: prana.NewObservable(0),
    }
}

func (c *Counter) Increment() {
    c.count.Update(func(n int) int { return n + 1 })
}

func (c *Counter) Render(ctx maya.RenderContext) maya.View {
    return maya.Column(
        maya.Text(fmt.Sprintf("Count: %d", c.count.Get())),
        maya.Button("+", c.Increment),
    )
}
```

### Todo List

```go
type Todo struct {
    Text string
    Done bool
}

type TodoList struct {
    items *prana.Observable[[]Todo]
}

func (t *TodoList) Add(text string) {
    t.items.Update(func(items []Todo) []Todo {
        return append(items, Todo{Text: text, Done: false})
    })
}

func (t *TodoList) Toggle(index int) {
    t.items.Update(func(items []Todo) []Todo {
        items[index].Done = !items[index].Done
        return items
    })
}
```

## Performance Considerations

### Thread Safety

All operations are thread-safe but come with a cost:

```go
// Fine for UI updates (infrequent)
obs.Set(value)

// For high-frequency updates, consider batching
batch.Add(func() {
    obs1.Set(v1)
    obs2.Set(v2)
})
```

### Watcher Cost

Each watcher adds overhead. Limit watchers:

```go
// Bad - too many watchers
for i := 0; i < 1000; i++ {
    obs.Watch(handlers[i])
}

// Good - single watcher with dispatch
obs.Watch(func(old, new any) {
    for _, handler := range handlers {
        handler(old, new)
    }
})
```

### Memory Leaks

Always unsubscribe:

```go
unwatch := obs.Watch(handler)
defer unwatch()  // Don't forget!
```

## Troubleshooting

### UI Not Updating

Ensure the observable is being modified, not replaced:

```go
// Bad
myComponent.items = newItems  // Direct assignment - no update!

// Good
myComponent.items.Set(newItems)  // Triggers update
```

### Infinite Loops

Avoid watching an observable and modifying it in the watcher:

```go
// Bad - infinite loop!
obs.Watch(func(old, new int) {
    obs.Set(new + 1)  // Triggers another change!
})

// Use computed instead
derived := obs.Derive(func(n int) int { return n + 1 })
```

## Related Modules

- **[Dravya](dravya.md)**: Connects observables to render system
- **[Maya](maya.md)**: Uses reactive state in components

## See Also

- [Counter Example](../../examples/02-counter/)
- [Reactive Patterns](../concepts.md#reactivity)
