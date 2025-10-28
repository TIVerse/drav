package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/TIVerse/drav/pkg/dravya"
)

var (
	version   = "dev"
	buildDate = "unknown"
	gitCommit = "unknown"
)

func main() {
	// Parse flags
	versionFlag := flag.Bool("version", false, "Print version information")
	verboseFlag := flag.Bool("verbose", false, "Enable verbose logging")
	debugFlag := flag.Bool("debug", false, "Enable debug logging")
	fpsFlag := flag.Int("fps", 60, "Target frames per second")
	pprofFlag := flag.Bool("pprof", false, "Enable pprof profiling")
	pprofAddr := flag.String("pprof-addr", "localhost:6060", "pprof server address")

	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Printf("DRAV (द्रव) - Dynamic Reactive Application View\n")
		fmt.Printf("Version:    %s\n", version)
		fmt.Printf("Build Date: %s\n", buildDate)
		fmt.Printf("Git Commit: %s\n", gitCommit)
		os.Exit(0)
	}

	// Configure logging
	logLevel := slog.LevelInfo
	if *verboseFlag {
		logLevel = slog.LevelDebug
	}
	if *debugFlag {
		logLevel = slog.LevelDebug
	}

	// Build app options
	appOpts := []dravya.AppOption{
		dravya.WithLogLevel(logLevel),
		dravya.WithFPSCap(*fpsFlag),
	}

	if *pprofFlag {
		appOpts = append(appOpts, dravya.WithPprof(true, *pprofAddr))
	}

	// Create application
	app := dravya.NewApp(appOpts...)

	// Set up signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("Received interrupt signal, shutting down...")
		cancel()
	}()

	// TODO: Load root component based on CLI args or config
	// For now, just run with no root to demonstrate the framework

	// Run application
	if err := app.Run(ctx); err != nil && err != context.Canceled {
		log.Fatalf("Application error: %v", err)
	}

	log.Println("DRAV shutdown complete")
}
