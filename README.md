# Timecapsule - Time-Based Data Storage for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/kolosys/timecapsule.svg)](https://pkg.go.dev/github.com/kolosys/timecapsule)
[![Go Report Card](https://goreportcard.com/badge/github.com/kolosys/timecapsule)](https://goreportcard.com/report/github.com/kolosys/timecapsule)

Timecapsule is a lightweight Go library that provides time-based data storage and retrieval. Store values that are only accessible after a specified time - like a "sealed envelope" or "time capsule" for objects, configurations, or application state.

## 🎯 Problem Statement

Small–mid companies often need to:

- Schedule **delayed actions** (e.g. send an email after 7 days)
- Store **future‑effective configs** (e.g. new pricing goes live next month)
- Implement **time‑locked features** (e.g. promo codes, trials)

Current approaches:

- Cron jobs → external, brittle
- Timers/goroutines → memory leaks, fragile across restarts
- DB "valid_from/valid_until" hacks → clunky boilerplate

There's **no simple Go‑native abstraction** for "don't unlock this value until X time."

## Features

- **Simple API** - Store and retrieve time-locked values with ease
- **Type Safety** - Full generics support for any data type
- **Context Support** - Proper timeout and cancellation handling
- **Thread Safe** - Concurrent access with read-write mutexes
- **Extensible** - Pluggable storage backends (in-memory included)
- **Minimal Dependencies** - Core functionality has zero external dependencies

## Quick Start

### Installation

```bash
go get github.com/kolosys/timecapsule@latest
```

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/kolosys/timecapsule"
)

func main() {
    // Create a new time capsule
    capsule := timecapsule.New[string]()

    // Store a value that unlocks in 1 second
    unlockTime := time.Now().Add(1 * time.Second)
    err := capsule.Store(context.Background(), "greeting", "Hello, World!", unlockTime)
    if err != nil {
        log.Fatal(err)
    }

    // Try to open immediately - should be locked
    if _, err := capsule.Open(context.Background(), "greeting"); err != nil {
        fmt.Println("Capsule is locked:", err)
    }

    // Wait for unlock
    time.Sleep(2 * time.Second)

    // Now open the capsule
    value, err := capsule.Open(context.Background(), "greeting")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Unlocked value:", value)
}
```

## Design Principles

- **Context-First** - All operations accept context for cancellation/timeouts
- **No Panics** - Library code returns errors instead of panicking  
- **Minimal Dependencies** - Core functionality has zero external dependencies
- **Thread-Safe** - All public APIs are safe for concurrent use
- **Type Safety** - Full generics support with compile-time type checking

## API Reference

### Core Types

```go
// TimeCapsule is the main interface
type TimeCapsule[T any] interface {
    Store(ctx context.Context, key string, value T, unlockTime time.Time) error
    Open(ctx context.Context, key string) (T, error)
    Peek(ctx context.Context, key string) (Metadata, error)
    Delay(ctx context.Context, key string, delay time.Duration) error
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) bool
    WaitForUnlock(ctx context.Context, key string) (T, error)
}

// Capsule represents a time-locked value
type Capsule[T any] struct {
    Value      T        `json:"value"`
    UnlockTime time.Time `json:"unlock_time"`
    CreatedAt  time.Time `json:"created_at"`
}

// Metadata contains information about a capsule
type Metadata struct {
    UnlockTime time.Time `json:"unlock_time"`
    CreatedAt  time.Time `json:"created_at"`
    IsLocked   bool      `json:"is_locked"`
}
```

### Methods

#### `New[T any]() TimeCapsule[T]`

Creates a new in-memory time capsule.

#### `Store(ctx, key, value, unlockTime) error`

Stores a value that will be unlocked at the specified time.

#### `Open(ctx, key) (T, error)`

Retrieves a value if it's unlocked. Returns an error if the capsule is still locked.

#### `Peek(ctx, key) (Metadata, error)`

Returns metadata about a capsule without opening it.

#### `Delay(ctx, key, delay) error`

Delays the unlock time of a capsule by the specified duration.

#### `Delete(ctx, key) error`

Removes a capsule from storage.

#### `Exists(ctx, key) bool`

Checks if a capsule exists.

#### `WaitForUnlock(ctx, key) (T, error)`

Blocks until a capsule is unlocked or the context is cancelled.

## Usage Examples

### Basic Usage

```go
capsule := timecapsule.New[string]()

