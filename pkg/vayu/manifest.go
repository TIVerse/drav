package vayu

import (
	"encoding/json"
	"fmt"
	"os"
)

// Manifest describes a plugin's metadata and requirements.
type Manifest struct {
	Name         string       `json:"name"`
	Version      string       `json:"version"`
	Author       string       `json:"author"`
	Description  string       `json:"description"`
	License      string       `json:"license"`
	Homepage     string       `json:"homepage"`
	EntryPoint   string       `json:"entry_point"`
	Type         string       `json:"type"` // "wasm", "process", "go-plugin"
	Capabilities Capabilities `json:"capabilities"`
	Dependencies []string     `json:"dependencies"`
}

// LoadManifest loads a manifest from a JSON file.
func LoadManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	return &manifest, nil
}

// Validate validates the manifest.
func (m *Manifest) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("manifest must have a name")
	}
	if m.Version == "" {
		return fmt.Errorf("manifest must have a version")
	}
	if m.EntryPoint == "" {
		return fmt.Errorf("manifest must have an entry point")
	}
	if m.Type == "" {
		return fmt.Errorf("manifest must specify a type")
	}
	if m.Type != "wasm" && m.Type != "process" && m.Type != "go-plugin" {
		return fmt.Errorf("invalid plugin type: %s", m.Type)
	}
	return nil
}

// ToMetadata converts manifest to plugin metadata.
func (m *Manifest) ToMetadata() PluginMetadata {
	return PluginMetadata{
		Name:        m.Name,
		Version:     m.Version,
		Author:      m.Author,
		Description: m.Description,
		License:     m.License,
		Homepage:    m.Homepage,
	}
}
