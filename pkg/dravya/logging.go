package dravya

import (
	"context"
	"log/slog"
	"os"
)

// Logger is the structured logger for DRAV.
type Logger struct {
	*slog.Logger
}

// NewLogger creates a new structured logger.
func NewLogger(level slog.Level, handler slog.Handler) *Logger {
	if handler == nil {
		handler = slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: level,
		})
	}
	return &Logger{
		Logger: slog.New(handler),
	}
}

// WithContext returns a logger with context values.
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{
		Logger: l.Logger.With("trace_id", ctx.Value("trace_id")),
	}
}

// Component returns a logger scoped to a specific component.
func (l *Logger) Component(name string) *Logger {
	return &Logger{
		Logger: l.Logger.With("component", name),
	}
}

// Performance logs performance-related information.
func (l *Logger) Performance(msg string, args ...any) {
	l.Logger.Info(msg, append([]any{"type", "performance"}, args...)...)
}

// Metric logs metric-related information.
func (l *Logger) Metric(msg string, args ...any) {
	l.Logger.Debug(msg, append([]any{"type", "metric"}, args...)...)
}
