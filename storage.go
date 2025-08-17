package timecapsule

import (
	"context"
	"time"
)

// Storage defines the interface for persistent storage backends
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

// PersistentTimeCapsule implements TimeCapsule using a persistent storage backend
type PersistentTimeCapsule[T any] struct {
	storage Storage
	codec   Codec[T]
}

// Codec defines how to serialize/deserialize values
type Codec[T any] interface {
	Encode(value T) ([]byte, error)
	Decode(data []byte) (T, error)
}

// NewWithStorage creates a new time capsule with persistent storage
func NewWithStorage[T any](storage Storage, codec Codec[T]) TimeCapsule[T] {
	return &PersistentTimeCapsule[T]{
		storage: storage,
		codec:   codec,
	}
}

// Store stores a value in a time capsule that will be unlocked at the specified time
func (tc *PersistentTimeCapsule[T]) Store(ctx context.Context, key string, value T, unlockTime time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if key == "" {
		return ErrInvalidKey
	}

	data, err := tc.codec.Encode(value)
	if err != nil {
		return err
	}

	return tc.storage.Store(ctx, key, data, unlockTime)
}

// Open retrieves a value from a time capsule if it's unlocked
func (tc *PersistentTimeCapsule[T]) Open(ctx context.Context, key string) (T, error) {
	if err := ctx.Err(); err != nil {
		var zero T
		return zero, err
	}

	if key == "" {
		var zero T
		return zero, ErrInvalidKey
	}

	data, err := tc.storage.Open(ctx, key)
	if err != nil {
		var zero T
		return zero, err
	}

	return tc.codec.Decode(data)
}

// Peek returns metadata about a capsule without opening it
func (tc *PersistentTimeCapsule[T]) Peek(ctx context.Context, key string) (Metadata, error) {
	if err := ctx.Err(); err != nil {
		return Metadata{}, err
	}

	if key == "" {
		return Metadata{}, ErrInvalidKey
	}

	return tc.storage.Peek(ctx, key)
}

// Delay delays the unlock time of a capsule
func (tc *PersistentTimeCapsule[T]) Delay(ctx context.Context, key string, delay time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if key == "" {
		return ErrInvalidKey
	}

	// Check if capsule exists
	_, err := tc.storage.Peek(ctx, key)
	if err != nil {
		return err
	}

	// Calculate new unlock time
	newUnlockTime := time.Now().Add(delay)

	// Get the current value
	data, err := tc.storage.Open(ctx, key)
	if err != nil {
		return err
	}

	// Re-store with new unlock time
	return tc.storage.Store(ctx, key, data, newUnlockTime)
}

// Delete removes a capsule from storage
func (tc *PersistentTimeCapsule[T]) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if key == "" {
		return ErrInvalidKey
	}

	return tc.storage.Delete(ctx, key)
}

// Exists checks if a capsule exists
func (tc *PersistentTimeCapsule[T]) Exists(ctx context.Context, key string) bool {
	if err := ctx.Err(); err != nil {
		return false
	}

	if key == "" {
		return false
	}

	return tc.storage.Exists(ctx, key)
}

// WaitForUnlock blocks until a capsule is unlocked or context is canceled
func (tc *PersistentTimeCapsule[T]) WaitForUnlock(ctx context.Context, key string) (T, error) {
	if err := ctx.Err(); err != nil {
		var zero T
		return zero, err
	}

	// First check if capsule exists
	metadata, err := tc.Peek(ctx, key)
	if err != nil {
		var zero T
		return zero, err
	}

	// If already unlocked, return immediately
	if !metadata.IsLocked {
		return tc.Open(ctx, key)
	}

	// Wait until unlock time or context cancellation
	timer := time.NewTimer(time.Until(metadata.UnlockTime))
	defer timer.Stop()

	select {
	case <-ctx.Done():
		var zero T
		return zero, ctx.Err()
	case <-timer.C:
		return tc.Open(ctx, key)
	}
}
