# Sri - Theme & Style

**Package:** `github.com/TIVerse/drav/pkg/sri`

## Overview

Sri (Sanskrit: श्री, "beauty, grace") is DRAV's theming and styling system. It provides consistent visual design, color palettes, animations, and accessibility features.

## Key Concepts

### Themes

Themes define the visual appearance:

```go
type Theme struct {
    Name        string
    Palette     Palette
    Styles      StyleSet
    Animations  AnimationSet
}
```

### Palettes

Color schemes for UI elements:

```go
type Palette struct {
    Primary     Color
    Secondary   Color
    Background  Color
    Foreground  Color
    Success     Color
    Warning     Color
    Error       Color
    Info        Color
}
```

### Styles

Visual properties for components:

```go
type Style struct {
    Foreground  Color
    Background  Color
    Border      BorderStyle
    Padding     Spacing
    Margin      Spacing
    FontWeight  FontWeight
    Underline   bool
    Bold        bool
    Italic      bool
}
```

## Built-in Themes

### Default Theme

```go
theme := sri.DefaultTheme()
```

### Dark Theme

```go
theme := sri.DarkTheme()
```

### Light Theme

```go
theme := sri.LightTheme()
```

### High Contrast

```go
theme := sri.HighContrastTheme()
```

## Creating Themes

### Custom Theme

```go
myTheme := sri.Theme{
    Name: "My Theme",
    Palette: sri.Palette{
        Primary:    sri.RGB(66, 135, 245),
        Secondary:  sri.RGB(156, 39, 176),
        Background: sri.RGB(18, 18, 18),
        Foreground: sri.RGB(255, 255, 255),
        Success:    sri.RGB(76, 175, 80),
        Warning:    sri.RGB(255, 152, 0),
        Error:      sri.RGB(244, 67, 54),
        Info:       sri.RGB(33, 150, 243),
    },
}
```

### Theme from File

```go
theme, err := sri.LoadTheme("./themes/my-theme.yaml")
```

YAML format:

```yaml
name: My Theme
palette:
  primary: "#4287f5"
  secondary: "#9c27b0"
  background: "#121212"
  foreground: "#ffffff"
  success: "#4caf50"
  warning: "#ff9800"
  error: "#f44336"
  info: "#2196f3"
```

## Applying Themes

### Set Active Theme

```go
themeMgr := sri.NewThemeManager()
themeMgr.SetActive(myTheme)
```

### Hot-Swapping

Switch themes at runtime:

```go
app.OnReady(func() {
    // Register keyboard shortcut
    eventHub := app.EventHub()
    eventHub.On(agni.EventTypeKey, func(ctx context.Context, event dravya.Event) error {
        keyEvent := event.(*agni.KeyEvent)
        
        // Ctrl+T to toggle theme
        if keyEvent.Rune == 't' && keyEvent.Modifiers&agni.ModCtrl != 0 {
            currentTheme := themeMgr.Current()
            if currentTheme.Name == "dark" {
                themeMgr.SetActive(sri.LightTheme())
            } else {
                themeMgr.SetActive(sri.DarkTheme())
            }
            app.RequestRender()
        }
        return nil
    })
})
```

## Styling Components

### Using Theme Colors

```go
func (c *MyComponent) Render(ctx maya.RenderContext) maya.View {
    theme := sri.CurrentTheme()
    
    return maya.Text("Hello",
        maya.WithForeground(theme.Palette.Primary),
        maya.WithBackground(theme.Palette.Background),
    )
}
```

### Component Styles

```go
buttonStyle := theme.Styles.Button
inputStyle := theme.Styles.Input
panelStyle := theme.Styles.Panel
```

### Dynamic Styles

```go
func (c *Button) style() sri.Style {
    theme := sri.CurrentTheme()
    base := theme.Styles.Button
    
    if c.focused {
        base.Border = sri.BorderDouble
        base.Foreground = theme.Palette.Primary
    }
    
    if c.disabled {
        base.Foreground = theme.Palette.Disabled
    }
    
    return base
}
```

