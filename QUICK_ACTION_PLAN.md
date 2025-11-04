# DRAV Quick Action Plan
**Priority: Critical Missing Features for v0.1**

## 🚨 Top 5 Blockers (Do First)

### 1. Event System Integration ⏱️ 3-4 days
**Why Critical:** Examples don't respond to keyboard input  
**Current State:** Event system exists but disconnected from app

**Tasks:**
```go
// pkg/dravya/event_loop.go - Create this file
- [ ] Poll tcell events in App.Run() loop
- [ ] Route key events to focused component
- [ ] Handle resize events → renderer
- [ ] Handle Ctrl+C gracefully

// Update examples/02-counter/main.go
- [ ] Add event handler for +/- keys
- [ ] Connect to counter increment/decrement
- [ ] Test interactive behavior
```

**Files to create/modify:**
- `pkg/dravya/event_loop.go` (new)
- `pkg/maya/event_handler.go` (new)
- `examples/02-counter/main.go` (modify)

---

### 2. Input Widget ⏱️ 2-3 days
**Why Critical:** Can't build any real apps without text input  
**Current State:** Missing completely

**Tasks:**
```go
// pkg/maya/widgets/input.go - Create this
type Input struct {
    value     *prana.Observable[string]
    cursor    int
    maxLength int
    placeholder string
}

- [ ] Text rendering with cursor
- [ ] Character insertion/deletion
- [ ] Cursor movement (left/right/home/end)
- [ ] Copy/paste support
- [ ] Input validation
```

**Example usage:**
```go
input := maya.Input(
    maya.WithPlaceholder("Enter name..."),
    maya.WithMaxLength(50),
)
```

---

### 3. Focus Management ⏱️ 2 days
**Why Critical:** Can't navigate between interactive elements  
**Current State:** Missing completely

**Tasks:**
```go
// pkg/maya/focus.go - Create this
type FocusManager struct {
    focused    string
    focusables []Focusable
    order      []string
}

- [ ] Track focusable components
- [ ] Tab/Shift+Tab navigation
- [ ] Focus enter/exit events
- [ ] Visual focus indicator
- [ ] Programmatic focus API
```

---

### 4. List Widget ⏱️ 2-3 days
**Why Critical:** Very common UI pattern  
**Current State:** Missing completely

**Tasks:**
```go
// pkg/maya/widgets/list.go
type List struct {
    items    *prana.Observable[[]string]
    selected int
    offset   int // for scrolling
}

- [ ] Render items with selection highlight
- [ ] Up/down arrow navigation
- [ ] Scrolling for overflow
- [ ] Multi-select support (optional)
- [ ] Custom item renderer
```

---

### 5. Panel Widget ⏱️ 1-2 days
**Why Critical:** Basic visual organization  
**Current State:** Missing completely

**Tasks:**
```go
// pkg/maya/widgets/panel.go
type Panel struct {
    title    string
    children []View
    border   BorderStyle
}

- [ ] Border rendering (single, double, rounded)
- [ ] Title in border
- [ ] Padding support
- [ ] Content clipping
```

---

## 📋 Week 1 Sprint Plan

### Day 1-2: Event Integration
- Set up event polling from tcell
- Route events to app components
- Make counter example interactive
- Test on Linux/Mac/Windows

### Day 3-4: Input + Focus
- Implement Input widget
- Implement FocusManager
- Create example: form with 3 inputs + Tab navigation

### Day 5: List + Panel
- Implement List widget
- Implement Panel widget
- Create example: file browser mockup

### Weekend: Testing
- Write unit tests for new widgets
- Integration test for event flow
- Update documentation

---

## 🎯 Quick Wins (Low Effort, High Impact)

### 1. Fix Command Palette Display (2 hours)
```go
// pkg/vak/palette_widget.go
- Add basic overlay component
- Show registered commands
- Handle Enter to execute
```

### 2. Add More View Options (1 hour)
```go
// pkg/maya/view.go - Add these options
- WithBorder(BorderStyle)
- WithTitle(string)
- WithMaxWidth(int)
- WithMaxHeight(int)
```

### 3. Connect Observable to Re-render (3 hours)
```go
// pkg/prana/observable.go
- Add App reference to trigger re-render
- Call app.RequestRender() on Set()
- Batch updates within same frame
```

### 4. Add Table Widget (4-6 hours)
```go
// pkg/maya/widgets/table.go
- Fixed column widths
- Header row
- Scrolling
- Cell styling
```

---

## 🐛 Known Issues to Fix

1. **Renderer doesn't handle resize** - Buffer size fixed at init
   - Add resize event handler
   - Reallocate buffers on resize

2. **No error boundaries** - Panics crash entire app
   - Add recover() in component Render()
   - Show error in UI instead of crash

3. **Plugin loader is all stubs** - WASM functions are TODOs
   - Either implement or remove from v0.1
   - Document as "coming in v0.2"

4. **Examples don't match README** - Code in README different from examples
   - Update README examples to match actual API
   - Or update examples to match README

---

## 📦 Recommended Dependencies to Add

```go
// For WASM plugins (Phase 2)
github.com/tetratelabs/wazero

// For fuzzy search
github.com/sahilm/fuzzy

// For better testing
github.com/stretchr/testify/assert
github.com/stretchr/testify/mock

// For terminal capabilities
github.com/muesli/termenv (optional, for better color detection)
```

---

## 🚀 After Week 1 - Next Priorities

### Week 2: Advanced Widgets
- Table widget
- Tabs widget
- Modal/Dialog widget
- Button widget

### Week 3: Command Palette UI
- Fuzzy search implementation
- Visual command palette
- Keyboard shortcuts system

### Week 4: Testing & Examples
- Achieve 70% test coverage
- Complete all 7 example directories
- Create 2-3 real-world examples

---

## ✅ Definition of Done for v0.1-alpha

- [ ] All 5 top blockers implemented and tested
- [ ] Counter example is fully interactive
- [ ] New example: Interactive form with validation
- [ ] New example: File browser with list navigation
- [ ] Test coverage > 60%
- [ ] Works on Linux, macOS, Windows
- [ ] Basic module documentation written
- [ ] CHANGELOG.md updated
- [ ] Tagged as v0.1.0-alpha

**Estimated Time:** 2-3 weeks full-time

---

## 🛠️ Development Commands

```bash
# Run specific example
go run ./examples/02-counter

# Run all tests
make test

# Run with race detector
go test -race ./...

# Build all examples
make examples

# Run linter
make lint

# Generate coverage report
make cover
```

---

## 📞 Getting Help

If stuck on:
- **Event system:** Check `pkg/agni/dispatcher.go` implementation
- **Rendering:** Look at `pkg/maya/renderer.go` renderView logic
- **State:** Reference `pkg/prana/observable_test.go` for patterns
- **tcell:** See `pkg/maya/driver_tcell.go` for terminal integration

---

**Remember:** Focus on making it work first, then make it right, then make it fast. Ship alpha early, iterate based on feedback.
