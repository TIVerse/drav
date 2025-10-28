package vayu

import (
	"context"
	"fmt"
	"os/exec"
)

// ProcessLoader loads plugins as separate processes.
type ProcessLoader struct {
	capChecker *CapabilityChecker
}

// NewProcessLoader creates a new process loader.
func NewProcessLoader() *ProcessLoader {
	return &ProcessLoader{}
}

// Load loads a plugin as a separate process.
func (pl *ProcessLoader) Load(path string, caps Capabilities) (Plugin, error) {
	pl.capChecker = NewCapabilityChecker(caps)

	return &processPlugin{
		path: path,
		caps: caps,
		meta: PluginMetadata{
			Name:    "process-plugin",
			Version: "0.1.0",
		},
	}, nil
}

// processPlugin runs a plugin in a separate process.
type processPlugin struct {
	path   string
	caps   Capabilities
	meta   PluginMetadata
	status PluginStatus
	cmd    *exec.Cmd
}

// Metadata returns plugin metadata.
func (p *processPlugin) Metadata() PluginMetadata {
	return p.meta
}

// Init initializes the plugin.
func (p *processPlugin) Init(ctx context.Context) error {
	p.status = PluginStatusLoaded
	return nil
}

// Start starts the plugin process.
func (p *processPlugin) Start(ctx context.Context) error {
	p.cmd = exec.CommandContext(ctx, p.path)

	// TODO: Set up IPC (stdin/stdout or socket)
	// TODO: Apply resource limits

	if err := p.cmd.Start(); err != nil {
		p.status = PluginStatusError
		return fmt.Errorf("failed to start plugin process: %w", err)
	}

	p.status = PluginStatusRunning
	return nil
}

// Stop stops the plugin process.
func (p *processPlugin) Stop(ctx context.Context) error {
	if p.cmd != nil && p.cmd.Process != nil {
		if err := p.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill plugin process: %w", err)
		}
	}
	p.status = PluginStatusUnloaded
	return nil
}

// Capabilities returns the plugin's capabilities.
func (p *processPlugin) Capabilities() Capabilities {
	return p.caps
}