## Color System

### RGB Colors

```go
color := sri.RGB(255, 100, 50)
```

### Hex Colors

```go
color := sri.Hex("#FF6432")
```

### Named Colors

```go
red := sri.Red
green := sri.Green
blue := sri.Blue
white := sri.White
black := sri.Black
```

### Color Manipulation

```go
// Lighten
lighter := sri.Lighten(color, 0.2)

// Darken
darker := sri.Darken(color, 0.2)

// Adjust alpha
transparent := sri.WithAlpha(color, 0.5)

// Mix colors
mixed := sri.Mix(color1, color2, 0.5)
```

## Animations

### Easing Functions

```go
type Easing int

const (
    EaseLinear Easing = iota
    EaseInQuad
    EaseOutQuad
    EaseInOutQuad
    EaseInCubic
    EaseOutCubic
    EaseInOutCubic
)
```

### Creating Animations

```go
anim := sri.Animation{
    Duration: time.Millisecond * 300,
    Easing:   sri.EaseInOutQuad,
    From:     0,
    To:       100,
}

// Get value at time t
value := anim.ValueAt(time.Millisecond * 150)
```

### Animated Component

```go
type AnimatedButton struct {
    animation *sri.Animation
    progress  float64
}

func (b *AnimatedButton) Render(ctx maya.RenderContext) maya.View {
    // Animate color
    theme := sri.CurrentTheme()
    color := sri.Interpolate(
        theme.Palette.Primary,
        theme.Palette.Secondary,
        b.progress,
    )
    
    return maya.Button(b.label, b.onClick,
        maya.WithForeground(color),
    )
}

func (b *AnimatedButton) update(delta time.Duration) {
    b.progress = b.animation.ValueAt(delta)
    if b.progress >= 1.0 {
        b.animation.Reset()
    }
}
```

## Accessibility

### High Contrast Mode

```go
if sri.IsHighContrastMode() {
    theme = sri.HighContrastTheme()
}
```

### Color Blind Safe

```go
palette := sri.ColorBlindSafePalette()
```

### WCAG Compliance

Check contrast ratios:

```go
ratio := sri.ContrastRatio(foreground, background)
if ratio < 4.5 {
    // WCAG AA fails for normal text
    log.Warning("Insufficient contrast")
}
```

### Screen Reader Support

```go
component := maya.Panel("Panel Title", content,
    maya.WithAriaLabel("Main content panel"),
    maya.WithAriaRole("region"),
)
```

## Style Sets

### Predefined Styles

```go
type StyleSet struct {
    Button  Style
    Input   Style
    Panel   Style
    Text    Style
    List    Style
    Table   Style
    Modal   Style
}
```

### Applying Style Sets

```go
button := maya.Button("Click", handler)
buttonWithStyle := sri.ApplyStyle(button, theme.Styles.Button)
```

## Best Practices

### 1. Use Theme Colors

Always use palette colors:

```go
// Good
color := theme.Palette.Primary

// Bad
color := sri.RGB(66, 135, 245)  // Hardcoded
```

### 2. Semantic Colors

Use semantic color names:

```go
// Good
errorText := maya.Text("Error!", 
    maya.WithForeground(theme.Palette.Error))

// Bad
errorText := maya.Text("Error!",
    maya.WithForeground(sri.Red))  // Less semantic
```

### 3. Respect User Preferences

Detect system theme:

```go
if sri.IsSystemDarkMode() {
    theme = sri.DarkTheme()
} else {
    theme = sri.LightTheme()
}
```

### 4. Smooth Transitions

Use animations for theme changes:

```go
themeMgr.SetActiveWithTransition(newTheme, time.Millisecond*300)
```

## Patterns

### Theme Switcher

