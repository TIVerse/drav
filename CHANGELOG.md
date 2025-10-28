# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project structure and build system
- Core modules: Dravya, Agni, Māyā, Prāṇa, Vāk, Vāyu, Śrī
- Event system with priority-based dispatch
- Reactive state management (Observable, Computed, Store)
- Command engine with parsing, history, and undo/redo
- Plugin system with capability-based security
- Theme system with dark and light themes
- Animation system with easing functions
- Layout engine (Row, Column, Flex, Grid, Stack)
- tcell-based renderer with diff algorithm
- 7 example applications demonstrating core features
- Comprehensive test infrastructure (unit, e2e, fuzz)
- CI/CD workflows (lint, test, security, release)
- Documentation structure and guides
- Security policy and vulnerability reporting
- Contributing guidelines and code of conduct

### Changed
- N/A (initial release)

### Deprecated
- N/A (initial release)

### Removed
- N/A (initial release)

### Fixed
- N/A (initial release)

### Security
- Input sanitization to prevent ANSI escape injection
- Capability-based plugin isolation
- Path traversal prevention
- Resource limits for plugins

## [0.1.0] - TBD

### Added
- First public alpha release
- Core framework functionality
- Basic examples and documentation

---

## Release Notes

### Version Numbering

- **Major.Minor.Patch** (e.g., 1.2.3)
- **Major**: Breaking changes
- **Minor**: New features, backwards compatible
- **Patch**: Bug fixes, backwards compatible

### Pre-1.0 Releases

During pre-1.0 development (0.x.x), minor version increments may include breaking changes. We will clearly document these in the changelog and provide migration guides.

### Support Policy

- **Latest minor version**: Full support with security patches and bug fixes
- **Previous minor version**: Security patches only for 6 months
- **Older versions**: No official support

---

[Unreleased]: https://github.com/TIVerse/drav/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/TIVerse/drav/releases/tag/v0.1.0
