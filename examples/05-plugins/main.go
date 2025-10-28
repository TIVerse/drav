package main

import (
	"context"
	"log"

	"github.com/TIVerse/drav/pkg/dravya"
	"github.com/TIVerse/drav/pkg/maya"
	"github.com/TIVerse/drav/pkg/vayu"
)

func main() {
	app := dravya.NewApp()

	// Define plugin capabilities
	caps := vayu.Capabilities{
		Filesystem: vayu.FSCapability{
			Read:  []string{"/tmp"},
			Write: []string{"/tmp/output"},
		},
		Network: vayu.NetworkCapability{
			AllowedDomains: []string{"api.example.com"},
		},
	}

	// Create plugin loader (WASM preferred)
	loader := vayu.NewWASMLoader()

	// Load plugin (this is a stub for now)
	plugin, err := loader.Load("./plugins/sample.wasm", caps)
	if err != nil {
		log.Printf("Failed to load plugin: %v", err)
	} else {
		log.Printf("Plugin loaded: %s", plugin.Metadata().Name)
	}

	// Create UI
	root := maya.Column(
		maya.Text("=== Plugin System Demo ==="),
		maya.Text(""),
		maya.Text("Plugin capabilities:"),
		maya.Text("  - Filesystem: /tmp (read), /tmp/output (write)"),
		maya.Text("  - Network: api.example.com"),
		maya.Text(""),
		maya.Text("WASM plugin loader initialized"),
		maya.Text(""),
		maya.Text("Press [Ctrl+C] to exit"),
	)

	app.SetRoot(maya.Stateless(root))

	if err := app.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
