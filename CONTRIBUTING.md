# Contributing to DRAV

Thank you for your interest in contributing to DRAV! We welcome contributions from the community.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contribution Workflow](#contribution-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Reporting Issues](#reporting-issues)

## Code of Conduct

This project adheres to the Contributor Covenant [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/drav.git`
3. Add upstream remote: `git remote add upstream https://github.com/TIVerse/drav.git`
4. Create a feature branch: `git checkout -b feature/your-feature-name`

## Development Setup

### Prerequisites

- Go 1.22 or later
- Make (optional but recommended)
- golangci-lint for linting

### Setup

```bash
# Install dependencies
go mod download

# Run tests
make test

# Run linter
make lint

# Build
make build
```

## Contribution Workflow

1. **Find or create an issue**: Look for existing issues or create a new one describing your proposed change
2. **Discuss**: For major changes, discuss your approach in the issue first
3. **Implement**: Write your code following our coding standards
4. **Test**: Add tests for your changes
5. **Document**: Update documentation as needed
6. **Submit**: Open a pull request

## Coding Standards

### Go Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` and `goimports` for formatting
- Run `golangci-lint` and fix all issues
- Write clear, idiomatic Go code

### Code Organization

- Keep functions small and focused
- Use descriptive variable and function names
- Add comments for exported functions and types
- Group related functionality in packages

### Error Handling

- Always handle errors explicitly
- Wrap errors with context using `fmt.Errorf`
- Don't panic in library code
- Return errors to callers

### Concurrency

- Document goroutine lifecycle
- Use channels for communication
- Avoid shared mutable state
- Use sync primitives when necessary

## Testing

### Test Requirements

- Write unit tests for all new code
- Maintain or improve code coverage
- Add integration tests for complex interactions
- Include edge cases and error paths

### Running Tests

```bash
# All tests
make test

# With coverage
make cover

# Specific package
go test ./pkg/dravya/...

# Benchmarks
make bench

# Fuzz tests
go test -fuzz=FuzzDiff ./tests/fuzz/
```

### Test Guidelines

- Use table-driven tests where appropriate
- Use subtests for related test cases
- Mock external dependencies
- Test both success and failure paths

## Pull Request Process

1. **Update your branch**: Rebase on latest main before submitting
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Ensure quality**:
   - All tests pass
   - Linter shows no issues
   - Code coverage maintained or improved
   - Documentation updated

3. **Write a clear PR description**:
   - What problem does it solve?
   - How does it solve it?
   - Are there breaking changes?
   - Link to related issues

4. **Review process**:
   - Maintainers will review your PR
   - Address feedback promptly
   - Be open to suggestions
   - Keep discussions respectful

5. **Merge**:
   - Maintainers will merge once approved
   - PR will be squashed into a single commit

## Reporting Issues

### Bug Reports

Include:
- DRAV version
- Go version
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Error messages and logs

### Feature Requests

Include:
- Problem description
- Proposed solution
- Use cases
- Alternatives considered

### Security Issues

**Do not open public issues for security vulnerabilities.**

See [SECURITY.md](SECURITY.md) for responsible disclosure process.

## Module-Specific Guidelines

### Dravya (Runtime)

- Ensure graceful shutdown
- Document lifecycle hooks
- Handle context cancellation

### Agni (Events)

- Non-blocking event dispatch
- Maintain priority ordering
- Document event types

### Māyā (Renderer)

- Optimize diff algorithm
- Minimize escape sequences
- Handle resize gracefully

### Prāṇa (State)

- Keep reducers pure
- Document side effects
- Ensure thread safety

### Vāk (Commands)

- Validate command definitions
- Support completion
- Implement undo where possible

### Vāyu (Plugins)

- Enforce capability boundaries
- Document security model
- Test isolation

### Śrī (Themes)

- Support color fallbacks
- Document theme structure
- Test on multiple terminals

## Documentation

- Update README.md for user-facing changes
- Add inline documentation for exported symbols
- Update module-specific docs in `docs/src/modules/`
- Include examples for new features

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Questions?

- Open a discussion on GitHub
- Join our community chat (coming soon)
- Email maintainers (see MAINTAINERS.md)

Thank you for contributing to DRAV! 🌊
