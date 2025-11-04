package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/TIVerse/drav/pkg/agni"
	"github.com/TIVerse/drav/pkg/dravya"
	"github.com/TIVerse/drav/pkg/maya"
	"github.com/TIVerse/drav/pkg/prana"
)

// Counter is a reactive counter component.
type Counter struct {
	count *prana.Observable[int]
	app   *dravya.App
}

// NewCounter creates a new counter.
func NewCounter(app *dravya.App) *Counter {
	return &Counter{
		count: prana.NewObservable(0),
		app:   app,
	}
}

// Render renders the counter.
func (c *Counter) Render(ctx maya.RenderContext) maya.View {
	return maya.Column(
		maya.Text("╔════════════════════════════════════╗"),
		maya.Text(fmt.Sprintf("║  Counter: %-20d  ║", c.count.Get())),
		maya.Text("╚════════════════════════════════════╝"),
		maya.Text(""),
		maya.Text("  [+] or [=]  Increment"),
		maya.Text("  [-]         Decrement"),
		maya.Text("  [r]         Reset"),
		maya.Text("  [Ctrl+C]    Exit"),
		maya.Text(""),
		maya.Text("This example demonstrates reactive state!"),
		maya.Text("The UI updates automatically when the count changes."),
	)
}

// Increment increments the counter.
func (c *Counter) Increment() {
	c.count.Update(func(n int) int { return n + 1 })
}

// Decrement decrements the counter.
func (c *Counter) Decrement() {
	c.count.Update(func(n int) int { return n - 1 })
}

// Reset resets the counter.
func (c *Counter) Reset() {
	c.count.Set(0)
}

// SetupEventHandlers registers event handlers for the counter.
func (c *Counter) SetupEventHandlers(eventHub dravya.EventHub) {
	// Handle key events
	eventHub.On(agni.EventTypeKey, func(ctx context.Context, event dravya.Event) error {
		keyEvent, ok := event.(*agni.KeyEvent)
		if !ok {
			return nil
		}

		switch {
		case keyEvent.Key == agni.KeyRune && (keyEvent.Rune == '+' || keyEvent.Rune == '='):
			c.Increment()
			c.app.Logger().Info("Counter incremented", "count", c.count.Get())
		case keyEvent.Key == agni.KeyRune && keyEvent.Rune == '-':
			c.Decrement()
			c.app.Logger().Info("Counter decremented", "count", c.count.Get())
		case keyEvent.Key == agni.KeyRune && (keyEvent.Rune == 'r' || keyEvent.Rune == 'R'):
			c.Reset()
			c.app.Logger().Info("Counter reset")
		}

		return nil
	})
}

func main() {
	// Create app with options
	app := dravya.NewApp(
		dravya.WithLogLevel(slog.LevelInfo),
	)

	// Create counter
	counter := NewCounter(app)

	// Watch for changes (demonstrates reactivity)
	counter.count.Watch(func(old, new int) {
		app.Logger().Info("Observable changed", "old", old, "new", new)
	})

	// Set up event handlers after app is ready
	app.OnReady(func() {
		eventHub := app.EventHub()
		if eventHub != nil {
			counter.SetupEventHandlers(eventHub)
			app.Logger().Info("Counter event handlers registered")
		}
	})

	// Set root component
	app.SetRoot(counter)

	// Run the application
	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
