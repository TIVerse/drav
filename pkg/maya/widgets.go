package maya

// PanelView represents a panel with a border and title.
type PanelView struct {
	BaseView
	Title   string
	Focused bool
}

// Panel creates a panel view with optional title.
func Panel(title string, children []View, opts ...ViewOption) View {
	attrs := Attributes{
		Border: BorderSingle,
	}
	for _, opt := range opts {
		opt(&attrs)
	}
	
	return &PanelView{
		BaseView: BaseView{
			ViewType:   "panel",
			ChildViews: children,
			Attributes: attrs,
		},
		Title:   title,
		Focused: attrs.Focused,
	}
}

// ButtonView represents a clickable button.
type ButtonView struct {
	BaseView
	Label    string
	OnClick  func()
	Disabled bool
	Focused  bool
}

// Button creates a button view.
func Button(label string, onClick func(), opts ...ViewOption) View {
	attrs := Attributes{
		Border: BorderSingle,
	}
	for _, opt := range opts {
		opt(&attrs)
	}
	
	return &ButtonView{
		BaseView: BaseView{
			ViewType:   "button",
			ChildViews: nil,
			Attributes: attrs,
		},
		Label:    label,
		OnClick:  onClick,
		Disabled: attrs.Disabled,
		Focused:  attrs.Focused,
	}
}

// InputView represents a text input field.
type InputView struct {
	BaseView
	Value       string
	Placeholder string
	OnChange    func(string)
	OnSubmit    func(string)
	Multiline   bool
	MaxLength   int
	CursorPos   int
	Focused     bool
	Disabled    bool
}

// Input creates a text input view.
func Input(value string, opts ...ViewOption) View {
	attrs := Attributes{
		Border: BorderSingle,
	}
	for _, opt := range opts {
		opt(&attrs)
	}
	
	return &InputView{
		BaseView: BaseView{
			ViewType:   "input",
			ChildViews: nil,
			Attributes: attrs,
		},
		Value:     value,
		CursorPos: len(value),
		Focused:   attrs.Focused,
		Disabled:  attrs.Disabled,
	}
}

// WithPlaceholder sets the input placeholder.
func WithPlaceholder(placeholder string) ViewOption {
	return func(a *Attributes) {
		// Store in a custom way - we'll handle this in the Input constructor
		// For now, this is a marker
	}
}

// WithOnChange sets the onChange handler.
func WithOnChange(handler func(string)) ViewOption {
	return func(a *Attributes) {
		// Custom attribute handling
	}
}

// WithOnSubmit sets the onSubmit handler.
func WithOnSubmit(handler func(string)) ViewOption {
	return func(a *Attributes) {
		// Custom attribute handling
	}
}

// ListItem represents an item in a list.
type ListItem struct {
	Label    string
	Value    interface{}
	Selected bool
	Disabled bool
}

// ListView represents a list of items.
type ListView struct {
	BaseView
	Items        []ListItem
	SelectedIdx  int
	OnSelect     func(int, ListItem)
	Multiselect  bool
	Searchable   bool
	VirtualScroll bool
}

// List creates a list view.
func List(items []ListItem, opts ...ViewOption) View {
	attrs := Attributes{
		Border: BorderSingle,
	}
	for _, opt := range opts {
		opt(&attrs)
	}
	
	return &ListView{
		BaseView: BaseView{
			ViewType:   "list",
			ChildViews: nil,
			Attributes: attrs,
		},
		Items:       items,
		SelectedIdx: -1,
	}
}

// WithItems sets list items.
func WithItems(items []ListItem) ViewOption {
	return func(a *Attributes) {
		// Custom attribute handling
	}
}

// TableColumn defines a table column.
type TableColumn struct {
	Header string
	Width  int
	Align  Alignment
}

// TableRow represents a row in a table.
type TableRow struct {
	Cells    []string
	Selected bool
	Disabled bool
}

// TableView represents a table.
type TableView struct {
	BaseView
	Columns      []TableColumn
	Rows         []TableRow
	SelectedRow  int
	Sortable     bool
	SortColumn   int
	SortAsc      bool
	OnRowSelect  func(int, TableRow)
}

// Table creates a table view.
func Table(columns []TableColumn, rows []TableRow, opts ...ViewOption) View {
	attrs := Attributes{
		Border: BorderSingle,
	}
	for _, opt := range opts {
		opt(&attrs)
	}
	
	return &TableView{
		BaseView: BaseView{
			ViewType:   "table",
			ChildViews: nil,
			Attributes: attrs,
		},
		Columns:     columns,
		Rows:        rows,
		SelectedRow: -1,
		SortAsc:     true,
	}
}

// TabItem represents a tab.
type TabItem struct {
	Label   string
	Content View
	Disabled bool
}

// TabsView represents a tabbed interface.
type TabsView struct {
	BaseView
	Tabs      []TabItem
	ActiveTab int
	OnChange  func(int)
}

// Tabs creates a tabs view.
func Tabs(tabs []TabItem, opts ...ViewOption) View {
	attrs := Attributes{}
	for _, opt := range opts {
		opt(&attrs)
	}
	
	return &TabsView{
		BaseView: BaseView{
			ViewType:   "tabs",
			ChildViews: nil,
			Attributes: attrs,
		},
		Tabs:      tabs,
		ActiveTab: 0,
	}
}

// ModalView represents a modal dialog.
type ModalView struct {
	BaseView
	Title     string
	Content   View
	Buttons   []View
	OnClose   func()
	Closeable bool
	Overlay   bool
}

// Modal creates a modal dialog.
func Modal(title string, content View, opts ...ViewOption) View {
	attrs := Attributes{
		Border: BorderDouble,
	}
	for _, opt := range opts {
		opt(&attrs)
	}
	
	return &ModalView{
		BaseView: BaseView{
			ViewType:   "modal",
			ChildViews: []View{content},
			Attributes: attrs,
		},
		Title:     title,
		Content:   content,
		Closeable: true,
		Overlay:   true,
	}
}

// WithTitle sets a title.
func WithTitle(title string) ViewOption {
	return func(a *Attributes) {
		// Custom attribute handling
	}
}

// WithCloseable sets whether a modal is closeable.
func WithCloseable(closeable bool) ViewOption {
	return func(a *Attributes) {
		// Custom attribute handling
	}
}

// WithOverlay sets whether to show an overlay.
func WithOverlay(overlay bool) ViewOption {
	return func(a *Attributes) {
		// Custom attribute handling
	}
}
