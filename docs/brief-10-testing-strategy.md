# Section 10: Testing Strategy

[← Back to Index](brief-index.md) | [Previous: Security Architecture](brief-09-security-architecture.md) | [Next: Developer Experience →](brief-11-developer-experience.md)

---

## 10.1 Testing Philosophy

### Testing Pyramid

```
                    ┌─────────┐
                    │   E2E   │  10%  (Slow, expensive)
                    └─────────┘
                 ┌───────────────┐
                 │  Integration  │  20%  (Medium speed)
                 └───────────────┘
              ┌─────────────────────┐
              │      Unit Tests      │  70%  (Fast, cheap)
              └─────────────────────┘
```

**Rationale**: Most tests should be fast unit tests. Integration and E2E tests catch integration issues but should be minority.

### Coverage Targets

| Component | Unit | Integration | E2E | Total |
|-----------|------|-------------|-----|-------|
| **Core Modules** | 95% | 80% | 50% | 90% |
| **Renderer** | 95% | 90% | 60% | 92% |
| **Commands** | 90% | 85% | 70% | 88% |
| **Events** | 95% | 85% | 60% | 90% |
| **State** | 98% | 90% | 50% | 95% |
| **Plugins** | 85% | 80% | 70% | 85% |
| **Overall** | 95% | 85% | 60% | 90% |

---

## 10.2 Unit Testing

### Module-Level Unit Tests

#### Testing Observables
```go
func TestObservable_SetNotifiesObservers(t *testing.T) {
    obs := NewObservable(0)
    
    var notified bool
    obs.Watch(func(v int) {
        notified = true
    })
    
    obs.Set(1)
    
    assert.True(t, notified, "observer should be notified")
}

func TestObservable_MultipleObservers(t *testing.T) {
    obs := NewObservable(0)
    
    count := 0
    obs.Watch(func(v int) { count++ })
    obs.Watch(func(v int) { count++ })
    
    obs.Set(1)
    
    assert.Equal(t, 2, count, "both observers should be notified")
}
```

#### Testing Command Parser
```go
func TestParser_ParseSimpleCommand(t *testing.T) {
    parser := NewParser()
    
    cmd, err := parser.Parse("save document.txt")
    
    assert.NoError(t, err)
    assert.Equal(t, "save", cmd.Name)
    assert.Equal(t, []string{"document.txt"}, cmd.Args)
}

func TestParser_ParseCommandWithFlags(t *testing.T) {
    parser := NewParser()
    
    cmd, err := parser.Parse("search pattern --case-sensitive --context=3")
    
    assert.NoError(t, err)
    assert.Equal(t, "search", cmd.Name)
    assert.Equal(t, []string{"pattern"}, cmd.Args)
    assert.Equal(t, true, cmd.Flags["case-sensitive"])
    assert.Equal(t, "3", cmd.Flags["context"])
}

func TestParser_QuotedArguments(t *testing.T) {
    parser := NewParser()
    
    cmd, err := parser.Parse(`save "my document.txt"`)
    
    assert.NoError(t, err)
    assert.Equal(t, []string{"my document.txt"}, cmd.Args)
}
```

#### Testing Diff Algorithm
```go
func TestDiff_NoChanges(t *testing.T) {
    oldBuf := createBuffer("Hello")
    newBuf := createBuffer("Hello")
    
    changes := Diff(oldBuf, newBuf)
    
    assert.Empty(t, changes, "no changes expected")
}

func TestDiff_SingleCellChange(t *testing.T) {
    oldBuf := createBuffer("Hello")
    newBuf := createBuffer("Hallo")
    
    changes := Diff(oldBuf, newBuf)
    
    assert.Len(t, changes, 1)
    assert.Equal(t, 1, changes[0].X)
    assert.Equal(t, 'a', changes[0].NewCell.Rune)
}

func BenchmarkDiff_LargeBuffer(b *testing.B) {
    oldBuf := createLargeBuffer(200, 50)
    newBuf := createLargeBuffer(200, 50)
    // Change 10% of cells
    modifyRandomCells(newBuf, 0.1)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        Diff(oldBuf, newBuf)
    }
}
```

### Test Helpers

```go
// Test environment setup
type TestEnv struct {
    app      *App
    renderer *MockRenderer
    events   *MockEventHub
}

func NewTestEnv() *TestEnv {
    return &TestEnv{
        app:      NewApp(),
        renderer: NewMockRenderer(),
        events:   NewMockEventHub(),
    }
}

func (te *TestEnv) Render(component Component) View {
    return component.Render()
}

func (te *TestEnv) SendKey(key Key) {
    te.events.Dispatch(KeyEvent{Key: key})
}
```

---

## 10.3 Integration Testing

### Module Integration Tests

