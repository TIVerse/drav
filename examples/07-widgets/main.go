package main

import (
	"context"
	"log"

	"github.com/TIVerse/drav/pkg/dravya"
	"github.com/TIVerse/drav/pkg/maya"
)

func main() {
	app := dravya.NewApp()

	// Create a layout with multiple widgets
	root := maya.Column(
		maya.Text("=== Widget Gallery ==="),
		maya.Text(""),
		maya.Row(
			maya.Text("Left Panel"),
			maya.Text(" | "),
			maya.Text("Right Panel"),
		),
		maya.Text(""),
		maya.Text("Available widgets:"),
		maya.Text("  ✓ Text"),
		maya.Text("  ✓ Row/Column layouts"),
		maya.Text("  - List (coming soon)"),
		maya.Text("  - Table (coming soon)"),
		maya.Text("  - Input (coming soon)"),
		maya.Text("  - Button (coming soon)"),
		maya.Text(""),
		maya.Text("Press [Ctrl+C] to exit"),
	)

	app.SetRoot(maya.Stateless(root))

	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
