# timecapsule API Reference
## Functions
### Example_basic

Example_basic demonstrates basic usage of the time capsule


```go
func Example_basic()
```
### Example_context

Example_context demonstrates using context for cancellation


```go
func Example_context()
```
### Example_management

Example_management demonstrates capsule management operations


```go
func Example_management()
```
### Example_struct

Example_struct demonstrates using structs with the time capsule


```go
func Example_struct()
```
### Example_waitForUnlock

Example_waitForUnlock demonstrates waiting for a capsule to unlock


```go
func Example_waitForUnlock()
```
## Types
### Capsule

Capsule represents a time-locked value


```go
type Capsule[T any] struct {
	Value      T         `json:"value"`
	UnlockTime time.Time `json:"unlock_time"`
	CreatedAt  time.Time `json:"created_at"`
}
```
#### Fields

| Field | Type | Description |
|-------|------|-------------|
| `Value` | `T` |  |
| `UnlockTime` | `time.Time` |  |
| `CreatedAt` | `time.Time` |  |
### Codec

Codec defines how to serialize/deserialize values


```go
type Codec[T any] interface {
	Encode(value T) ([]byte, error)
	Decode(data []byte) (T, error)
}
```
### JSONCodec

JSONCodec implements Codec using JSON encoding


```go
type JSONCodec[T any] struct{}
```
#### Methods
##### Decode

Decode deserializes JSON bytes to a value


```go
func (c *JSONCodec[T]) Decode(data []byte) (T, error)
```
##### Encode

Encode serializes a value to JSON bytes


```go
func (c *JSONCodec[T]) Encode(value T) ([]byte, error)
```
### MemoryTimeCapsule

MemoryTimeCapsule implements TimeCapsule using in-memory storage


```go
type MemoryTimeCapsule[T any] struct {
	capsules map[string]Capsule[T]
	mu       sync.RWMutex
}
```
#### Fields

| Field | Type | Description |
|-------|------|-------------|
| `capsules` | `map[string]Capsule[T]` |  |
| `mu` | `sync.RWMutex` |  |
#### Methods
##### Delay

Delay delays the unlock time of a capsule


```go
func (tc *MemoryTimeCapsule[T]) Delay(ctx context.Context, key string, delay time.Duration) error
```
##### Delete

Delete removes a capsule from storage


```go
func (tc *MemoryTimeCapsule[T]) Delete(ctx context.Context, key string) error
```
##### Exists

Exists checks if a capsule exists


```go
func (tc *MemoryTimeCapsule[T]) Exists(ctx context.Context, key string) bool
```
##### Open

Open retrieves a value from a time capsule if it's unlocked


```go
func (tc *MemoryTimeCapsule[T]) Open(ctx context.Context, key string) (T, error)
```
##### Peek

Peek returns metadata about a capsule without opening it


```go
func (tc *MemoryTimeCapsule[T]) Peek(ctx context.Context, key string) (Metadata, error)
```
##### Store

Store stores a value in a time capsule that will be unlocked at the specified time


```go
func (tc *MemoryTimeCapsule[T]) Store(ctx context.Context, key string, value T, unlockTime time.Time) error
```
##### WaitForUnlock

WaitForUnlock blocks until a capsule is unlocked or context is canceled


```go
func (tc *MemoryTimeCapsule[T]) WaitForUnlock(ctx context.Context, key string) (T, error)
```
### Metadata

Metadata contains information about a capsule without exposing its value


```go
type Metadata struct {
	UnlockTime time.Time `json:"unlock_time"`
	CreatedAt  time.Time `json:"created_at"`
	IsLocked   bool      `json:"is_locked"`
}
```
#### Fields

| Field | Type | Description |
|-------|------|-------------|
| `UnlockTime` | `time.Time` |  |
| `CreatedAt` | `time.Time` |  |
| `IsLocked` | `bool` |  |
### PersistentTimeCapsule

PersistentTimeCapsule implements TimeCapsule using a persistent storage backend


```go
type PersistentTimeCapsule[T any] struct {
	storage Storage
	codec   Codec[T]
}
```
#### Fields

| Field | Type | Description |
|-------|------|-------------|
| `storage` | `Storage` |  |
| `codec` | `Codec[T]` |  |
#### Methods
##### Delay

Delay delays the unlock time of a capsule


```go
func (tc *PersistentTimeCapsule[T]) Delay(ctx context.Context, key string, delay time.Duration) error
```
##### Delete

Delete removes a capsule from storage


```go
func (tc *PersistentTimeCapsule[T]) Delete(ctx context.Context, key string) error
```
##### Exists

Exists checks if a capsule exists


```go
func (tc *PersistentTimeCapsule[T]) Exists(ctx context.Context, key string) bool
```
##### Open

Open retrieves a value from a time capsule if it's unlocked


```go
func (tc *PersistentTimeCapsule[T]) Open(ctx context.Context, key string) (T, error)
```
##### Peek

Peek returns metadata about a capsule without opening it


```go
func (tc *PersistentTimeCapsule[T]) Peek(ctx context.Context, key string) (Metadata, error)
```
##### Store

Store stores a value in a time capsule that will be unlocked at the specified time


```go
func (tc *PersistentTimeCapsule[T]) Store(ctx context.Context, key string, value T, unlockTime time.Time) error
```
##### WaitForUnlock

WaitForUnlock blocks until a capsule is unlocked or context is canceled


```go
func (tc *PersistentTimeCapsule[T]) WaitForUnlock(ctx context.Context, key string) (T, error)
```
### Storage

Storage defines the interface for persistent storage backends


```go
type Storage interface {
	// Store stores a value with its unlock time
	Store(ctx context.Context, key string, value []byte, unlockTime time.Time) error

	// Open retrieves a value if it's unlocked
	Open(ctx context.Context, key string) ([]byte, error)

	// Peek returns metadata about a capsule without opening it
	Peek(ctx context.Context, key string) (Metadata, error)

	// Delete removes a capsule
	Delete(ctx context.Context, key string) error

	// Exists checks if a capsule exists
	Exists(ctx context.Context, key string) bool

	// Close closes the storage connection
	Close() error
}
```
### TimeCapsule

TimeCapsule is the main interface for storing and retrieving time-locked values


```go
type TimeCapsule[T any] interface {
	Store(ctx context.Context, key string, value T, unlockTime time.Time) error
	Open(ctx context.Context, key string) (T, error)
	Peek(ctx context.Context, key string) (Metadata, error)
	Delay(ctx context.Context, key string, delay time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) bool
	WaitForUnlock(ctx context.Context, key string) (T, error)
}
```
## Variables
### ErrCapsuleNotFound

Common errors


```go
ErrCapsuleNotFound = errors.New("capsule not found")
```
### ErrCapsuleLocked

Common errors


```go
ErrCapsuleLocked = errors.New("capsule is still locked")
```
### ErrInvalidKey

Common errors


```go
ErrInvalidKey = errors.New("invalid key")
```
