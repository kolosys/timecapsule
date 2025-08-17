# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

## [0.1.0] - 2025-08-17

### Added

- Initial release of go-timecapsule library
- Core TimeCapsule interface with generic support
- MemoryTimeCapsule implementation with thread-safe operations
- PersistentTimeCapsule for pluggable storage backends
- JSONCodec for serialization/deserialization
- Comprehensive test suite with examples
- Demo application showcasing library features
- Delay functionality for extending unlock times
- Context support for timeout and cancellation
- GitHub Actions CI/CD workflows
- Automated release process with GoReleaser
- Comprehensive linting with golangci-lint
- Security scanning with gosec
- Version information in demo binary
- CHANGELOG.md for tracking changes
- Full documentation and README

### Features

- Store time-locked values with any data type
- Retrieve values only after specified unlock time
- Peek at capsule metadata without opening
- Wait for capsules to unlock with context support
- Delay unlock times dynamically
- Thread-safe concurrent access
- Extensible storage backend architecture
- Type-safe generics support
