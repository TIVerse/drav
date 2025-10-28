package main

import (
	"context"
	"log"
	"log/slog"

	"github.com/TIVerse/drav/pkg/dravya"
	"github.com/TIVerse/drav/pkg/maya"
)

func main() {
	// Create application
	app := dravya.NewApp(
		dravya.WithLogLevel(slog.LevelInfo),
	)

	// Create a simple text component
	root := maya.Text("Hello, DRAV! 🌊\n\nPress Ctrl+C to exit.")

	// Set the root component
	app.SetRoot(maya.Stateless(root))

	// Run the application
	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
