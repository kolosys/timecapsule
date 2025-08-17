package timecapsule

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Common errors
var (
	ErrCapsuleNotFound = errors.New("capsule not found")
	ErrCapsuleLocked   = errors.New("capsule is still locked")
	ErrInvalidKey      = errors.New("invalid key")
)

// Capsule represents a time-locked value
type Capsule[T any] struct {
	Value      T         `json:"value"`
	UnlockTime time.Time `json:"unlock_time"`
	CreatedAt  time.Time `json:"created_at"`
}

// Metadata contains information about a capsule without exposing its value
type Metadata struct {
	UnlockTime time.Time `json:"unlock_time"`
	CreatedAt  time.Time `json:"created_at"`
	IsLocked   bool      `json:"is_locked"`
}

// TimeCapsule is the main interface for storing and retrieving time-locked values
type TimeCapsule[T any] interface {
	Store(ctx context.Context, key string, value T, unlockTime time.Time) error
	Open(ctx context.Context, key string) (T, error)
	Peek(ctx context.Context, key string) (Metadata, error)
	Delay(ctx context.Context, key string, delay time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) bool
	WaitForUnlock(ctx context.Context, key string) (T, error)
}

// MemoryTimeCapsule implements TimeCapsule using in-memory storage
type MemoryTimeCapsule[T any] struct {
	capsules map[string]Capsule[T]
	mu       sync.RWMutex
}

// New creates a new in-memory time capsule
func New[T any]() TimeCapsule[T] {
	return &MemoryTimeCapsule[T]{
		capsules: make(map[string]Capsule[T]),
	}
}

// Store stores a value in a time capsule that will be unlocked at the specified time
func (tc *MemoryTimeCapsule[T]) Store(ctx context.Context, key string, value T, unlockTime time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if key == "" {
		return ErrInvalidKey
	}

	tc.mu.Lock()
	defer tc.mu.Unlock()

	capsule := Capsule[T]{
		Value:      value,
		UnlockTime: unlockTime,
		CreatedAt:  time.Now(),
	}

	tc.capsules[key] = capsule
	return nil
}

// Open retrieves a value from a time capsule if it's unlocked
func (tc *MemoryTimeCapsule[T]) Open(ctx context.Context, key string) (T, error) {
	if err := ctx.Err(); err != nil {
		var zero T
		return zero, err
	}

	if key == "" {
		var zero T
		return zero, ErrInvalidKey
	}

	tc.mu.RLock()
	capsule, exists := tc.capsules[key]
	tc.mu.RUnlock()

	if !exists {
		var zero T
		return zero, ErrCapsuleNotFound
	}

	if time.Now().Before(capsule.UnlockTime) {
		var zero T
		return zero, ErrCapsuleLocked
	}

	return capsule.Value, nil
}

// Peek returns metadata about a capsule without opening it
func (tc *MemoryTimeCapsule[T]) Peek(ctx context.Context, key string) (Metadata, error) {
	if err := ctx.Err(); err != nil {
		return Metadata{}, err
	}

	if key == "" {
		return Metadata{}, ErrInvalidKey
	}

	tc.mu.RLock()
	capsule, exists := tc.capsules[key]
	tc.mu.RUnlock()

	if !exists {
		return Metadata{}, ErrCapsuleNotFound
	}

	now := time.Now()
	return Metadata{
		UnlockTime: capsule.UnlockTime,
		CreatedAt:  capsule.CreatedAt,
		IsLocked:   now.Before(capsule.UnlockTime),
	}, nil
}

// Delay delays the unlock time of a capsule
func (tc *MemoryTimeCapsule[T]) Delay(ctx context.Context, key string, delay time.Duration) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if key == "" {
		return ErrInvalidKey
	}

	tc.mu.Lock()
	defer tc.mu.Unlock()

	capsule, exists := tc.capsules[key]
	if !exists {
		return ErrCapsuleNotFound
	}

	capsule.UnlockTime = time.Now().Add(delay)
	tc.capsules[key] = capsule
	return nil
}

// Delete removes a capsule from storage
func (tc *MemoryTimeCapsule[T]) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if key == "" {
		return ErrInvalidKey
	}

	tc.mu.Lock()
	defer tc.mu.Unlock()

	if _, exists := tc.capsules[key]; !exists {
		return ErrCapsuleNotFound
	}

	delete(tc.capsules, key)
	return nil
}

// Exists checks if a capsule exists
func (tc *MemoryTimeCapsule[T]) Exists(ctx context.Context, key string) bool {
	if err := ctx.Err(); err != nil {
		return false
	}

	if key == "" {
		return false
	}

	tc.mu.RLock()
	defer tc.mu.RUnlock()

	_, exists := tc.capsules[key]
	return exists
}

// WaitForUnlock blocks until a capsule is unlocked or context is canceled
func (tc *MemoryTimeCapsule[T]) WaitForUnlock(ctx context.Context, key string) (T, error) {
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
