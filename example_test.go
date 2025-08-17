package timecapsule

import (
	"context"
	"fmt"
	"time"
)

// Example_basic demonstrates basic usage of the time capsule
func Example_basic() {
	// Create a new time capsule
	capsule := New[string]()

	// Store a value that unlocks in 1 second
	unlockTime := time.Now().Add(1 * time.Second)
	err := capsule.Store(context.Background(), "greeting", "Hello, World!", unlockTime)
	if err != nil {
		panic(err)
	}

	// Try to open immediately - should be locked
	if _, openErr := capsule.Open(context.Background(), "greeting"); openErr != nil {
		fmt.Println("Capsule is locked:", openErr)
	}

	// Wait for unlock
	time.Sleep(2 * time.Second)

	// Now open the capsule
	value, err := capsule.Open(context.Background(), "greeting")
	if err != nil {
		panic(err)
	}
	fmt.Println("Unlocked value:", value)

	// Output:
	// Capsule is locked: capsule is still locked
	// Unlocked value: Hello, World!
}

// Example_struct demonstrates using structs with the time capsule
func Example_struct() {
	type Promo struct {
		Code     string `json:"code"`
		Discount int    `json:"discount"`
		Valid    bool   `json:"valid"`
	}

	// Create a time capsule for Promo structs
	capsule := New[Promo]()

	// Store a promotional code that unlocks tomorrow
	tomorrow := time.Now().Add(24 * time.Hour)
	promo := Promo{
		Code:     "HOLIDAY50",
		Discount: 50,
		Valid:    true,
	}

	err := capsule.Store(context.Background(), "holiday-sale", promo, tomorrow)
	if err != nil {
		panic(err)
	}

	// Check metadata without opening
	metadata, err := capsule.Peek(context.Background(), "holiday-sale")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Promo unlocks in: %v\n", time.Until(metadata.UnlockTime).Round(time.Hour))
	fmt.Printf("Is locked: %v\n", metadata.IsLocked)

	// Output:
	// Promo unlocks in: 24h0m0s
	// Is locked: true
}

// Example_waitForUnlock demonstrates waiting for a capsule to unlock
func Example_waitForUnlock() {
	capsule := New[int]()

	// Store a value that unlocks in 100ms
	unlockTime := time.Now().Add(100 * time.Millisecond)
	err := capsule.Store(context.Background(), "count", 42, unlockTime)
	if err != nil {
		panic(err)
	}

	// Wait for unlock with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	value, err := capsule.WaitForUnlock(ctx, "count")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Waited for unlock, got value: %d\n", value)

	// Output:
	// Waited for unlock, got value: 42
}

// Example_context demonstrates using context for cancellation
func Example_context() {
	capsule := New[string]()

	// Store a value that unlocks in 5 seconds
	unlockTime := time.Now().Add(5 * time.Second)
	err := capsule.Store(context.Background(), "slow", "This takes time", unlockTime)
	if err != nil {
		panic(err)
	}

	// Create a context that cancels after 1 second
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Try to wait for unlock, but context will cancel first
	_, err = capsule.WaitForUnlock(ctx, "slow")
	if err != nil {
		fmt.Println("Context canceled:", err)
	}

	// Output:
	// Context canceled: context deadline exceeded
}

// Example_management demonstrates capsule management operations
func Example_management() {
	capsule := New[float64]()

	// Store multiple values
	now := time.Now()
	var err error
	err = capsule.Store(context.Background(), "price1", 19.99, now.Add(1*time.Hour))
	if err != nil {
		panic(err)
	}
	err = capsule.Store(context.Background(), "price2", 29.99, now.Add(2*time.Hour))
	if err != nil {
		panic(err)
	}

	// Check if capsules exist
	fmt.Printf("price1 exists: %v\n", capsule.Exists(context.Background(), "price1"))
	fmt.Printf("price3 exists: %v\n", capsule.Exists(context.Background(), "price3"))

	// Delete a capsule
	err = capsule.Delete(context.Background(), "price1")
	if err != nil {
		panic(err)
	}

	fmt.Printf("After delete, price1 exists: %v\n", capsule.Exists(context.Background(), "price1"))

	// Output:
	// price1 exists: true
	// price3 exists: false
	// After delete, price1 exists: false
}

// Example_delay demonstrates delaying a capsule's unlock time
func Example_delay() {
	capsule := New[string]()
	ctx := context.Background()

	// Store a value that unlocks in 1 hour
	unlockTime := time.Now().Add(1 * time.Hour)
	err := capsule.Store(ctx, "secret", "confidential", unlockTime)
	if err != nil {
		panic(err)
	}

	// Check initial unlock time
	metadata, _ := capsule.Peek(ctx, "secret")
	fmt.Printf("Original unlock time: %v\n", metadata.UnlockTime.Format("15:04"))

	// Delay by 2 more hours
	err = capsule.Delay(ctx, "secret", 2*time.Hour)
	if err != nil {
		panic(err)
	}

	// Check new unlock time
	metadata, _ = capsule.Peek(ctx, "secret")
	fmt.Printf("New unlock time: %v\n", metadata.UnlockTime.Format("15:04"))

	// Output:
	// Original unlock time: 16:30
	// New unlock time: 17:30
}
