package maya

// View represents a virtual UI tree node.
type View interface {
	Type() string
	Children() []View
	Attrs() Attributes
}

// Attributes holds view attributes.
type Attributes struct {
	ID          string
	Width       int
	Height      int
	MinWidth    int
	MinHeight   int
	MaxWidth    int
	MaxHeight   int
	Padding     Spacing
	Margin      Spacing
	Border      BorderStyle
	Background  Color
	Foreground  Color
	Align       Alignment
	Justify     Justification
	Grow        int
	Shrink      int
	Basis       int
	ZIndex      int
	Focused     bool
	Disabled    bool
	Hidden      bool
	Wrap        bool
}

// Spacing defines padding or margin.
type Spacing struct {
	Top    int
	Right  int
	Bottom int
	Left   int
}

// AllSides creates uniform spacing.
func AllSides(n int) Spacing {
	return Spacing{Top: n, Right: n, Bottom: n, Left: n}
}

// BorderStyle defines border rendering.
type BorderStyle int

const (
	BorderNone BorderStyle = iota
	BorderSingle
	BorderDouble
	BorderRounded
	BorderThick
)

// Color represents a terminal color.
type Color struct {
	R, G, B uint8
	Default bool
}

// DefaultColor returns the default terminal color.
func DefaultColor() Color {
	return Color{Default: true}
}

// RGB creates an RGB color.
func RGB(r, g, b uint8) Color {
	return Color{R: r, G: g, B: b, Default: false}
}

// Alignment defines content alignment.
type Alignment int

const (
	AlignStart Alignment = iota
	AlignCenter
	AlignEnd
	AlignStretch
)

// Justification defines content justification.
type Justification int

const (
	JustifyStart Justification = iota
	JustifyCenter
	JustifyEnd
	JustifySpaceBetween
	JustifySpaceAround
	JustifySpaceEvenly
)

// BaseView provides a base implementation for views.
type BaseView struct {
	ViewType   string
	ChildViews []View
	Attributes Attributes
}

// Type returns the view type.
func (v *BaseView) Type() string {
	return v.ViewType
}

// Children returns child views.
func (v *BaseView) Children() []View {
	return v.ChildViews
}

// Attrs returns view attributes.
func (v *BaseView) Attrs() Attributes {
	return v.Attributes
}

// TextView represents a text view.
type TextView struct {
	BaseView
	Content string
}

// Text creates a text view.
func Text(content string, opts ...ViewOption) View {
	v := &TextView{
		BaseView: BaseView{
			ViewType:   "text",
			ChildViews: nil,
			Attributes: Attributes{},
		},
		Content: content,
	}
	for _, opt := range opts {
		opt(&v.Attributes)
	}
	return v
}

// ContainerView represents a container view.
type ContainerView struct {
	BaseView
	Direction Direction
}

// Direction defines layout direction.
type Direction int

const (
	DirectionRow Direction = iota
	DirectionColumn
)

// Row creates a horizontal container.
func Row(children ...View) View {
	return &ContainerView{
		BaseView: BaseView{
			ViewType:   "row",
			ChildViews: children,
			Attributes: Attributes{},
		},
		Direction: DirectionRow,
	}
}

// Column creates a vertical container.
func Column(children ...View) View {
	return &ContainerView{
		BaseView: BaseView{
			ViewType:   "column",
			ChildViews: children,
			Attributes: Attributes{},
		},
		Direction: DirectionColumn,
	}
}

// ViewOption configures view attributes.
type ViewOption func(*Attributes)

// WithID sets the view ID.
func WithID(id string) ViewOption {
	return func(a *Attributes) {
		a.ID = id
	}
}

// WithSize sets width and height.
func WithSize(width, height int) ViewOption {
	return func(a *Attributes) {
		a.Width = width
		a.Height = height
	}
}

// WithPadding sets padding.
func WithPadding(padding Spacing) ViewOption {
	return func(a *Attributes) {
		a.Padding = padding
	}
}

// WithMargin sets margin.
func WithMargin(margin Spacing) ViewOption {
	return func(a *Attributes) {
		a.Margin = margin
	}
}

// WithBorder sets the border style.
func WithBorder(style BorderStyle) ViewOption {
	return func(a *Attributes) {
		a.Border = style
	}
}

// WithBackground sets the background color.
func WithBackground(color Color) ViewOption {
	return func(a *Attributes) {
		a.Background = color
	}
}

// WithForeground sets the foreground color.
func WithForeground(color Color) ViewOption {
	return func(a *Attributes) {
		a.Foreground = color
	}
}

// WithAlign sets alignment.
func WithAlign(align Alignment) ViewOption {
	return func(a *Attributes) {
		a.Align = align
	}
}

// WithGrow sets the grow factor.
func WithGrow(grow int) ViewOption {
	return func(a *Attributes) {
		a.Grow = grow
	}
}

// WithFocus sets focus state.
func WithFocus(focused bool) ViewOption {
	return func(a *Attributes) {
		a.Focused = focused
	}
}
