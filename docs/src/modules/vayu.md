# Vayu - Plugin System

**Package:** `github.com/TIVerse/drav/pkg/vayu`

## Overview

Vayu (Sanskrit: वायु, "wind, carrier") is DRAV's plugin system. It enables loading external functionality via WASM modules or separate processes with capability-based security.

## Key Concepts

### Plugins

Plugins extend DRAV functionality without modifying core code:

- **WASM Plugins**: Sandboxed WebAssembly modules
- **Process Plugins**: Separate OS processes with IPC
- **Native Plugins**: Go plugins (limited support)

### Capabilities

Fine-grained permission system:

```go
type Capability int

const (
    CapFileRead      Capability = 1 << iota  // Read files
    CapFileWrite                              // Write files
    CapNetwork                                // Network access
    CapProcess                                // Spawn processes
    CapUI                                     // UI rendering
    CapState                                  // Access app state
)
```

### Sandboxing

Plugins run in restricted environments:
- WASM: Complete memory isolation
- Process: OS-level isolation
- Capabilities: Explicit permission grants

## Core API

### Plugin Interface

```go
type Plugin interface {
    Name() string
    Version() string
    Init(ctx context.Context) error
    Start(ctx context.Context) error
    Stop() error
}
```

### Loading Plugins

```go
manager := vayu.NewManager()

// Load WASM plugin
plugin, err := manager.LoadWASM("./plugins/my-plugin.wasm", vayu.Capabilities{
    FileRead: true,
    UI: true,
})

// Load process plugin
plugin, err := manager.LoadProcess("./plugins/my-plugin", vayu.Capabilities{
    Network: true,
})
```

### Plugin Manifest

Plugins declare metadata in manifest:

```yaml
# plugin.yaml
name: syntax-highlighter
version: 1.0.0
description: Syntax highlighting for code
author: DRAV Team
capabilities:
  - ui
  - file_read
entry: plugin.wasm
```

## WASM Plugins

### Creating a WASM Plugin

```go
// plugin/main.go
package main

//export init
func init() int32 {
    // Initialize plugin
    return 0
}

//export highlight
func highlight(code *byte, lang *byte) *byte {
    // Highlight code
    highlighted := syntaxHighlight(toString(code), toString(lang))
    return toBytes(highlighted)
}

func main() {}
```

Build:

```bash
GOOS=wasip1 GOARCH=wasm go build -o plugin.wasm plugin/main.go
```

### Host Functions

DRAV provides host functions to plugins:

```go
// Available to WASM plugins
drav.log(message)           // Logging
drav.readFile(path)         // File I/O (if permitted)
drav.httpGet(url)           // HTTP requests (if permitted)
drav.renderUI(component)    // UI rendering (if permitted)
```

## Process Plugins

### Creating a Process Plugin

```go
// plugin/main.go
package main

import (
    "github.com/TIVerse/drav/pkg/vayu/ipc"
)

func main() {
    plugin := &MyPlugin{}
    
    server := ipc.NewServer(plugin)
    server.Listen()
}

type MyPlugin struct{}

func (p *MyPlugin) Name() string {
    return "my-plugin"
}

func (p *MyPlugin) Execute(method string, args []any) (any, error) {
    switch method {
    case "process":
        return p.process(args[0].(string)), nil
    default:
        return nil, fmt.Errorf("unknown method: %s", method)
    }
}
```

### IPC Communication

Plugins communicate via JSON-RPC:

```json
// Request
{
    "jsonrpc": "2.0",
    "method": "process",
    "params": ["data"],
    "id": 1
}

// Response
{
    "jsonrpc": "2.0",
    "result": "processed data",
    "id": 1
}
```

## Capability System

### Defining Capabilities

```go
caps := vayu.Capabilities{
    FileRead:  true,
    FileWrite: false,  // No write access
    Network:   true,
    UI:        true,
}

plugin, err := manager.LoadWASM("plugin.wasm", caps)
```

### Checking Capabilities

```go
if plugin.HasCapability(vayu.CapFileRead) {
    content, err := plugin.ReadFile("/path/to/file")
}
```

### Runtime Enforcement

Capabilities are enforced by the runtime:

```go
// Plugin tries to write file without permission
err := plugin.WriteFile("/path", data)
// Returns: ErrPermissionDenied
```

## Plugin Manager

### Managing Plugins

```go
manager := vayu.NewManager()

// Load plugin
plugin, err := manager.Load("my-plugin", vayu.LoadOptions{
    Type:         vayu.PluginTypeWASM,
    Path:         "./plugins/my-plugin.wasm",
    Capabilities: caps,
})

// List loaded plugins
plugins := manager.List()

// Get specific plugin
plugin := manager.Get("my-plugin")

// Unload plugin
manager.Unload("my-plugin")
```

### Plugin Lifecycle

```go
// 1. Load: Read and validate plugin
plugin, err := manager.Load("my-plugin", opts)

// 2. Init: Initialize plugin resources
err = plugin.Init(ctx)

// 3. Start: Begin plugin execution
err = plugin.Start(ctx)

// 4. Stop: Cleanup and shutdown
err = plugin.Stop()

// 5. Unload: Remove from manager
manager.Unload("my-plugin")
```

## Patterns

### Plugin Registry

