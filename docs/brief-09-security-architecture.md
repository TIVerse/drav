# Section 9: Security Architecture

[← Back to Index](brief-index.md) | [Previous: Performance Engineering](brief-08-performance-engineering.md) | [Next: Testing Strategy →](brief-10-testing-strategy.md)

---

## 9.1 Threat Model

### Attack Surface Analysis

**DRAV applications face threats from**:

1. **Malicious Plugins** - Third-party code execution
2. **Command Injection** - User input in commands
3. **Terminal Escape Injection** - Malicious terminal sequences
4. **Resource Exhaustion** - Memory/CPU DoS
5. **State Corruption** - Invalid state manipulation
6. **File System Access** - Unauthorized file operations

### Trust Boundaries

```
┌─────────────────────────────────────────┐
│           Untrusted Zone                 │
│  ┌─────────────────────────────────┐    │
│  │   User Input                     │    │
│  │   Plugins                        │    │
│  │   External Data                  │    │
│  └─────────────────────────────────┘    │
└────────────────┬────────────────────────┘
                 │ Validation
                 ▼
┌─────────────────────────────────────────┐
│           Trusted Zone                   │
│  ┌─────────────────────────────────┐    │
│  │   DRAV Core                      │    │
│  │   Validated State                │    │
│  │   Sanitized Output               │    │
│  └─────────────────────────────────┘    │
└─────────────────────────────────────────┘
```

---

## 9.2 Plugin Security

### Sandboxing Strategy

**Multi-Layer Defense**:

#### Layer 1: Capability System

```go
type PluginCapabilities struct {
    FileSystem   bool
    Network      bool
    StateWrite   bool
    CommandExec  bool
    UIModify     bool
}

type SandboxedPlugin struct {
    plugin      Plugin
    caps        PluginCapabilities
    resourceLimits ResourceLimits
}

func (sp *SandboxedPlugin) CheckCapability(cap string) error {
    switch cap {
    case "filesystem":
        if !sp.caps.FileSystem {
            return ErrCapabilityDenied
        }
    case "network":
        if !sp.caps.Network {
            return ErrCapabilityDenied
        }
    }
    return nil
}
```

#### Layer 2: Resource Limits

```go
type ResourceLimits struct {
    MaxMemory     int64         // Bytes
    MaxCPU        time.Duration // CPU time
    MaxGoroutines int
    MaxFileHandles int
    Timeout       time.Duration
}

func (rl *ResourceLimits) Enforce(plugin Plugin) error {
    // Memory limit
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    if m.Alloc > uint64(rl.MaxMemory) {
        return ErrMemoryLimitExceeded
    }
    
    // Goroutine limit
    if runtime.NumGoroutine() > rl.MaxGoroutines {
        return ErrGoroutineLimitExceeded
    }
    
    return nil
}
```

#### Layer 3: API Isolation

```go
// Plugin only gets restricted API
type PluginAPI interface {
    Log(level LogLevel, message string)
    GetConfig(key string) (string, error)
    RegisterCommand(name string, handler CommandHandler)
    // No direct access to internal state
}

// Internal API not exposed to plugins
type InternalAPI interface {
    DirectMemoryAccess()
    RawTerminalControl()
    CoreStateManipulation()
}
```

### Plugin Verification

**Signature Verification**:
```go
type PluginSignature struct {
    PluginHash  []byte
    Signature   []byte
    PublicKey   []byte
    SignedBy    string
    SignedAt    time.Time
}

func VerifyPlugin(pluginPath string, sig PluginSignature) error {
    // Read plugin bytes
    data, err := os.ReadFile(pluginPath)
    if err != nil {
        return err
    }
    
    // Verify hash
    hash := sha256.Sum256(data)
    if !bytes.Equal(hash[:], sig.PluginHash) {
        return ErrHashMismatch
    }
    
    // Verify signature
    pub, err := x509.ParsePKIXPublicKey(sig.PublicKey)
    if err != nil {
        return err
    }
    
    return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA256, hash[:], sig.Signature)
}
```