#### Testing State → UI Integration
```go
func TestStateUIIntegration(t *testing.T) {
    env := NewTestEnv()
    
    // Create observable state
    counter := NewObservable(0)
    
    // Create component that observes state
    component := &Counter{count: counter}
    env.app.SetRoot(component)
    
    // Initial render
    view := env.Render(component)
    assert.Contains(t, viewToString(view), "Count: 0")
    
    // Change state
    counter.Set(5)
    
    // Wait for re-render
    time.Sleep(50 * time.Millisecond)
    
    // Assert UI updated
    view = env.Render(component)
    assert.Contains(t, viewToString(view), "Count: 5")
}
```

#### Testing Command → State Integration
```go
func TestCommandStateIntegration(t *testing.T) {
    env := NewTestEnv()
    
    counter := NewObservable(0)
    
    // Register command that modifies state
    env.app.RegisterCommand("increment", Command{
        Handler: func(ctx Context, args []string) error {
            counter.Set(counter.Get() + 1)
            return nil
        },
    })
    
    // Execute command
    err := env.app.Execute("increment")
    assert.NoError(t, err)
    
    // Assert state changed
    assert.Equal(t, 1, counter.Get())
}
```

#### Testing Event → Command Integration
```go
func TestEventCommandIntegration(t *testing.T) {
    env := NewTestEnv()
    
    executed := false
    
    // Bind key to command
    env.app.OnKey(KeyEnter, func(e KeyEvent) {
        env.app.Execute("mycommand")
    })
    
    env.app.RegisterCommand("mycommand", Command{
        Handler: func(ctx Context, args []string) error {
            executed = true
            return nil
        },
    })
    
    // Send key event
    env.SendKey(KeyEnter)
    
    // Wait for async execution
    time.Sleep(50 * time.Millisecond)
    
    assert.True(t, executed)
}
```

---

## 10.4 End-to-End Testing

### Full Application Tests

```go
func TestE2E_CounterApp(t *testing.T) {
    // Start application in test mode
    app := NewApp(WithTestMode())
    
    // Set up counter app
    counter := NewObservable(0)
    app.SetRoot(&CounterComponent{count: counter})
    
    // Start app in goroutine
    go app.Run()
    defer app.Shutdown()
    
    // Wait for startup
    time.Sleep(100 * time.Millisecond)
    
    // Simulate user pressing increment button
    app.SendKey(KeyEnter)
    time.Sleep(50 * time.Millisecond)
    
    // Take screenshot
    screenshot := app.Screenshot()
    
    // Assert counter incremented in UI
    assert.Contains(t, screenshot, "Count: 1")
    
    // Press again
    app.SendKey(KeyEnter)
    time.Sleep(50 * time.Millisecond)
    
    screenshot = app.Screenshot()
    assert.Contains(t, screenshot, "Count: 2")
}
```

### Visual Regression Testing

```go
func TestVisualRegression_Dashboard(t *testing.T) {
    app := setupDashboardApp()
    
    // Render to virtual terminal
    output := app.RenderToString()
    
    // Compare with golden snapshot
    goldenPath := "testdata/dashboard.golden"
    
    if os.Getenv("UPDATE_GOLDEN") == "1" {
        // Update golden file
        os.WriteFile(goldenPath, []byte(output), 0644)
    } else {
        // Compare with golden
        golden, _ := os.ReadFile(goldenPath)
        assert.Equal(t, string(golden), output)
    }
}
```

---

## 10.5 Property-Based Testing

### Using gopter

```go
import "github.com/leanovate/gopter"

func TestObservable_Properties(t *testing.T) {
    properties := gopter.NewProperties(nil)
    
    // Property: Get always returns last Set value
    properties.Property("Get returns last Set value", 
        prop.ForAll(
            func(value int) bool {
                obs := NewObservable(0)
                obs.Set(value)
                return obs.Get() == value
            },
            gen.Int(),
        ),
    )
    
    // Property: Multiple Sets, last one wins
    properties.Property("Multiple Sets, last wins",
        prop.ForAll(
            func(values []int) bool {
                if len(values) == 0 {
                    return true
                }
                obs := NewObservable(0)
                for _, v := range values {
                    obs.Set(v)
                }
                return obs.Get() == values[len(values)-1]
            },
            gen.SliceOf(gen.Int()),
        ),
    )
    
    properties.TestingRun(t)
}
```

---

## 10.6 Performance Testing

### Benchmark Suite

```go
func BenchmarkRenderFullScreen(b *testing.B) {
    app := setupComplexApp()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        app.Render()
    }
}

func BenchmarkObservableUpdate(b *testing.B) {
    obs := NewObservable(0)
    obs.Watch(func(v int) {})
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        obs.Set(i)
    }
}

func BenchmarkCommandParsing(b *testing.B) {
    parser := NewParser()
    input := "command arg1 arg2 --flag1 --flag2=value"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        parser.Parse(input)
    }
}
```

### Load Testing

