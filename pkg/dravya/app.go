package dravya

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"time"

	"github.com/TIVerse/drav/pkg/maya"
)

// App is the main DRAV application.
type App struct {
	config      *appConfig
	logger      *Logger
	lifecycle   *Lifecycle
	loop        *Loop
	root        Component
	renderer    Renderer
	eventHub    EventHub
	stateStore  StateStore
	cmdRegistry CommandRegistry
	pluginMgr   PluginManager
	themeMgr    ThemeManager
	mu          sync.RWMutex
	startTime   time.Time
	tasks       sync.WaitGroup
} // Component is a renderable UI component (alias for maya.Component).
type Component = maya.Component

// RenderContext provides context for rendering (alias for maya.RenderContext).
type RenderContext = maya.RenderContext

// View represents a virtual UI tree node (alias for maya.View).
type View = maya.View

// Renderer handles screen rendering.
type Renderer interface {
	Init() error
	Render(ctx context.Context, view interface{}) error
	Clear() error
	Shutdown() error
}

// EventHub manages event dispatch.
type EventHub interface {
	Emit(ctx context.Context, event Event) error
	On(eventType string, handler EventHandler) (unsubscribe func())
	Start(ctx context.Context) error
	Stop() error
}

// Event represents a system or user event.
type Event interface {
	Type() string
	Time() time.Time
}

// EventHandler processes events.
type EventHandler func(ctx context.Context, event Event) error

// StateStore manages reactive state.
type StateStore interface {
	Get(key string) (any, bool)
	Set(key string, value any) error
	Watch(key string, callback func(old, new any)) (unwatch func())
}

// CommandRegistry manages commands.
type CommandRegistry interface {
	Register(cmd Command) error
	Execute(ctx context.Context, input string) (Result, error)
	List() []Command
}

// Command represents an executable command.
type Command interface {
	Name() string
	Execute(ctx context.Context, args []string) (Result, error)
}

// Result is the result of command execution.
type Result interface {
	Success() bool
	Message() string
}

// PluginManager manages plugins.
type PluginManager interface {
	Load(path string) error
	Unload(name string) error
	List() []Plugin
}

// Plugin represents a loaded plugin.
type Plugin interface {
	Name() string
	Version() string
}

// ThemeManager manages themes.
type ThemeManager interface {
	Load(name string) error
	Current() Theme
}

// Theme represents a UI theme.
type Theme interface {
	Name() string
}

// NewApp creates a new DRAV application with options.
func NewApp(opts ...AppOption) *App {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	logger := NewLogger(cfg.logLevel, slog.NewTextHandler(cfg.logOutput, &slog.HandlerOptions{
		Level: cfg.logLevel,
	}))

	app := &App{
		config:    cfg,
		logger:    logger,
		lifecycle: NewLifecycle(),
		loop:      NewLoop(cfg.fpsTarget),
	}

	// Register lifecycle hooks
	app.lifecycle.OnState(StateInitializing, app.onInitializing)
	app.lifecycle.OnState(StateRunning, app.onRunning)
	app.lifecycle.OnState(StateShuttingDown, app.onShuttingDown)

	return app
}

// SetRoot sets the root component for rendering (accepts any maya.Component).
func (a *App) SetRoot(root Component) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.root = root
}

// Run starts the application main loop.
func (a *App) Run(ctx context.Context) error {
	a.startTime = time.Now()
	a.logger.Info("Starting DRAV application")

	// Initialize
	if err := a.lifecycle.SetState(ctx, StateInitializing); err != nil {
		return fmt.Errorf("initialization failed: %w", err)
	}

	// Start pprof if enabled
	if a.config.enablePprof {
		go func() {
			a.logger.Info("Starting pprof server", "addr", a.config.pprofAddr)
			if err := http.ListenAndServe(a.config.pprofAddr, nil); err != nil {
				a.logger.Error("pprof server failed", "error", err)
			}
		}()
	}

	// Transition to running
	if err := a.lifecycle.SetState(ctx, StateRunning); err != nil {
		return fmt.Errorf("failed to start: %w", err)
	}

	// Start main loop
	errCh := make(chan error, 1)
	go func() {
		errCh <- a.loop.Start(ctx)
	}()

	// Wait for shutdown signal or context cancellation
	select {
	case <-ctx.Done():
		a.logger.Info("Context cancelled, shutting down")
		return a.Shutdown(ctx.Err())
	case <-a.lifecycle.ShutdownSignal():
		a.logger.Info("Shutdown signal received")
		return a.Shutdown(a.lifecycle.ShutdownError())
	case err := <-errCh:
		if err != nil {
			a.logger.Error("Main loop error", "error", err)
			return a.Shutdown(err)
		}
	}

	return nil
}