// Store a value
err := capsule.Store(context.Background(), "secret", "confidential",
    time.Now().Add(24*time.Hour))

// Check if it's locked
metadata, _ := capsule.Peek(context.Background(), "secret")
fmt.Printf("Is locked: %v\n", metadata.IsLocked)

// Try to open (will fail if locked)
value, err := capsule.Open(context.Background(), "secret")
if err != nil {
    fmt.Println("Still locked:", err)
}
```

### Struct Values

```go
type Promo struct {
    Code     string `json:"code"`
    Discount int    `json:"discount"`
}

capsule := timecapsule.New[Promo]()

promo := Promo{Code: "HOLIDAY50", Discount: 50}
err := capsule.Store(context.Background(), "holiday-sale", promo,
    time.Now().Add(24*time.Hour))

// Later...
retrievedPromo, err := capsule.Open(context.Background(), "holiday-sale")
```

### Waiting for Unlock

```go
capsule := timecapsule.New[int]()

// Store a value that unlocks in 5 seconds
err := capsule.Store(context.Background(), "count", 42,
    time.Now().Add(5*time.Second))

// Wait for unlock with timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

value, err := capsule.WaitForUnlock(ctx, "count")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Got value: %d\n", value)
```

### Context Cancellation

```go
capsule := timecapsule.New[string]()

// Store a value
err := capsule.Store(context.Background(), "slow", "takes time",
    time.Now().Add(time.Hour))

// Create a context that cancels after 1 second
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()

// This will fail due to context cancellation
_, err = capsule.WaitForUnlock(ctx, "slow")
if err != nil {
    fmt.Println("Context cancelled:", err)
}
```

### Delaying Capsules

```go
capsule := timecapsule.New[string]()

// Store a value that unlocks in 1 hour
err := capsule.Store(context.Background(), "secret", "confidential",
    time.Now().Add(time.Hour))

// Delay by 2 more hours
err = capsule.Delay(context.Background(), "secret", 2*time.Hour)

// Check the new unlock time
metadata, _ := capsule.Peek(context.Background(), "secret")
fmt.Printf("New unlock time: %v\n", metadata.UnlockTime)
```

## Architecture

### In-Memory Storage

The default implementation uses an in-memory map with read-write mutexes for thread safety:

```go
type MemoryTimeCapsule[T any] struct {
    capsules map[string]Capsule[T]
    mu       sync.RWMutex
}
```

### Extensible Storage

The library supports pluggable storage backends through the `Storage` interface:

```go
type Storage interface {
    Store(ctx context.Context, key string, value []byte, unlockTime time.Time) error
    Open(ctx context.Context, key string) ([]byte, error)
    Peek(ctx context.Context, key string) (Metadata, error)
    Delete(ctx context.Context, key string) error
    Exists(ctx context.Context, key string) bool
    Close() error
}
```

Future backends will include:

- Redis
- PostgreSQL
- SQLite
- File system

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

Licensed under the [MIT License](LICENSE).

## Use Cases

### Delayed Actions

```go
// Schedule a welcome email for tomorrow
capsule.Store(ctx, "welcome-email-user123", emailData,
    time.Now().Add(24*time.Hour))
```

### Feature Flags

```go
// Enable new feature next week
capsule.Store(ctx, "new-ui-feature", true,
    time.Now().Add(7*24*time.Hour))
```

### Promotional Codes

```go
// Holiday sale starts on Black Friday
capsule.Store(ctx, "black-friday-sale", promoCode,
    blackFridayDate)
```

### Configuration Changes

```go
// New pricing goes live next month
capsule.Store(ctx, "new-pricing", pricingConfig,
    time.Now().Add(30*24*time.Hour))
```


