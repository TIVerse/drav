package main

import (
	"context"
	"log"

	"github.com/TIVerse/drav/pkg/dravya"
	"github.com/TIVerse/drav/pkg/maya"
	"github.com/TIVerse/drav/pkg/sri"
)

func main() {
	// Create themes
	darkTheme := sri.DefaultDark()
	lightTheme := sri.DefaultLight()

	log.Printf("Dark theme: %s", darkTheme.Name)
	log.Printf("Light theme: %s", lightTheme.Name)

	app := dravya.NewApp(
		dravya.WithTheme(darkTheme.Name),
	)

	// Create UI with styled text
	root := maya.Column(
		maya.Text("=== Theme Demo ==="),
		maya.Text(""),
		maya.Text("Using dark theme by default"),
		maya.Text(""),
		maya.Text("Available themes:"),
		maya.Text("  - default-dark (active)"),
		maya.Text("  - default-light"),
		maya.Text(""),
		maya.Text("Theme switching coming soon!"),
		maya.Text(""),
		maya.Text("Press [Ctrl+C] to exit"),
	)

	app.SetRoot(maya.Stateless(root))

	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
