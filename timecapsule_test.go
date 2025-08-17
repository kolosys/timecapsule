package timecapsule

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	capsule := New[string]()
	assert.NotNil(t, capsule)

	// Test that it's empty initially
	assert.False(t, capsule.Exists(context.Background(), "test"))
}

func TestStoreAndOpen(t *testing.T) {
	capsule := New[string]()
	ctx := context.Background()

	// Store a value that unlocks in the future
	unlockTime := time.Now().Add(100 * time.Millisecond)
	err := capsule.Store(ctx, "test", "hello", unlockTime)
	require.NoError(t, err)

	// Verify it exists
	assert.True(t, capsule.Exists(ctx, "test"))

	// Try to open immediately - should be locked
	_, err = capsule.Open(ctx, "test")
	assert.ErrorIs(t, err, ErrCapsuleLocked)

	// Wait for unlock
	time.Sleep(150 * time.Millisecond)

	// Now open the capsule
	value, err := capsule.Open(ctx, "test")
	require.NoError(t, err)
	assert.Equal(t, "hello", value)
}

func TestStoreAndOpenImmediate(t *testing.T) {
	capsule := New[string]()
	ctx := context.Background()

	// Store a value that unlocks immediately
	unlockTime := time.Now().Add(-1 * time.Second)
	err := capsule.Store(ctx, "test", "hello", unlockTime)
	require.NoError(t, err)

	// Should be able to open immediately
	value, err := capsule.Open(ctx, "test")
	require.NoError(t, err)
	assert.Equal(t, "hello", value)
}

func TestOpenNonExistent(t *testing.T) {
	capsule := New[string]()
	ctx := context.Background()

	_, err := capsule.Open(ctx, "nonexistent")
	assert.ErrorIs(t, err, ErrCapsuleNotFound)
}

func TestStoreInvalidKey(t *testing.T) {
	capsule := New[string]()
	ctx := context.Background()

	err := capsule.Store(ctx, "", "hello", time.Now().Add(time.Hour))
	assert.ErrorIs(t, err, ErrInvalidKey)
}

func TestOpenInvalidKey(t *testing.T) {
	capsule := New[string]()
	ctx := context.Background()

	_, err := capsule.Open(ctx, "")
	assert.ErrorIs(t, err, ErrInvalidKey)
}

func TestPeek(t *testing.T) {
	capsule := New[string]()
	ctx := context.Background()

	// Store a value
	unlockTime := time.Now().Add(time.Hour)
	err := capsule.Store(ctx, "test", "hello", unlockTime)
	require.NoError(t, err)

	// Peek at metadata
	metadata, err := capsule.Peek(ctx, "test")
	require.NoError(t, err)
	assert.True(t, metadata.IsLocked)
	assert.Equal(t, unlockTime.Unix(), metadata.UnlockTime.Unix())
	assert.WithinDuration(t, time.Now(), metadata.CreatedAt, 2*time.Second)

	// Peek at non-existent capsule
	_, err = capsule.Peek(ctx, "nonexistent")
	assert.ErrorIs(t, err, ErrCapsuleNotFound)
}

func TestDelete(t *testing.T) {
	capsule := New[string]()
	ctx := context.Background()

	// Store a value
	err := capsule.Store(ctx, "test", "hello", time.Now().Add(time.Hour))
	require.NoError(t, err)

	// Verify it exists
	assert.True(t, capsule.Exists(ctx, "test"))

	// Delete it
	err = capsule.Delete(ctx, "test")
	require.NoError(t, err)

	// Verify it's gone
	assert.False(t, capsule.Exists(ctx, "test"))

	// Try to delete non-existent
	err = capsule.Delete(ctx, "nonexistent")
	assert.ErrorIs(t, err, ErrCapsuleNotFound)
}

func TestExists(t *testing.T) {
	capsule := New[string]()
	ctx := context.Background()

	// Initially doesn't exist
	assert.False(t, capsule.Exists(ctx, "test"))

	// Store a value
	err := capsule.Store(ctx, "test", "hello", time.Now().Add(time.Hour))
	require.NoError(t, err)

	// Now exists
	assert.True(t, capsule.Exists(ctx, "test"))

	// Delete it
	err = capsule.Delete(ctx, "test")
	require.NoError(t, err)

	// No longer exists
	assert.False(t, capsule.Exists(ctx, "test"))
}

