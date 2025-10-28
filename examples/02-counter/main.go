package main

import (
	"context"
	"fmt"
	"log"

	"github.com/TIVerse/drav/pkg/dravya"
	"github.com/TIVerse/drav/pkg/maya"
	"github.com/TIVerse/drav/pkg/prana"
)

// Counter is a reactive counter component.
type Counter struct {
	count *prana.Observable[int]
}

// NewCounter creates a new counter.
func NewCounter() *Counter {
	return &Counter{
		count: prana.NewObservable(0),
	}
}

// Render renders the counter.
func (c *Counter) Render(ctx maya.RenderContext) maya.View {
	return maya.Column(
		maya.Text(fmt.Sprintf("Count: %d", c.count.Get())),
		maya.Text(""),
		maya.Text("Press [+] to increment"),
		maya.Text("Press [-] to decrement"),
		maya.Text("Press [r] to reset"),
		maya.Text("Press [Ctrl+C] to exit"),
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

func main() {
	app := dravya.NewApp()
	counter := NewCounter()

	// Watch for changes and trigger re-render
	counter.count.Watch(func(old, new int) {
		log.Printf("Count changed: %d -> %d", old, new)
	})

	app.SetRoot(counter)

	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
