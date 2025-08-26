# Contributing to Kolosys/Timecapsule

Thank you for your interest in contributing to Timecapsule! This document provides guidelines and information for contributors.

## 🎯 Project Overview

Timecapsule is a lightweight Go library for storing time-locked values. It provides a simple, type-safe API for creating "time capsules" that unlock values at specified times.

## 🚀 Getting Started

### Prerequisites

- **Go 1.21+** (check with `go version`)
- **Git** for version control
- **Make** (optional, for using the Makefile)

### Development Setup

1. **Fork and Clone**

   ```bash
   git clone https://github.com/YOUR_USERNAME/timecapsule.git
   cd timecapsule
   ```

2. **Install Dependencies**

   ```bash
   go mod download
   go mod tidy
   ```

3. **Run Tests**

   ```bash
   go test -v
   go test -race -v  # Test for race conditions
   ```

4. **Run with Coverage**

   ```bash
   go test -v -cover
   go test -v -coverprofile=coverage.out
   go tool cover -html=coverage.out  # View coverage in browser
   ```

5. **Format and Lint**
   ```bash
   go fmt ./...
   go vet ./...
   ```

## 📝 Contributing Guidelines

### Issue Reporting

Before creating an issue, please:

1. **Search existing issues** to avoid duplicates
2. **Use a clear, descriptive title**
3. **Provide reproduction steps** for bugs
4. **Include Go version and OS** information
5. **Add relevant code samples** when applicable

**Bug Report Template:**

````markdown
## Bug Description

Brief description of the issue

## Steps to Reproduce

1. Step 1
2. Step 2
3. Step 3

## Expected Behavior

What should happen

## Actual Behavior

What actually happens

## Environment

- Go version: `go version`
- OS:
- timecapsule version:

## Code Sample

```go
// Minimal code to reproduce the issue
```
````

### Feature Requests

When requesting features:

1. **Explain the use case** and problem being solved
2. **Provide examples** of the proposed API
3. **Consider backwards compatibility**
4. **Check the roadmap** to see if it's already planned

### Pull Request Process

1. **Create a feature branch**

   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/issue-description
   ```

2. **Make your changes**

   - Follow the [coding standards](#coding-standards)
   - Add tests for new functionality
   - Update documentation if needed
   - Ensure all tests pass

3. **Commit your changes**

   ```bash
   git add .
   git commit -m "feat: add new time capsule feature"
   # or
   git commit -m "fix: resolve unlock time calculation bug"
   ```

4. **Push and create PR**

   ```bash
   git push origin feature/your-feature-name
   ```

5. **Fill out the PR template** with:
   - Description of changes
   - Related issue links
   - Testing details
   - Breaking changes (if any)

## 📋 Coding Standards

### Code Style

- **Follow standard Go formatting**: Use `go fmt` and `go vet`
- **Use meaningful variable names**: `unlockTime` not `t`
- **Add comments for exported functions**: Follow Go documentation conventions
- **Keep functions focused**: Single responsibility principle
- **Prefer composition over inheritance**

### Function Documentation

```go
// Store saves a time-locked value that becomes accessible at unlockTime.
// The value is serialized using the configured codec and stored in the backend.
// Returns an error if the key already exists or storage fails.
func (tc *MemoryTimeCapsule[T]) Store(ctx context.Context, key string, value T, unlockTime time.Time) error {
    // implementation
}
```

### Error Handling

- **Use sentinel errors** for common error types:

  ```go
  var (
      ErrCapsuleLocked   = errors.New("capsule is still locked")
      ErrCapsuleNotFound = errors.New("capsule not found")
      ErrInvalidKey      = errors.New("invalid key")
  )
  ```

- **Wrap errors with context**:
  ```go
  if err != nil {
      return fmt.Errorf("failed to store capsule %q: %w", key, err)
  }
  ```

### Testing Standards

1. **Test Coverage**: Aim for >90% coverage
2. **Table-driven tests** for multiple scenarios:

   ```go
   func TestStore(t *testing.T) {
       tests := []struct {
           name        string
           key         string
           value       string
           unlockTime  time.Time
           expectError bool
       }{
           {
               name:       "valid storage",
               key:        "test-key",
               value:      "test-value",
               unlockTime: time.Now().Add(time.Hour),
           },
           // more test cases...
       }

       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               // test implementation
           })
       }
   }
   ```

3. **Test with context cancellation**:

   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
   defer cancel()
   ```

4. **Test race conditions**:

   ```bash
   go test -race ./...
   ```

