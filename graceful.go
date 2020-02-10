// Package graceful provides a server graceful draining helper.
package graceful

import "context"

// Drainer will drain connections and shutdown gracefully.
type Drainer interface {
	Drain() <-chan struct{}
	DrainWithContext(context.Context) <-chan struct{}
}

// Drain takes Drainer and returns a future to hold until fully drained.
func Drain(d Drainer) <-chan struct{} {
	return d.Drain()
}

// DrainWithContext behaves the same as Drain but can be canceled with ctx.
func DrainWithContext(ctx context.Context, d Drainer) <-chan struct{} {
	return d.DrainWithContext(ctx)
}
