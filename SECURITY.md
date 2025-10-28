# Security Policy

## Supported Versions

We release patches for security vulnerabilities in the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 0.5.x   | :white_check_mark: |
| < 0.5   | :x:                |

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to **security@tiverse.dev**.

You should receive a response within 48 hours. If for some reason you do not, please follow up via email to ensure we received your original message.

Please include the following information:

- Type of issue (e.g., buffer overflow, SQL injection, cross-site scripting, etc.)
- Full paths of source file(s) related to the manifestation of the issue
- The location of the affected source code (tag/branch/commit or direct URL)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit it

This information will help us triage your report more quickly.

## Security Model

### Threat Model

DRAV considers the following threats:

1. **Malicious terminal input**: ANSI escape sequence injection
2. **Plugin attacks**: Unauthorized filesystem/network access
3. **Resource exhaustion**: Memory/CPU DoS via plugins
4. **Path traversal**: Unauthorized file access
5. **Command injection**: Via command parser

### Security Boundaries

#### Terminal Input
- All terminal input is sanitized to remove dangerous escape sequences
- See `internal/ansi/sanitize.go` for implementation
- XSS-like attacks via terminal escape codes are prevented

#### Plugin System
- Plugins run with capability-based security (see `pkg/vayu/capability.go`)
- Filesystem access is restricted to explicitly allowed paths
- Network access is limited to whitelisted domains/ports
- Resource limits enforced (memory, goroutines, timeouts)
- Plugin isolation via WASM or out-of-process execution

#### Filesystem Operations
- Path traversal prevention in `internal/osutil/pathsafe.go`
- All paths validated against allowed roots
- Symbolic link attacks mitigated

#### Command Execution
- Command input sanitized and validated
- No shell execution - commands run directly
- Arguments properly escaped

### Security Features

#### Input Sanitization
```go
// All terminal input is sanitized
import "github.com/TIVerse/drav/internal/ansi"
safe := ansi.Sanitize(userInput)
```

#### Capability-Based Plugin Security
```go
caps := vayu.Capabilities{
    Filesystem: vayu.FSCapability{
        Read:  []string{"/allowed/path"},
        Write: []string{"/allowed/output"},
    },
    Network: vayu.NetworkCapability{
        AllowedDomains: []string{"api.safe.com"},
        RateLimit:      100, // requests per minute
    },
}
```

#### Safe Path Operations
```go
import "github.com/TIVerse/drav/internal/osutil"
safePath, err := osutil.SafePath(baseDir, userPath)
```

### Best Practices for Users

1. **Plugin Sources**: Only load plugins from trusted sources
2. **Capability Review**: Audit plugin capabilities before loading
3. **Updates**: Keep DRAV updated to receive security patches
4. **Input Validation**: Validate all user input in your applications
5. **Resource Limits**: Configure appropriate resource limits for plugins

### Security Audits

- Code is regularly scanned with gosec
- Dependencies checked with govulncheck
- Container images scanned with Trivy
- All PRs undergo security review

### Dependency Management

We use:
- `go.mod` for dependency tracking
- Dependabot for automated updates
- Regular security audits of dependencies
- Minimal dependency footprint

### Known Limitations

1. **WASM Isolation**: WASM plugin support is experimental
2. **Terminal Escape Sequences**: Complex escape sequences may not be fully sanitized
3. **Resource Limits**: Best-effort enforcement on some platforms
4. **Go Plugin System**: Linux-only and has known security limitations

### Security Checklist for Contributors

- [ ] Input validation for all user data
- [ ] No use of `unsafe` package without justification
- [ ] Error messages don't leak sensitive information
- [ ] No hardcoded secrets or credentials
- [ ] Proper cleanup of sensitive data
- [ ] Thread-safe code with documented concurrency
- [ ] Tests include security edge cases

## Disclosure Policy

When we receive a security bug report, we will:

1. Confirm the problem and determine affected versions
2. Audit code to find similar problems
3. Prepare fixes for all supported versions
4. Release patches as soon as possible

## Public Disclosure Timing

We prefer to fully disclose the vulnerability as soon as a fix is available. We will coordinate with you on the disclosure timing and give you credit for the discovery unless you prefer to remain anonymous.

## Comments on This Policy

If you have suggestions on how this process could be improved, please submit a pull request or open an issue.

## Security Hall of Fame

We maintain a list of security researchers who have responsibly disclosed vulnerabilities:

_No entries yet - be the first!_

---

**Last Updated**: 2025-10-29

Thank you for helping keep DRAV and its users safe!