func TestWaitForUnlock(t *testing.T) {
	capsule := New[string]()
	ctx := context.Background()

	// Store a value that unlocks in 50ms
	unlockTime := time.Now().Add(50 * time.Millisecond)
	err := capsule.Store(ctx, "test", "hello", unlockTime)
	require.NoError(t, err)

	// Wait for unlock
	value, err := capsule.WaitForUnlock(ctx, "test")
	require.NoError(t, err)
	assert.Equal(t, "hello", value)
}

func TestWaitForUnlockAlreadyUnlocked(t *testing.T) {
	capsule := New[string]()
	ctx := context.Background()

	// Store a value that's already unlocked
	unlockTime := time.Now().Add(-1 * time.Second)
	err := capsule.Store(ctx, "test", "hello", unlockTime)
	require.NoError(t, err)

	// Should return immediately
	value, err := capsule.WaitForUnlock(ctx, "test")
	require.NoError(t, err)
	assert.Equal(t, "hello", value)
}

func TestWaitForUnlockContextCancelled(t *testing.T) {
	capsule := New[string]()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Store a value that unlocks in 1 second
	unlockTime := time.Now().Add(time.Second)
	err := capsule.Store(ctx, "test", "hello", unlockTime)
	require.NoError(t, err)

	// Context should cancel before unlock
	_, err = capsule.WaitForUnlock(ctx, "test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "deadline exceeded")
}

func TestWaitForUnlockNonExistent(t *testing.T) {
	capsule := New[string]()
	ctx := context.Background()

	_, err := capsule.WaitForUnlock(ctx, "nonexistent")
	assert.ErrorIs(t, err, ErrCapsuleNotFound)
}

func TestConcurrentAccess(t *testing.T) {
	capsule := New[int]()
	ctx := context.Background()

	// Store multiple values concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			key := fmt.Sprintf("key-%d", id)
			unlockTime := time.Now().Add(time.Duration(id) * time.Millisecond)
			err := capsule.Store(ctx, key, id, unlockTime)
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	// Wait for all stores to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all values exist
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key-%d", i)
		assert.True(t, capsule.Exists(ctx, key))
	}
}

func TestStructValues(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	capsule := New[TestStruct]()
	ctx := context.Background()

	testData := TestStruct{
		Name:  "test",
		Value: 42,
	}

	// Store struct
	unlockTime := time.Now().Add(-1 * time.Second) // Already unlocked
	err := capsule.Store(ctx, "struct", testData, unlockTime)
	require.NoError(t, err)

	// Retrieve struct
	value, err := capsule.Open(ctx, "struct")
	require.NoError(t, err)
	assert.Equal(t, testData, value)
}

func TestDelay(t *testing.T) {
	capsule := New[string]()
	ctx := context.Background()

	// Store a value that unlocks in 1 hour
	unlockTime := time.Now().Add(time.Hour)
	err := capsule.Store(ctx, "test", "hello", unlockTime)
	require.NoError(t, err)

	// Verify it's locked
	metadata, err := capsule.Peek(ctx, "test")
	require.NoError(t, err)
	assert.True(t, metadata.IsLocked)

	// Delay by 2 hours
	err = capsule.Delay(ctx, "test", 2*time.Hour)
	require.NoError(t, err)

	// Verify it's still locked but with new unlock time
	metadata, err = capsule.Peek(ctx, "test")
	require.NoError(t, err)
	assert.True(t, metadata.IsLocked)
	assert.True(t, metadata.UnlockTime.After(unlockTime))

	// Try to delay non-existent capsule
	err = capsule.Delay(ctx, "nonexistent", time.Hour)
	assert.ErrorIs(t, err, ErrCapsuleNotFound)
}

func TestDelayInvalidKey(t *testing.T) {
	capsule := New[string]()
	ctx := context.Background()

	err := capsule.Delay(ctx, "", time.Hour)
	assert.ErrorIs(t, err, ErrInvalidKey)
}

func TestContextCancellation(t *testing.T) {
	capsule := New[string]()
	ctx, cancel := context.WithCancel(context.Background())

	// Store a value
	err := capsule.Store(ctx, "test", "hello", time.Now().Add(time.Hour))
	require.NoError(t, err)

	// Cancel context
	cancel()

	// Operations should fail with context error
	_, err = capsule.Open(ctx, "test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")

	_, err = capsule.Peek(ctx, "test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")

	err = capsule.Delay(ctx, "test", time.Hour)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")

	err = capsule.Delete(ctx, "test")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "context canceled")
}
