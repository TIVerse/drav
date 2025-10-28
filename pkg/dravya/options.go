package dravya

import (
	"io"
	"log/slog"
	"os"
	"time"
)

// AppOption configures the DRAV application.
type AppOption func(*appConfig)

// appConfig holds internal configuration for the App.
type appConfig struct {
	logLevel       slog.Level
	logOutput      io.Writer
	fpsTarget      int
	enableMetrics  bool
	enablePprof    bool
	pprofAddr      string
	inputMode      InputMode
	mouseEnabled   bool
	pasteEnabled   bool
	altScreenMode  bool
	themeName      string
	pluginDirs     []string
	debugOverlay   bool
	maxConcurrency int
}

// InputMode defines how the terminal receives input.
type InputMode int

const (
	// InputModeCooked is standard line-buffered input.
	InputModeCooked InputMode = iota
	// InputModeRaw is character-by-character unbuffered input.
	InputModeRaw
)

// defaultConfig returns the default application configuration.
func defaultConfig() *appConfig {
	return &appConfig{
		logLevel:       slog.LevelInfo,
		logOutput:      os.Stderr,
		fpsTarget:      60,
		enableMetrics:  false,
		enablePprof:    false,
		pprofAddr:      "localhost:6060",
		inputMode:      InputModeRaw,
		mouseEnabled:   true,
		pasteEnabled:   true,
		altScreenMode:  true,
		themeName:      "default-dark",
		pluginDirs:     []string{"./plugins"},
		debugOverlay:   false,
		maxConcurrency: 10,
	}
}

// WithLogLevel sets the logging level.
func WithLogLevel(level slog.Level) AppOption {
	return func(c *appConfig) {
		c.logLevel = level
	}
}

// WithLogOutput sets the log output writer.
func WithLogOutput(w io.Writer) AppOption {
	return func(c *appConfig) {
		c.logOutput = w
	}
}

// WithFPSCap sets the target frames per second cap.
func WithFPSCap(fps int) AppOption {
	return func(c *appConfig) {
		if fps > 0 && fps <= 240 {
			c.fpsTarget = fps
		}
	}
}

// WithMetrics enables metrics collection.
func WithMetrics(enabled bool) AppOption {
	return func(c *appConfig) {
		c.enableMetrics = enabled
	}
}

// WithPprof enables pprof profiling endpoints.
func WithPprof(enabled bool, addr string) AppOption {
	return func(c *appConfig) {
		c.enablePprof = enabled
		if addr != "" {
			c.pprofAddr = addr
		}
	}
}

// WithInputMode sets the terminal input mode.
func WithInputMode(mode InputMode) AppOption {
	return func(c *appConfig) {
		c.inputMode = mode
	}
}

// WithMouse enables or disables mouse input.
func WithMouse(enabled bool) AppOption {
	return func(c *appConfig) {
		c.mouseEnabled = enabled
	}
}

// WithPaste enables or disables bracketed paste mode.
func WithPaste(enabled bool) AppOption {
	return func(c *appConfig) {
		c.pasteEnabled = enabled
	}
}

// WithAltScreen enables or disables alternate screen mode.
func WithAltScreen(enabled bool) AppOption {
	return func(c *appConfig) {
		c.altScreenMode = enabled
	}
}

// WithTheme sets the theme name.
func WithTheme(name string) AppOption {
	return func(c *appConfig) {
		c.themeName = name
	}
}

// WithPluginDirs sets the plugin search directories.
func WithPluginDirs(dirs ...string) AppOption {
	return func(c *appConfig) {
		c.pluginDirs = dirs
	}
}

// WithDebugOverlay enables the debug overlay (toggle with F12).
func WithDebugOverlay(enabled bool) AppOption {
	return func(c *appConfig) {
		c.debugOverlay = enabled
	}
}

// WithMaxConcurrency sets the maximum number of concurrent background tasks.
func WithMaxConcurrency(n int) AppOption {
	return func(c *appConfig) {
		if n > 0 {
			c.maxConcurrency = n
		}
	}
}

// frameDuration returns the target duration per frame.
func (c *appConfig) frameDuration() time.Duration {
	return time.Second / time.Duration(c.fpsTarget)
}
