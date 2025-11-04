# Maya - Renderer & UI

**Package:** `github.com/TIVerse/drav/pkg/maya`

## Overview

Maya (Sanskrit: माया, "illusion, appearance") is DRAV's rendering engine. It manages the virtual UI tree, performs efficient diffing, and renders to the terminal using tcell.

## Key Concepts

### Components

Components implement the `Component` interface:

```go
type Component interface {
    Render(ctx RenderContext) View
}
```

### Views

Views are immutable virtual UI nodes returned from `Render()`:

```go
func (c *MyComponent) Render(ctx maya.RenderContext) maya.View {
    return maya.Column(
        maya.Text("Hello"),
        maya.Text("World"),
    )
}
```

### Rendering Pipeline

1. Component renders to View tree
2. Differ compares with previous tree
3. Only changed cells are updated
4. Double-buffering prevents flicker

## Built-in Widgets

### Text

Display static or dynamic text:

```go
maya.Text("Hello, World!")
maya.Text(fmt.Sprintf("Count: %d", count))
```

### Panel

Container with border and optional title:

```go
maya.Panel("My Panel", []maya.View{
    maya.Text("Content"),
}, maya.WithBorder(maya.BorderSingle))
```

### Button

Clickable button with callback:

```go
maya.Button("Click Me", func() {
    // Handle click
}, maya.WithFocus(true))
```

### Input

Text input field:

```go
maya.Input(value,
    maya.WithPlaceholder("Enter text..."),
    maya.WithOnChange(handleChange),
)
```

### List

Scrollable list of items:

```go
maya.List([]maya.ListItem{
    {Label: "Item 1", Value: "val1"},
    {Label: "Item 2", Value: "val2"},
})
```

### Table

Data table with columns and rows:

```go
maya.Table(
    []maya.TableColumn{
        {Header: "Name", Width: 20},
        {Header: "Age", Width: 10},
    },
    []maya.TableRow{
        {Cells: []string{"Alice", "30"}},
        {Cells: []string{"Bob", "25"}},
    },
)
```

### Tabs

Tabbed interface:

```go
maya.Tabs([]maya.TabItem{
    {Label: "Tab 1", Content: view1},
    {Label: "Tab 2", Content: view2},
})
```

### Modal

Dialog overlay:

```go
maya.Modal("Confirm", contentView,
    maya.WithCloseable(true),
)
```

## Layout Containers

### Row

Horizontal layout:

```go
maya.Row(
    maya.Text("Left"),
    maya.Text("Right"),
)
```

### Column

Vertical layout:

```go
maya.Column(
    maya.Text("Top"),
    maya.Text("Bottom"),
)
```

## View Options

### Sizing

```go
maya.WithSize(width, height)
maya.WithGrow(1)  // Flex grow factor
```

### Spacing

```go
maya.WithPadding(maya.AllSides(2))
maya.WithMargin(maya.Spacing{Top: 1, Bottom: 1})
```

### Styling

```go
maya.WithForeground(maya.RGB(255, 255, 255))
maya.WithBackground(maya.RGB(0, 0, 255))
maya.WithBorder(maya.BorderDouble)
```

### Alignment

```go
maya.WithAlign(maya.AlignCenter)
```

### Focus

```go
maya.WithFocus(true)
```

## Focus Management

Register focusable components:

```go
focusMgr := app.FocusManager()
focusMgr.Register("my-button")

// Check if focused
if focusMgr.IsFocused("my-button") {
    // Render with focus styling
}
```

Navigation:
- **Tab**: Focus next component
- **Shift+Tab**: Focus previous component

## Rendering Context

The `RenderContext` provides terminal information:

```go
func (c *MyComponent) Render(ctx maya.RenderContext) maya.View {
    width := ctx.Width()
    height := ctx.Height()
    focused := ctx.Focused()
    
    return maya.Text(fmt.Sprintf("Size: %dx%d", width, height))
}
```

## Color System

### RGB Colors

```go
red := maya.RGB(255, 0, 0)
```

### Default Color

```go
defaultColor := maya.DefaultColor()
```

## Border Styles

- `BorderNone`: No border
- `BorderSingle`: Single-line box (─│┌┐└┘)
- `BorderDouble`: Double-line box (═║╔╗╚╝)
- `BorderRounded`: Rounded corners (─│╭╮╰╯)
- `BorderThick`: Thick lines

## Best Practices

### 1. Immutable Views

Never mutate views after creation. Always create new ones:

```go
// Good
return maya.Text(newText)

// Bad - don't try to mutate
view.Content = newText
```

### 2. Efficient Rendering

Cache expensive view construction:

```go
type MyComponent struct {
    cachedView maya.View
    lastData   Data
}

func (c *MyComponent) Render(ctx maya.RenderContext) maya.View {
    if c.data != c.lastData {
        c.cachedView = c.buildView()
        c.lastData = c.data
    }
    return c.cachedView
}
```

### 3. Responsive Layouts

Use context dimensions:

```go
func (c *MyComponent) Render(ctx maya.RenderContext) maya.View {
    if ctx.Width() < 80 {
        return c.renderCompact()
    }
    return c.renderFull()
}
```

## Examples

### Simple Component

```go
type Greeting struct {
    name string
}

func (g *Greeting) Render(ctx maya.RenderContext) maya.View {
    return maya.Panel("Greeting", []maya.View{
        maya.Text(fmt.Sprintf("Hello, %s!", g.name)),
    })
}
```

### Interactive Form

```go
type LoginForm struct {
    username *prana.Observable[string]
    password *prana.Observable[string]
}

func (f *LoginForm) Render(ctx maya.RenderContext) maya.View {
    return maya.Column(
        maya.Text("Login"),
        maya.Input(f.username.Get(),
            maya.WithPlaceholder("Username"),
        ),
        maya.Input(f.password.Get(),
            maya.WithPlaceholder("Password"),
        ),
        maya.Button("Submit", f.handleLogin),
    )
}
```

## Performance Tips

- **Minimize View Depth**: Flatten nested structures where possible
- **Avoid Large Tables**: Use virtual scrolling for >1000 rows
- **Optimize Strings**: Pre-format expensive strings outside `Render()`
- **Use Keys**: Provide unique IDs for list items to aid diffing

## Troubleshooting

### Rendering Artifacts

Ensure terminal supports required features:

```bash
echo $TERM  # Should be xterm-256color or similar
```

### Layout Issues

Check size constraints:

```go
maya.Text("Long text...").WithSize(20, 1)  // May truncate
```

### Focus Not Working

Register components with focus manager:

```go
app.FocusManager().Register(componentID)
```

## Related Modules

- **[Dravya](dravya.md)**: Application runtime
- **[Prana](prana.md)**: State management for reactive UIs
- **[Sri](sri.md)**: Theming and styling

## See Also

- [Layout Examples](../../examples/03-dashboard/)
- [Widget Gallery](../../examples/)