### Plugin Isolation

**Process Isolation** (future):
```go
// Run plugins in separate processes
type IsolatedPlugin struct {
    cmd     *exec.Cmd
    stdin   io.WriteCloser
    stdout  io.ReadCloser
}

func LoadIsolatedPlugin(path string) (*IsolatedPlugin, error) {
    cmd := exec.Command("drav-plugin-runner", path)
    stdin, _ := cmd.StdinPipe()
    stdout, _ := cmd.StdoutPipe()
    
    cmd.Start()
    
    return &IsolatedPlugin{cmd, stdin, stdout}, nil
}
```

---

## 9.3 Input Validation

### Command Injection Prevention

```go
// Sanitize user input
func SanitizeInput(input string) string {
    // Remove control characters
    input = strings.Map(func(r rune) rune {
        if unicode.IsControl(r) && r != '\n' && r != '\t' {
            return -1
        }
        return r
    }, input)
    
    return input
}

// Parameterized commands (like prepared statements)
func ExecuteSafeCommand(cmdName string, args ...string) error {
    cmd := registry.Get(cmdName)
    if cmd == nil {
        return ErrCommandNotFound
    }
    
    // Args are already separated, no injection possible
    return cmd.Execute(args...)
}
```

### Terminal Escape Injection

**Problem**: Malicious terminal sequences
```
User inputs: \033]0;Evil\007\033[2J
Result: Changes terminal title, clears screen
```