```go
type ThemeSwitcher struct {
    manager  *sri.ThemeManager
    themes   []sri.Theme
    current  int
}

func (s *ThemeSwitcher) Next() {
    s.current = (s.current + 1) % len(s.themes)
    s.manager.SetActive(s.themes[s.current])
}

func (s *ThemeSwitcher) Render(ctx maya.RenderContext) maya.View {
    theme := s.manager.Current()
    return maya.Button(
        fmt.Sprintf("Theme: %s", theme.Name),
        s.Next,
    )
}
```

### Conditional Styling

```go
func (c *Component) style() sri.Style {
    base := sri.CurrentTheme().Styles.Base
    
    switch c.state {
    case StateNormal:
        return base
    case StateHover:
        return sri.Merge(base, sri.Style{
            Foreground: theme.Palette.Primary,
        })
    case StateActive:
        return sri.Merge(base, sri.Style{
            Background: theme.Palette.Primary,
            Foreground: theme.Palette.Background,
        })
    case StateDisabled:
        return sri.Merge(base, sri.Style{
            Foreground: theme.Palette.Disabled,
        })
    }
    
    return base
}
```

### Custom Palette

```go
func createBrandPalette(brandColor sri.Color) sri.Palette {
    return sri.Palette{
        Primary:    brandColor,
        Secondary:  sri.Lighten(brandColor, 0.2),
        Background: sri.Darken(brandColor, 0.9),
        Foreground: sri.White,
        Success:    sri.Green,
        Warning:    sri.Orange,
        Error:      sri.Red,
        Info:       sri.Blue,
    }
}
```

## Examples

### Custom Theme Example

```go
func main() {
    // Create custom theme
    myTheme := sri.Theme{
        Name: "Ocean",
        Palette: sri.Palette{
            Primary:    sri.Hex("#0077BE"),
            Secondary:  sri.Hex("#00BCD4"),
            Background: sri.Hex("#FFFFFF"),
            Foreground: sri.Hex("#212121"),
        },
    }
    
    // Apply theme
    themeMgr := sri.NewThemeManager()
    themeMgr.SetActive(myTheme)
    
    app := dravya.NewApp()
    app.SetRoot(&MyComponent{})
    app.Run(context.Background())
}
```

### Animated Component

```go
type PulsingButton struct {
    pulse *sri.Animation
}

func NewPulsingButton() *PulsingButton {
    return &PulsingButton{
        pulse: sri.NewAnimation(
            time.Second,
            sri.EaseInOutQuad,
            0, 1,
        ),
    }
}

func (b *PulsingButton) Render(ctx maya.RenderContext) maya.View {
    theme := sri.CurrentTheme()
    progress := b.pulse.Progress()
    
    color := sri.Interpolate(
        theme.Palette.Primary,
        theme.Palette.Secondary,
        progress,
    )
    
    return maya.Button("Click Me", b.onClick,
        maya.WithForeground(color),
    )
}
```

## Performance Considerations

### Color Caching

Cache computed colors:

```go
type ThemeCache struct {
    colors map[string]sri.Color
}

func (c *ThemeCache) Get(key string, compute func() sri.Color) sri.Color {
    if color, ok := c.colors[key]; ok {
        return color
    }
    color := compute()
    c.colors[key] = color
    return color
}
```

### Minimize Style Changes

Apply styles at component level:

```go
// Good
styledPanel := sri.ApplyStyle(panel, style)

// Bad - many individual properties
panel.WithForeground(color1)
     .WithBackground(color2)
     .WithBorder(border)
     // ...
```

## Troubleshooting

### Colors Not Displaying

Check terminal color support:

```bash
echo $TERM
# Should be xterm-256color or similar
```

### Theme Not Applying

Ensure theme is set before rendering:

```go
app.OnReady(func() {
    themeMgr.SetActive(myTheme)
})
```

### Animation Stuttering

Use fixed time step:

```go
anim.Update(time.Millisecond * 16)  // ~60 FPS
```

## Related Modules

- **[Maya](maya.md)**: Applies styles to components
- **[Dravya](dravya.md)**: Theme manager integration

## See Also

- [Theme Examples](../../examples/)
- [Color Theory](../concepts.md#colors)
- [Accessibility Guide](../accessibility.md)
