package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/TIVerse/drav/pkg/dravya"
	"github.com/TIVerse/drav/pkg/maya"
	"github.com/TIVerse/drav/pkg/prana"
)

// Dashboard displays system metrics.
type Dashboard struct {
	cpu    *prana.Observable[float64]
	memory *prana.Observable[float64]
	uptime *prana.Observable[time.Duration]
}

// NewDashboard creates a new dashboard.
func NewDashboard() *Dashboard {
	return &Dashboard{
		cpu:    prana.NewObservable(0.0),
		memory: prana.NewObservable(0.0),
		uptime: prana.NewObservable(time.Duration(0)),
	}
}

// Render renders the dashboard.
func (d *Dashboard) Render(ctx maya.RenderContext) maya.View {
	return maya.Column(
		maya.Text("=== System Dashboard ==="),
		maya.Text(""),
		maya.Text(fmt.Sprintf("CPU:    %.1f%%", d.cpu.Get())),
		maya.Text(fmt.Sprintf("Memory: %.1f%%", d.memory.Get())),
		maya.Text(fmt.Sprintf("Uptime: %s", d.uptime.Get())),
		maya.Text(""),
		maya.Text("Press [Ctrl+C] to exit"),
	)
}

// Update updates the dashboard metrics (simulated).
func (d *Dashboard) Update() {
	// Simulate metrics
	d.cpu.Set(50.0 + float64(time.Now().Unix()%50))
	d.memory.Set(60.0 + float64(time.Now().Unix()%40))
	d.uptime.Update(func(t time.Duration) time.Duration {
		return t + time.Second
	})
}

func main() {
	app := dravya.NewApp()
	dashboard := NewDashboard()

	// Update metrics every second
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			dashboard.Update()
		}
	}()

	app.SetRoot(dashboard)

	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