**Solution**: Escape sanitization
```go
// Strip dangerous escape sequences
func SanitizeTerminalOutput(s string) string {
    // Whitelist safe sequences
    safeSequences := []string{
        "\033[0m",    // Reset
        "\033[1m",    // Bold
        "\033[3m",    // Italic
        "\033[4m",    // Underline
        "\033[3[0-9]m", // Colors
    }
    
    // Remove all escape sequences except whitelisted
    re := regexp.MustCompile(`\033\[[^m]*m`)
    return re.ReplaceAllStringFunc(s, func(seq string) string {
        for _, safe := range safeSequences {
            if matched, _ := regexp.MatchString(safe, seq); matched {
                return seq
            }
        }
        return ""
    })
}
```

---

## 9.4 State Security

### State Validation

```go
type Validator[T any] interface {
    Validate(value T) error
}

type ValidatedObservable[T any] struct {
    *Observable[T]
    validator Validator[T]
}

func (vo *ValidatedObservable[T]) Set(value T) error {
    // Validate before setting
    if err := vo.validator.Validate(value); err != nil {
        return fmt.Errorf("invalid state: %w", err)
    }
    
    vo.Observable.Set(value)
    return nil
}

// Example
type PositiveIntValidator struct{}

func (v PositiveIntValidator) Validate(value int) error {
    if value < 0 {
        return errors.New("value must be non-negative")
    }
    return nil
}

counter := NewValidatedObservable(0, PositiveIntValidator{})
counter.Set(-1)  // Error: invalid state
```

### State Encryption (Sensitive Data)

```go
type EncryptedObservable[T any] struct {
    encrypted []byte
    key       []byte
}

func (eo *EncryptedObservable[T]) Get() (T, error) {
    var value T
    
    // Decrypt
    plaintext, err := decrypt(eo.encrypted, eo.key)
    if err != nil {
        return value, err
    }
    
    // Deserialize
    err = json.Unmarshal(plaintext, &value)
    return value, err
}

func (eo *EncryptedObservable[T]) Set(value T) error {
    // Serialize
    plaintext, err := json.Marshal(value)
    if err != nil {
        return err
    }
    
    // Encrypt
    eo.encrypted, err = encrypt(plaintext, eo.key)
    return err
}
```

---

## 9.5 File System Security

### Path Traversal Prevention

```go
func SafeFilePath(basePath, userPath string) (string, error) {
    // Clean and join paths
    fullPath := filepath.Join(basePath, filepath.Clean(userPath))
    
    // Check if still within base
    if !strings.HasPrefix(fullPath, basePath) {
        return "", ErrPathTraversal
    }
    
    return fullPath, nil
}

// Example
dataDir := "/app/data"
userFile := "../../etc/passwd"  // Attack attempt
safePath, err := SafeFilePath(dataDir, userFile)
// err == ErrPathTraversal
```

### File Access Control

```go
type FileAccessPolicy struct {
    AllowRead  []string  // Allowed directories
    AllowWrite []string
    Deny       []string  // Explicit denials
}

func (fap *FileAccessPolicy) CheckRead(path string) error {
    // Check deny list first
    for _, denied := range fap.Deny {
        if strings.HasPrefix(path, denied) {
            return ErrAccessDenied
        }
    }
    
    // Check allow list
    for _, allowed := range fap.AllowRead {
        if strings.HasPrefix(path, allowed) {
            return nil
        }
    }
    
    return ErrAccessDenied
}
```

---

## 9.6 Network Security (Future)

### Outbound Connection Control

```go
type NetworkPolicy struct {
    AllowedHosts []string
    AllowedPorts []int
    MaxConnections int
}

func (np *NetworkPolicy) CheckConnection(host string, port int) error {
    // Check host whitelist
    allowed := false
    for _, allowedHost := range np.AllowedHosts {
        if host == allowedHost || strings.HasSuffix(host, "."+allowedHost) {
            allowed = true
            break
        }
    }
    
    if !allowed {
        return ErrHostNotAllowed
    }
    
    // Check port whitelist
    portAllowed := false
    for _, allowedPort := range np.AllowedPorts {
        if port == allowedPort {
            portAllowed = true
            break
        }
    }
    
    if !portAllowed {
        return ErrPortNotAllowed
    }
    
    return nil
}
```

### TLS Verification

```go
func SecureHTTPClient() *http.Client {
    return &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                MinVersion: tls.VersionTLS12,
                CipherSuites: []uint16{
                    tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
                    tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
                },
            },
        },
        Timeout: 30 * time.Second,
    }
}
```

---

## 9.7 Dependency Security

### Vulnerability Scanning

```yaml
# .github/workflows/security.yml
name: Security Scan

on: [push, pull_request]

jobs:
  scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Run Gosec
        uses: securego/gosec@master
        with:
          args: ./...
      
      - name: Run Trivy
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          scan-ref: '.'
```

### Dependency Pinning

```go
// go.mod with specific versions
require (
    github.com/gdamore/tcell/v2 v2.6.0
    // Not: github.com/gdamore/tcell/v2 latest
)
```

### Supply Chain Security

```bash
# Verify checksums
go mod verify

# Use Go's checksum database
export GOSUMDB=sum.golang.org
```

---

## 9.8 Secrets Management

### Environment Variables

```go
// Secure secret loading
func LoadSecret(name string) (string, error) {
    value := os.Getenv(name)
    if value == "" {
        return "", fmt.Errorf("secret %s not found", name)
    }
    
    // Don't log secrets
    log.Printf("Loaded secret: %s", name)  // ✓ Safe
    // log.Printf("Secret value: %s", value)  // ✗ Dangerous
    
    return value, nil
}
```

### Configuration Security

```go
type SecureConfig struct {
    APIKey    string `json:"-"` // Never serialize
    Password  string `json:"-"`
    PublicURL string `json:"public_url"`
}

func (sc *SecureConfig) Save(path string) error {
    // Save to file with restricted permissions
    data, err := json.Marshal(sc)
    if err != nil {
        return err
    }
    
    // 0600 = read/write for owner only
    return os.WriteFile(path, data, 0600)
}
```

---

## 9.9 Error Handling Security

### Avoid Information Disclosure

```go
// Bad: Exposes internal paths
func (app *App) HandleError(err error) {
    fmt.Printf("Error: %v\n", err)
    // Output: Error: open /home/user/.secret: permission denied
}

// Good: Generic message to user, detailed log internally
func (app *App) HandleError(err error) {
    // User sees generic message
    fmt.Println("An error occurred. Please check logs.")
    
    // Detailed error logged securely
    app.logger.Error("Error details", 
        "error", err,
        "stack", debug.Stack(),
    )
}
```

### Stack Trace Sanitization

```go
func SanitizeStackTrace(trace string) string {
    // Remove absolute paths
    re := regexp.MustCompile(`/.*?/go/src/`)
    return re.ReplaceAllString(trace, "")
}
```

---

## 9.10 Security Best Practices

### Principle of Least Privilege

```go
// Default: No permissions
type DefaultPluginAPI struct{}

// Plugins must request specific permissions
type PrivilegedPluginAPI struct {
    DefaultPluginAPI
    fileAccess FileAccess  // Granted explicitly
}
```

### Defense in Depth

**Multiple layers**:
1. Input validation
2. Output encoding
3. Resource limits
4. Capability system
5. Process isolation (future)

### Fail Securely

```go
func AuthenticateUser(token string) (*User, error) {
    user, err := validateToken(token)
    if err != nil {
        // Fail closed: Deny access on error
        return nil, ErrUnauthorized
    }
    return user, nil
}
```

### Security Logging

```go
type SecurityLogger struct {
    logger *log.Logger
}

func (sl *SecurityLogger) LogSecurityEvent(event SecurityEvent) {
    sl.logger.Printf("[SECURITY] %s: %s (user=%s, time=%s)",
        event.Type,
        event.Description,
        event.User,
        event.Timestamp,
    )
}

// Events to log:
// - Failed authentication attempts
// - Permission denials
// - Resource limit violations
// - Plugin load/unload
// - Configuration changes
```

---

## 9.11 Compliance & Auditing

### Audit Trail

```go
type AuditLog struct {
    entries []AuditEntry
    mu      sync.Mutex
}

type AuditEntry struct {
    Timestamp time.Time
    User      string
    Action    string
    Resource  string
    Result    string
    Details   map[string]interface{}
}

func (al *AuditLog) Record(entry AuditEntry) {
    al.mu.Lock()
    defer al.mu.Unlock()
    
    entry.Timestamp = time.Now()
    al.entries = append(al.entries, entry)
    
    // Persist to secure storage
    al.persist(entry)
}
```

### GDPR Considerations

```go
// User data deletion
func (app *App) DeleteUserData(userID string) error {
    // Remove all user data
    app.state.Delete("user:" + userID)
    app.cache.Invalidate("user:" + userID)
    app.logs.Anonymize(userID)
    
    return nil
}

// Data export
func (app *App) ExportUserData(userID string) ([]byte, error) {
    data := app.state.GetAll("user:" + userID)
    return json.Marshal(data)
}
```

---

## Summary

**Security Principles**:
1. **Defense in Depth**: Multiple security layers
2. **Least Privilege**: Minimal permissions by default
3. **Input Validation**: Sanitize all external input
4. **Secure by Default**: Safe defaults, opt-in for risky features
5. **Fail Securely**: Deny access on error

**Key Security Features**:
- Plugin sandboxing with capabilities
- Resource limits enforcement
- Input/output sanitization
- Path traversal prevention
- Vulnerability scanning in CI
- Audit logging

**Security Checklist**:
- ✅ Plugin signature verification
- ✅ Resource limits per plugin
- ✅ Input validation on all boundaries
- ✅ Output sanitization for terminal
- ✅ File system access control
- ✅ Dependency vulnerability scanning
- ✅ Security audit logging

---

[← Back to Index](brief-index.md) | [Previous: Performance Engineering](brief-08-performance-engineering.md) | [Next: Testing Strategy →](brief-10-testing-strategy.md)
