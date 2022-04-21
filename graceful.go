// Package graceful provides a server graceful draining helper.
package graceful

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Drainer will drain connections and shutdown gracefully.
type Drainer interface {
	Drain(context.Context) error
}

// Drain takes Drainers, waits for signals or context cancellation
// It will then block until all Drainers were shutdown or timedout.
func Drain(ctx context.Context, ds ...Drainer) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, os.Interrupt)
	defer stop()

	// Wait for signals or context to be done
	<-ctx.Done()

	log.Printf("Shutting down gracefully: %v", ctx.Err())

	// Server Drainer context
	sctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	for _, d := range ds {
		wg.Add(1)
		go func(d Drainer, wg *sync.WaitGroup) {
			defer wg.Done()
			if err := d.Drain(sctx); err != nil {
				log.Printf("Failed to drain: %v", err)
			}
		}(d, &wg)
	}
	wg.Wait()
	return ctx.Err()
}
