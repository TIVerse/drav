package vayu

import (
	"context"
)

// Plugin represents a loaded plugin.
type Plugin interface {
	// Metadata returns plugin metadata.
	Metadata() PluginMetadata

	// Init initializes the plugin.
	Init(ctx context.Context) error

	// Start starts the plugin.
	Start(ctx context.Context) error

	// Stop stops the plugin.
	Stop(ctx context.Context) error

	// Capabilities returns the required capabilities.
	Capabilities() Capabilities
}

// PluginMetadata contains plugin metadata.
type PluginMetadata struct {
	Name        string
	Version     string
	Author      string
	Description string
	License     string
	Homepage    string
}

// BasePlugin provides a base implementation of Plugin.
type BasePlugin struct {
	meta PluginMetadata
	caps Capabilities
}

// NewBasePlugin creates a new base plugin.
func NewBasePlugin(meta PluginMetadata, caps Capabilities) *BasePlugin {
	return &BasePlugin{
		meta: meta,
		caps: caps,
	}
}

// Metadata returns plugin metadata.
func (p *BasePlugin) Metadata() PluginMetadata {
	return p.meta
}

// Init initializes the plugin (default no-op).
func (p *BasePlugin) Init(ctx context.Context) error {
	return nil
}

// Start starts the plugin (default no-op).
func (p *BasePlugin) Start(ctx context.Context) error {
	return nil
}

// Stop stops the plugin (default no-op).
func (p *BasePlugin) Stop(ctx context.Context) error {
	return nil
}

// Capabilities returns the required capabilities.
func (p *BasePlugin) Capabilities() Capabilities {
	return p.caps
}

// PluginInfo provides summary information about a plugin.
type PluginInfo struct {
	Name    string
	Version string
	Status  PluginStatus
	Path    string
}

// PluginStatus represents the status of a plugin.
type PluginStatus int

const (
	// PluginStatusUnloaded means the plugin is not loaded.
	PluginStatusUnloaded PluginStatus = iota
	// PluginStatusLoaded means the plugin is loaded but not started.
	PluginStatusLoaded
	// PluginStatusRunning means the plugin is running.
	PluginStatusRunning
	// PluginStatusError means the plugin encountered an error.
	PluginStatusError
)

// String returns the string representation of plugin status.
func (s PluginStatus) String() string {
	switch s {
	case PluginStatusUnloaded:
		return "unloaded"
	case PluginStatusLoaded:
		return "loaded"
	case PluginStatusRunning:
		return "running"
	case PluginStatusError:
		return "error"
	default:
		return "unknown"
	}
}
