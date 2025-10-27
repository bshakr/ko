package signals

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// SetupCancellableContext creates a context that can be cancelled by interrupt signals (Ctrl+C, SIGTERM).
// It returns the context and a cleanup function that must be called (typically via defer) to:
// - Stop signal notifications
// - Cancel the context
// - Wait for the signal handler goroutine to exit
//
// Example usage:
//
//	ctx, cleanup := signals.SetupCancellableContext()
//	defer cleanup()
//
//	// Use ctx for long-running operations
//	if err := someOperation(ctx); err != nil {
//	    return err
//	}
func SetupCancellableContext() (context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())

	// Setup signal channel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Signal handler goroutine
	done := make(chan struct{})
	go func() {
		defer close(done)
		select {
		case <-sigChan:
			fmt.Println("\nOperation cancelled by user")
			cancel()
		case <-ctx.Done():
			// Context cancelled or operation completed, exit goroutine
		}
	}()

	// Cleanup function
	cleanup := func() {
		signal.Stop(sigChan) // Stop receiving signals
		cancel()             // Signal goroutine to exit if still running
		<-done               // Wait for signal handler to finish
	}

	return ctx, cleanup
}