```go
func TestLoad_1000Components(t *testing.T) {
    app := NewApp()
    
    // Create 1000 components
    components := make([]Component, 1000)
    for i := 0; i < 1000; i++ {
        components[i] = NewTextComponent(fmt.Sprintf("Component %d", i))
    }
    
    app.SetRoot(Column(components...))
    
    // Measure render time
    start := time.Now()
    app.Render()
    duration := time.Since(start)
    
    // Should render in < 100ms
    assert.Less(t, duration, 100*time.Millisecond)
}

func TestLoad_10000Observables(t *testing.T) {
    // Create 10,000 observables with observers
    observables := make([]*Observable[int], 10000)
    for i := 0; i < 10000; i++ {
        obs := NewObservable(0)
        obs.Watch(func(v int) {})
        observables[i] = obs
    }
    
    // Update all
    start := time.Now()
    for _, obs := range observables {
        obs.Set(1)
    }
    duration := time.Since(start)
    
    // Should complete in < 1s
    assert.Less(t, duration, 1*time.Second)
}
```

---

## 10.7 Fuzzing

### Fuzz Testing Input Parser

```go
func FuzzCommandParser(f *testing.F) {
    // Seed corpus
    f.Add("save file.txt")
    f.Add("search pattern --flag")
    f.Add(`command "quoted arg"`)
    
    parser := NewParser()
    
    f.Fuzz(func(t *testing.T, input string) {
        // Should not panic
        _, err := parser.Parse(input)
        
        // If no error, result should be valid
        if err == nil {
            // Additional validation
        }
    })
}
```

---

## 10.8 Test Organization

### Directory Structure

```
test/
├── unit/
│   ├── maya_test.go
│   ├── vak_test.go
│   ├── agni_test.go
│   ├── prana_test.go
│   └── ...
├── integration/
│   ├── state_ui_test.go
│   ├── command_state_test.go
│   └── event_command_test.go
├── e2e/
│   ├── counter_app_test.go
│   ├── dashboard_app_test.go
│   └── editor_app_test.go
├── benchmark/
│   ├── render_bench_test.go
│   ├── state_bench_test.go
│   └── parser_bench_test.go
├── testdata/
│   ├── golden/
│   │   ├── dashboard.golden
│   │   └── editor.golden
│   └── fixtures/
└── helpers/
    ├── test_env.go
    ├── mocks.go
    └── assertions.go
```

### Test Tags

```go
// +build unit

package test

// Unit tests only
```

```go
// +build integration

package test

// Integration tests
```

Run specific tests:
```bash
go test -tags=unit ./...
go test -tags=integration ./...
go test -tags=e2e ./...
```

---

## 10.9 Mocking

### Mock Interfaces

```go
type MockRenderer struct {
    mock.Mock
}

func (m *MockRenderer) Render(view View) error {
    args := m.Called(view)
    return args.Error(0)
}

func (m *MockRenderer) Size() (int, int) {
    args := m.Called()
    return args.Int(0), args.Int(1)
}

// Usage
func TestComponentRender(t *testing.T) {
    renderer := new(MockRenderer)
    renderer.On("Render", mock.Anything).Return(nil)
    
    component := NewComponent()
    component.Render()
    
    renderer.AssertExpectations(t)
}
```

---

## 10.10 Continuous Testing

### GitHub Actions Workflow

```yaml
name: Tests

on: [push, pull_request]

jobs:
  unit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: go test -tags=unit -race -coverprofile=coverage.txt ./...
      - uses: codecov/codecov-action@v3

  integration:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - run: go test -tags=integration ./...

  e2e:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - run: go test -tags=e2e ./...
```

---

## 10.11 Test Quality Metrics

### Measuring Test Effectiveness

```go
// Mutation testing
// Change implementation, tests should fail
func TestMutationCoverage(t *testing.T) {
    // Original: counter++
    // Mutation: counter--
    // Tests should catch this
}
```

### Test Metrics Dashboard

- **Code Coverage**: 95%+ target
- **Test Execution Time**: < 5 minutes
- **Flaky Test Rate**: < 1%
- **Test-to-Code Ratio**: ~1:1 lines

---

## Summary

**Testing Strategy**:
- 70% unit, 20% integration, 10% E2E
- 95% overall coverage target
- Continuous testing in CI
- Property-based and fuzz testing
- Performance benchmarks
- Visual regression testing

**Key Practices**:
- Test-driven development
- Mock external dependencies
- Fast feedback loops
- Comprehensive test suite
- Automated testing pipeline

**Tools**:
- Go testing framework
- testify for assertions
- gopter for property testing
- pprof for benchmarking
- Golden files for snapshots

---

[← Back to Index](brief-index.md) | [Previous: Security Architecture](brief-09-security-architecture.md) | [Next: Developer Experience →](brief-11-developer-experience.md)
