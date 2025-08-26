package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/kolosys/timecapsule"
)

// These variables are set during build time via ldflags
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

type PromoCode struct {
	Code     string `json:"code"`
	Discount int    `json:"discount"`
	Valid    bool   `json:"valid"`
}

func main() {
	// Check for version flag
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("timecapsule version %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Build Date: %s\n", date)
		return
	}

	fmt.Println("🚀 timecapsule Demo")
	fmt.Println("=====================")

	// Create a new time capsule
	capsule := timecapsule.New[PromoCode]()
	ctx := context.Background()

	// Demo 1: Store a promotional code that unlocks in 3 seconds
	fmt.Println("\n1. Storing a promotional code that unlocks in 3 seconds...")
	promo := PromoCode{
		Code:     "HOLIDAY50",
		Discount: 50,
		Valid:    true,
	}

	unlockTime := time.Now().Add(3 * time.Second)
	err := capsule.Store(ctx, "holiday-sale", promo, unlockTime)
	if err != nil {
		fmt.Printf("❌ Failed to store promo: %v\n", err)
		return
	}

	// Demo 2: Try to peek at the capsule
	fmt.Println("\n2. Peeking at the capsule metadata...")
	metadata, err := capsule.Peek(ctx, "holiday-sale")
	if err != nil {
		fmt.Printf("❌ Failed to peek: %v\n", err)
		return
	}

	fmt.Printf("   Unlock time: %v\n", metadata.UnlockTime.Format("15:04:05"))
	fmt.Printf("   Is locked: %v\n", metadata.IsLocked)
	fmt.Printf("   Time until unlock: %v\n", time.Until(metadata.UnlockTime).Round(time.Second))

	// Demo 3: Try to open the capsule (should fail)
	fmt.Println("\n3. Trying to open the capsule now (should fail)...")
	if _, openErr := capsule.Open(ctx, "holiday-sale"); openErr != nil {
		fmt.Printf("   ❌ Expected error: %v\n", openErr)
	}

	// Demo 4: Wait for unlock
	fmt.Println("\n4. Waiting for the capsule to unlock...")
	value, err := capsule.WaitForUnlock(ctx, "holiday-sale")
	if err != nil {
		fmt.Printf("❌ Failed to wait for unlock: %v\n", err)
		return
	}

	fmt.Printf("   ✅ Unlocked! Promo code: %s, Discount: %d%%\n", value.Code, value.Discount)

	// Demo 5: Store multiple values with different unlock times
	fmt.Println("\n5. Storing multiple values with different unlock times...")

	// Value that unlocks immediately
	err = capsule.Store(ctx, "immediate", PromoCode{Code: "NOW", Discount: 10, Valid: true}, time.Now().Add(-1*time.Second))
	if err != nil {
		fmt.Printf("❌ Failed to store immediate: %v\n", err)
		return
	}

	// Value that unlocks in 2 seconds
	err = capsule.Store(ctx, "soon", PromoCode{Code: "SOON", Discount: 20, Valid: true}, time.Now().Add(2*time.Second))
	if err != nil {
		fmt.Printf("❌ Failed to store soon: %v\n", err)
		return
	}

	// Value that unlocks in 5 seconds
	err = capsule.Store(ctx, "later", PromoCode{Code: "LATER", Discount: 30, Valid: true}, time.Now().Add(5*time.Second))
	if err != nil {
		fmt.Printf("❌ Failed to store later: %v\n", err)
		return
	}

	// Demo 6: Check which capsules exist and their status
	fmt.Println("\n6. Checking capsule status...")
	keys := []string{"immediate", "soon", "later", "holiday-sale"}

	for _, key := range keys {
		if capsule.Exists(ctx, key) {
			meta, _ := capsule.Peek(ctx, key)
			status := "🔒 LOCKED"
			if !meta.IsLocked {
				status = "🔓 UNLOCKED"
			}
			fmt.Printf("   %s: %s\n", key, status)
		} else {
			fmt.Printf("   %s: ❌ NOT FOUND\n", key)
		}
	}

	// Demo 7: Open immediate capsule
	fmt.Println("\n7. Opening immediate capsule...")
	if value, openErr := capsule.Open(ctx, "immediate"); openErr == nil {
		fmt.Printf("   ✅ %s: %d%% discount\n", value.Code, value.Discount)
	}

	// Demo 8: Wait for "soon" capsule with timeout
	fmt.Println("\n8. Waiting for 'soon' capsule with timeout...")
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if value, waitErr := capsule.WaitForUnlock(timeoutCtx, "soon"); waitErr == nil {
		fmt.Printf("   ✅ %s: %d%% discount\n", value.Code, value.Discount)
	} else {
		fmt.Printf("   ❌ Timeout: %v\n", waitErr)
	}

	// Demo 9: Delay a capsule
	fmt.Println("\n9. Delaying the 'later' capsule...")
	originalMeta, _ := capsule.Peek(ctx, "later")
	fmt.Printf("   Original unlock time: %v\n", originalMeta.UnlockTime.Format("15:04:05"))

	err = capsule.Delay(ctx, "later", 1*time.Hour)
	if err != nil {
		fmt.Printf("❌ Failed to delay capsule: %v\n", err)
		return
	}

	newMeta, _ := capsule.Peek(ctx, "later")
	fmt.Printf("   New unlock time: %v\n", newMeta.UnlockTime.Format("15:04:05"))

	// Demo 10: Clean up
	fmt.Println("\n10. Cleaning up...")
	for _, key := range keys {
		if err := capsule.Delete(ctx, key); err == nil {
			fmt.Printf("   ✅ Deleted: %s\n", key)
		}
	}

	fmt.Println("\n🎉 Demo completed successfully!")
}