// Shutdown gracefully shuts down the application.
func (a *App) Shutdown(err error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.lifecycle.Shutdown(err)

	if stateErr := a.lifecycle.SetState(ctx, StateShuttingDown); stateErr != nil {
		a.logger.Error("Failed to transition to shutdown state", "error", stateErr)
	}

	// Stop the main loop
	a.loop.Stop()

	// Wait for background tasks with timeout
	done := make(chan struct{})
	go func() {
		a.tasks.Wait()
		close(done)
	}()

	select {
	case <-done:
		a.logger.Info("All tasks completed")
	case <-ctx.Done():
		a.logger.Warn("Shutdown timeout, some tasks may not have completed")
	}

	// Transition to terminated
	if stateErr := a.lifecycle.SetState(ctx, StateTerminated); stateErr != nil {
		a.logger.Error("Failed to transition to terminated state", "error", stateErr)
	}

	uptime := time.Since(a.startTime)
	stats := a.loop.Stats()
	a.logger.Info("Application terminated",
		"uptime", uptime,
		"frames", stats.FrameCount,
		"dropped_frames", stats.DroppedFrames,
		"error", err,
	)

	return err
}

// Logger returns the application logger.
func (a *App) Logger() *Logger {
	return a.logger
}

// Lifecycle returns the lifecycle manager.
func (a *App) Lifecycle() *Lifecycle {
	return a.lifecycle
}

// onInitializing is called during initialization.
func (a *App) onInitializing(ctx context.Context) error {
	a.logger.Info("Initializing application")

	// Initialize renderer if set
	if a.renderer != nil {
		if err := a.renderer.Init(); err != nil {
			return fmt.Errorf("renderer init failed: %w", err)
		}
	}

	// Start event hub if set
	if a.eventHub != nil {
		a.tasks.Add(1)
		go func() {
			defer a.tasks.Done()
			if err := a.eventHub.Start(ctx); err != nil {
				a.logger.Error("Event hub error", "error", err)
			}
		}()
	}

	return nil
}

// onRunning is called when entering running state.
func (a *App) onRunning(ctx context.Context) error {
	a.logger.Info("Application running")

	// Register frame callback for rendering
	a.loop.OnFrame(a.renderFrame)

	return nil
}

// onShuttingDown is called during shutdown.
func (a *App) onShuttingDown(ctx context.Context) error {
	a.logger.Info("Shutting down application")

	// Stop event hub
	if a.eventHub != nil {
		if err := a.eventHub.Stop(); err != nil {
			a.logger.Error("Failed to stop event hub", "error", err)
		}
	}

	// Shutdown renderer
	if a.renderer != nil {
		if err := a.renderer.Shutdown(); err != nil {
			a.logger.Error("Failed to shutdown renderer", "error", err)
		}
	}

	return nil
}

// renderFrame is called each frame to render the UI.
func (a *App) renderFrame(ctx context.Context, frameTime time.Time, delta time.Duration) error {
	a.mu.RLock()
	root := a.root
	renderer := a.renderer
	a.mu.RUnlock()

	if root == nil || renderer == nil {
		return nil
	}

	// Render the root component
	// For now, we skip actual rendering as it requires maya integration
	// This will be implemented when we create the maya package

	return nil
}

// Stats returns application statistics.
func (a *App) Stats() AppStats {
	loopStats := a.loop.Stats()
	return AppStats{
		Uptime:        time.Since(a.startTime),
		State:         a.lifecycle.State().String(),
		FrameCount:    loopStats.FrameCount,
		DroppedFrames: loopStats.DroppedFrames,
		FPS:           a.config.fpsTarget,
	}
}

// AppStats contains application statistics.
type AppStats struct {
	Uptime        time.Duration
	State         string
	FrameCount    uint64
	DroppedFrames uint64
	FPS           int
}

// ErrShutdown is returned when the application is shutting down.
var ErrShutdown = errors.New("application is shutting down")
