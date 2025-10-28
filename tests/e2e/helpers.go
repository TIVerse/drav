package e2e

import (
	"context"
	"time"

	"github.com/TIVerse/drav/pkg/dravya"
)

// RunWithTimeout runs an app with a timeout.
func RunWithTimeout(app *dravya.App, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return app.Run(ctx)
}

// WaitForState waits for a lifecycle state.
func WaitForState(app *dravya.App, targetState dravya.LifecycleState, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if app.Lifecycle().State() == targetState {
			return true
		}
		time.Sleep(10 * time.Millisecond)
	}
	return false
}
