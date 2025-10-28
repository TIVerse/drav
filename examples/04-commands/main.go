package main

import (
	"context"
	"fmt"
	"log"

	"github.com/TIVerse/drav/pkg/dravya"
	"github.com/TIVerse/drav/pkg/maya"
	"github.com/TIVerse/drav/pkg/vak"
)

func main() {
	app := dravya.NewApp()
	registry := vak.NewRegistry()

	// Register commands
	registry.Register(vak.Command{
		Name:    "hello",
		Summary: "Print a greeting",
		Usage:   "hello [name]",
		Execute: func(ctx context.Context, args []string) (vak.Result, error) {
			name := "World"
			if len(args) > 0 {
				name = args[0]
			}
			return vak.SuccessResult(fmt.Sprintf("Hello, %s!", name)), nil
		},
	})

	registry.Register(vak.Command{
		Name:    "echo",
		Summary: "Echo back the arguments",
		Usage:   "echo <text>",
		Execute: func(ctx context.Context, args []string) (vak.Result, error) {
			if len(args) == 0 {
				return vak.ErrorResult("no arguments provided"), nil
			}
			return vak.SuccessResult(fmt.Sprintf("Echo: %v", args)), nil
		},
	})

	// Create UI
	root := maya.Column(
		maya.Text("=== Command Demo ==="),
		maya.Text(""),
		maya.Text("Available commands:"),
		maya.Text("  - hello [name]"),
		maya.Text("  - echo <text>"),
		maya.Text(""),
		maya.Text("Command palette integration coming soon!"),
		maya.Text(""),
		maya.Text("Press [Ctrl+C] to exit"),
	)

	app.SetRoot(maya.Stateless(root))

	// Test command execution
	result, err := registry.Execute(context.Background(), "hello DRAV")
	if err != nil {
		log.Printf("Command error: %v", err)
	} else {
		log.Printf("Command result: %s", result.Message())
	}

	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
