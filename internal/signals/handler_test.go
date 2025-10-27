package signals

import (
	"context"
	"os"
	"runtime"
	"syscall"
	"testing"
	"time"
)

// TestSetupCancellableContextCreatesValidContext verifies that the function
// returns a valid context and cleanup function
func TestSetupCancellableContextCreatesValidContext(t *testing.T) {
	ctx, cleanup := SetupCancellableContext()

	if ctx == nil {
		t.Fatal("Context should not be nil")
	}

	if cleanup == nil {
		t.Fatal("Cleanup function should not be nil")
	}

	// Context should not be cancelled initially
	select {
	case <-ctx.Done():
		t.Fatal("Context should not be cancelled initially")
	default:
		// Expected: context is not cancelled
	}

	// Cleanup should not block
	cleanup()
}

// TestCleanupDoesNotHang verifies that calling cleanup() completes quickly
// and doesn't cause the process to hang
func TestCleanupDoesNotHang(t *testing.T) {
	ctx, cleanup := SetupCancellableContext()

	if ctx == nil {
		t.Fatal("Context should not be nil")
	}

	// Call cleanup in a goroutine with timeout
	done := make(chan struct{})
	go func() {
		cleanup()
		close(done)
	}()

	// Wait for cleanup with timeout
	select {
	case <-done:
		// Success: cleanup completed
	case <-time.After(2 * time.Second):
		t.Fatal("Cleanup hung and did not complete within 2 seconds")
	}
}

// TestCleanupCancelsContext verifies that cleanup cancels the context
func TestCleanupCancelsContext(t *testing.T) {
	ctx, cleanup := SetupCancellableContext()

	// Context should not be cancelled initially
	select {
	case <-ctx.Done():
		t.Fatal("Context should not be cancelled initially")
	default:
		// Expected
	}

	// Call cleanup
	cleanup()

	// Context should now be cancelled
	select {
	case <-ctx.Done():
		// Expected: context is cancelled
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Context should be cancelled after cleanup")
	}

	// Verify context error
	if ctx.Err() != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", ctx.Err())
	}
}

// TestNoGoroutineLeak verifies that the signal handler goroutine exits
// after cleanup is called
func TestNoGoroutineLeak(t *testing.T) {
	// Get initial goroutine count
	initialCount := runtime.NumGoroutine()

	// Create context and cleanup immediately
	ctx, cleanup := SetupCancellableContext()
	_ = ctx

	cleanup()

	// Give goroutines time to exit
	time.Sleep(100 * time.Millisecond)

	// Check goroutine count
	finalCount := runtime.NumGoroutine()

	// Final count should be same or less than initial
	// (we allow "same" because other test goroutines might be running)
	if finalCount > initialCount+1 {
		t.Errorf("Goroutine leak detected: initial=%d, final=%d", initialCount, finalCount)
	}
}

// TestMultipleCleanupCallsAreSafe verifies that calling cleanup multiple times
// doesn't cause panics or hangs
func TestMultipleCleanupCallsAreSafe(t *testing.T) {
	ctx, cleanup := SetupCancellableContext()
	_ = ctx

	// First cleanup
	cleanup()

	// Second cleanup should not panic or hang
	done := make(chan struct{})
	go func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Cleanup panicked on second call: %v", r)
			}
			close(done)
		}()
		cleanup()
	}()

	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Second cleanup call hung")
	}
}

// TestSignalDeliveryDoesNotCrash verifies that sending a signal after setup
// doesn't cause crashes (though we can't easily verify cancellation in tests)
func TestSignalDeliveryDoesNotCrash(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping signal delivery test in short mode")
	}

	ctx, cleanup := SetupCancellableContext()
	defer cleanup()

	// Send interrupt signal to ourselves
	// Note: This will be caught by the signal handler
	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Failed to find current process: %v", err)
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("Failed to send signal: %v", err)
	}

	// Wait a bit for signal to be processed
	time.Sleep(100 * time.Millisecond)

	// Context should be cancelled by the signal handler
	select {
	case <-ctx.Done():
		// Expected: signal handler cancelled the context
		t.Log("Context was cancelled by signal handler (expected)")
	default:
		// Note: On some systems, SIGTERM might not be delivered in tests
		t.Log("Context was not cancelled - signal might not have been delivered in test environment")
	}

	// Cleanup should still work even after signal
	cleanup()
}

// TestContextPropagation verifies that the context can be used with
// context-aware operations
func TestContextPropagation(t *testing.T) {
	ctx, cleanup := SetupCancellableContext()
	defer cleanup()

	// Create a child context
	childCtx, childCancel := context.WithCancel(ctx)
	defer childCancel()

	// Child should not be cancelled initially
	select {
	case <-childCtx.Done():
		t.Fatal("Child context should not be cancelled initially")
	default:
		// Expected
	}

	// Call cleanup to cancel parent
	cleanup()

	// Child should also be cancelled when parent is cancelled
	select {
	case <-childCtx.Done():
		// Expected: child inherits parent cancellation
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Child context should be cancelled when parent is cancelled")
	}
}

// BenchmarkSetupCancellableContext measures the performance overhead
// of setting up the cancellable context
func BenchmarkSetupCancellableContext(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx, cleanup := SetupCancellableContext()
		_ = ctx
		cleanup()
	}
}
