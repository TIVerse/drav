package maya

// Component is a renderable UI component.
type Component interface {
	Render(ctx RenderContext) View
}

// RenderContext provides context for rendering.
type RenderContext interface {
	Width() int
	Height() int
	Focused() bool
}

// BaseRenderContext provides a basic render context implementation.
type BaseRenderContext struct {
	WidthVal   int
	HeightVal  int
	FocusedVal bool
}

// Width returns the available width.
func (c *BaseRenderContext) Width() int {
	return c.WidthVal
}

// Height returns the available height.
func (c *BaseRenderContext) Height() int {
	return c.HeightVal
}

// Focused returns whether the component is focused.
func (c *BaseRenderContext) Focused() bool {
	return c.FocusedVal
}

// FuncComponent wraps a function as a component.
type FuncComponent struct {
	RenderFunc func(ctx RenderContext) View
}

// Render renders the function component.
func (fc *FuncComponent) Render(ctx RenderContext) View {
	return fc.RenderFunc(ctx)
}

// Func creates a functional component.
func Func(render func(ctx RenderContext) View) Component {
	return &FuncComponent{
		RenderFunc: render,
	}
}

// StatelessComponent is a component without state.
type StatelessComponent struct {
	view View
}

// Render returns the static view.
func (sc *StatelessComponent) Render(ctx RenderContext) View {
	return sc.view
}

// Stateless creates a stateless component.
func Stateless(view View) Component {
	return &StatelessComponent{
		view: view,
	}
}