```go
type PluginRegistry struct {
    manager *vayu.Manager
    plugins map[string]Plugin
}

func (r *PluginRegistry) LoadAll(dir string) error {
    files, err := os.ReadDir(dir)
    if err != nil {
        return err
    }
    
    for _, file := range files {
        if filepath.Ext(file.Name()) == ".wasm" {
            plugin, err := r.manager.LoadWASM(
                filepath.Join(dir, file.Name()),
                defaultCapabilities,
            )
            if err != nil {
                log.Printf("Failed to load %s: %v", file.Name(), err)
                continue
            }
            
            r.plugins[plugin.Name()] = plugin
        }
    }
    
    return nil
}
```

### Plugin Events

```go
type PluginEventHandler struct {
    eventHub dravya.EventHub
    plugins  []Plugin
}

func (h *PluginEventHandler) register() {
    h.eventHub.On("document-save", func(ctx context.Context, event Event) error {
        // Notify all plugins of save event
        for _, plugin := range h.plugins {
            if err := plugin.OnDocumentSave(ctx, event); err != nil {
                log.Printf("Plugin %s error: %v", plugin.Name(), err)
            }
        }
        return nil
    })
}
```

### Plugin Configuration

```go
type PluginConfig struct {
    Enabled      bool
    Capabilities vayu.Capabilities
    Settings     map[string]any
}

func loadPluginConfig(name string) (*PluginConfig, error) {
    data, err := os.ReadFile(fmt.Sprintf("plugins/%s/config.json", name))
    if err != nil {
        return nil, err
    }
    
    var config PluginConfig
    if err := json.Unmarshal(data, &config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

## Security

### Sandboxing

WASM plugins run in complete isolation:

- No access to host memory
- No direct system calls
- All I/O via host functions
- CPU and memory limits enforced

### Capability Validation

```go
func validateCapabilities(requested, allowed vayu.Capabilities) error {
    if requested.FileWrite && !allowed.FileWrite {
        return fmt.Errorf("file write not permitted")
    }
    
    if requested.Network && !allowed.Network {
        return fmt.Errorf("network access not permitted")
    }
    
    return nil
}
```

### Resource Limits

```go
opts := vayu.LoadOptions{
    Type: vayu.PluginTypeWASM,
    Limits: vayu.ResourceLimits{
        MaxMemory:     100 * 1024 * 1024,  // 100 MB
        MaxCPUTime:    time.Second * 5,     // 5 seconds
        MaxFileSize:   10 * 1024 * 1024,    // 10 MB
        MaxNetworkOps: 100,                 // 100 requests
    },
}
```

## Best Practices

### 1. Minimal Capabilities

Grant only required permissions:

```go
// Good - minimal caps
caps := vayu.Capabilities{
    FileRead: true,
}

// Bad - excessive caps
caps := vayu.Capabilities{
    FileRead:  true,
    FileWrite: true,
    Network:   true,
    Process:   true,
}
```

### 2. Error Handling

Handle plugin failures gracefully:

```go
plugin, err := manager.Load("my-plugin", opts)
if err != nil {
    log.Printf("Failed to load plugin: %v", err)
    // Continue without plugin
    return nil
}

if err := plugin.Start(ctx); err != nil {
    log.Printf("Failed to start plugin: %v", err)
    manager.Unload(plugin.Name())
}
```

### 3. Timeout Protection

Set timeouts for plugin operations:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := plugin.Execute(ctx, "process", data)
if err == context.DeadlineExceeded {
    log.Printf("Plugin timeout")
    return nil
}
```

### 4. Version Compatibility

Check plugin API version:

```go
if plugin.APIVersion() != "1.0" {
    return fmt.Errorf("incompatible plugin API version: %s", plugin.APIVersion())
}
```

## Examples

### Simple WASM Plugin

```go
// Host code
manager := vayu.NewManager()
plugin, err := manager.LoadWASM("./hello.wasm", vayu.Capabilities{})
if err != nil {
    log.Fatal(err)
}

result, err := plugin.Call("greet", "World")
fmt.Println(result)  // "Hello, World!"
```

### Plugin with File Access

```go
caps := vayu.Capabilities{
    FileRead: true,
}

plugin, err := manager.LoadWASM("./analyzer.wasm", caps)

// Plugin can now read files
result, err := plugin.Call("analyzeFile", "/path/to/file.txt")
```

## Performance Considerations

### WASM Overhead

WASM has ~10-20% overhead vs native:

- Use for untrusted code
- Consider native plugins for performance-critical code

### Plugin Caching

Cache plugin instances:

```go
type PluginCache struct {
    cache map[string]Plugin
    mu    sync.RWMutex
}

func (c *PluginCache) Get(name string) (Plugin, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    plugin, ok := c.cache[name]
    return plugin, ok
}
```

### Lazy Loading

Load plugins on-demand:

```go
func (m *Manager) GetOrLoad(name string) (Plugin, error) {
    if plugin := m.Get(name); plugin != nil {
        return plugin, nil
    }
    
    return m.Load(name, defaultOpts)
}
```

## Troubleshooting

### Plugin Won't Load

Check WASM compatibility:

```bash
wasm-validate plugin.wasm
```

### Permission Errors

Verify capabilities:

```go
log.Printf("Plugin capabilities: %+v", plugin.Capabilities())
```

### IPC Failures

Check process plugin logs:

```bash
./plugins/my-plugin --log-level=debug
```

## Related Modules

- **[Dravya](dravya.md)**: Plugin manager integration
- **[Vak](vak.md)**: Plugin commands

## See Also

- [Plugin Examples](../../examples/05-plugins/)
- [Plugin Development Guide](../concepts.md#plugins)
- [Security Architecture](../../SECURITY.md)