5. **Example tests for documentation**:
   ```go
   func ExampleTimeCapsule_Store() {
       capsule := timecapsule.New[string]()
       err := capsule.Store(context.Background(), "greeting", "Hello!", time.Now().Add(time.Hour))
       fmt.Println(err == nil)
       // Output: true
   }
   ```

## 🏗️ Architecture Guidelines

### Package Structure

The library follows a simple, focused structure:

```
timecapsule/
├── timecapsule.go      # Core interfaces and types
├── storage.go          # Storage interface and implementations
├── codec.go            # Serialization interface and implementations
├── example_test.go     # Example tests for documentation
├── timecapsule_test.go # Unit tests
└── cmd/demo/           # Demo application
```

### Interface Design

- **Keep interfaces small** and focused
- **Use generics** for type safety: `TimeCapsule[T any]`
- **Accept interfaces, return structs**
- **Make zero values useful**

### Storage Backends

When adding new storage backends:

1. **Implement the `Storage` interface**
2. **Handle context cancellation** properly
3. **Add comprehensive tests**
4. **Document configuration options**
5. **Consider connection pooling and cleanup**

Example:

```go
type RedisStorage struct {
    client *redis.Client
    prefix string
}

func (r *RedisStorage) Store(ctx context.Context, key string, value []byte, unlockTime time.Time) error {
    // Implementation with proper context handling
}
```

## 🧪 Testing

### Running Tests

```bash
# Basic tests
go test -v

# With race detection
go test -race -v

# With coverage
go test -v -cover

# Specific test
go test -run TestStore -v

# Benchmarks
go test -bench=. -v
```

### Test Categories

1. **Unit Tests**: Test individual functions and methods
2. **Integration Tests**: Test storage backends with real dependencies
3. **Race Tests**: Ensure thread safety
4. **Example Tests**: Validate documentation examples
5. **Benchmark Tests**: Performance validation

### Mocking and Test Utilities

- Use **dependency injection** for testability
- Create **test helpers** for common setup:
  ```go
  func setupTestCapsule(t *testing.T) TimeCapsule[string] {
      t.Helper()
      return New[string]()
  }
  ```

## 📚 Documentation

### Code Documentation

- **Document all exported types and functions**
- **Use examples in documentation**:
  ```go
  // Store saves a value with a future unlock time.
  //
  // Example:
  //   capsule := timecapsule.New[string]()
  //   err := capsule.Store(ctx, "key", "value", time.Now().Add(time.Hour))
  func Store(ctx context.Context, key string, value T, unlockTime time.Time) error
  ```

### README Updates

When adding features:

- Update the feature list
- Add usage examples
- Update the API reference
- Add to use cases if applicable

### Changelog

Follow [Keep a Changelog](https://keepachangelog.com/) format:

```markdown
## [1.2.0] - 2024-01-15

### Added

- Redis storage backend
- Bulk operations support

### Changed

- Improved error messages

### Fixed

- Race condition in memory storage
```

## 🔄 Development Workflow

### Branch Naming

- **Features**: `feature/add-redis-backend`
- **Bug fixes**: `fix/unlock-time-calculation`
- **Documentation**: `docs/improve-readme`
- **Refactoring**: `refactor/storage-interface`

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat: add Redis storage backend`
- `fix: resolve race condition in memory storage`
- `docs: update API documentation`
- `test: add integration tests for storage`
- `refactor: simplify codec interface`

### Release Process

1. **Update version** in relevant files
2. **Update CHANGELOG.md**
3. **Create release notes**
4. **Tag the release**: `git tag v1.2.0`
5. **Push tags**: `git push origin --tags`

## 🎯 Priority Areas

We're especially interested in contributions for:

1. **Storage Backends**

   - Redis implementation
   - PostgreSQL/SQLite support
   - File system backend
   - Cloud storage options

2. **Features**

   - TTL cleanup and auto-purging
   - Bulk operations
   - Encryption support
   - Webhook notifications

3. **Performance**

   - Benchmarking and optimization
   - Memory usage improvements
   - Concurrent access optimization

4. **Developer Experience**
   - Better error messages
   - More examples and tutorials
   - CLI tools and utilities

## ❓ Getting Help

- **Documentation**: Check the [README](README.md) and code comments
- **Issues**: Search [existing issues](https://github.com/kolosys/timecapsule/issues)
- **Discussions**: Use [GitHub Discussions](https://github.com/kolosys/timecapsule/discussions)
- **Chat**: Join our community discussions

## 📄 License

By contributing to timecapsule, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to timecapsule! 🚀
