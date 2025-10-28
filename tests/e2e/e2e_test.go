package e2e

import (
	"context"
	"testing"
	"time"

	"github.com/TIVerse/drav/pkg/dravya"
	"github.com/TIVerse/drav/pkg/maya"
)

func TestBasicAppLifecycle(t *testing.T) {
	app := dravya.NewApp()

	// Set a simple root component
	root := maya.Text("Test")
	app.SetRoot(maya.Stateless(root))

	// Run app in background
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run(ctx)
	}()

	// Wait for context timeout or error
	select {
	case err := <-errCh:
		if err != nil && err != context.DeadlineExceeded {
			t.Fatalf("App run failed: %v", err)
		}
	case <-ctx.Done():
		// Expected timeout
	}
}

func TestAppShutdown(t *testing.T) {
	app := dravya.NewApp()

	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run(ctx)
	}()

	// Give it time to start
	time.Sleep(100 * time.Millisecond)

	// Cancel context
	cancel()

	// Wait for shutdown
	select {
	case err := <-errCh:
		if err != nil && err != context.Canceled {
			t.Fatalf("Unexpected error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Shutdown timeout")
	}
}
