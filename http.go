// Package graceful provides a server graceful draining helper.
package graceful

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type HTTPDrain struct {
	*http.Server
}

// Drain takes a http.Server and drains on certain os.Signals
// and returns a future to block the server until fully drained.
func (d *HTTPDrain) Drain() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		q := make(chan os.Signal, 1)
		signal.Notify(q, syscall.SIGTERM, os.Interrupt)
		<-q
		log.Println("[graceful] Starting shutdown sequence")
		if err := d.Shutdown(context.Background()); err != nil {
			log.Printf("[graceful] Could not shutdown gracefully: %v", err)
		}
		close(done)
	}()
	return done
}

// DrainWithContext behaves the same as Drain but can be canceled with ctx.
func (d *HTTPDrain) DrainWithContext(ctx context.Context) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		q := make(chan os.Signal, 1)
		signal.Notify(q, syscall.SIGTERM, os.Interrupt)

		// Wait for signal either from ctx or q
		select {
		case signal := <-q:
			log.Printf("[graceful] Shutting down due to signal: %v", signal)
		case <-ctx.Done():
			log.Printf("[graceful] Shutting down due to ctx: %v", ctx.Err())
		}

		if err := d.Shutdown(context.Background()); err != nil {
			log.Printf("[graceful] Could not shutdown gracefully: %v", err)
		}
		close(done)
	}()
	return done
}
