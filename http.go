// Package graceful provides a server graceful draining helper.
package graceful

import (
	"context"
	"net/http"
)

type DrainHTTP struct {
	*http.Server
}

// Drain takes a http.Server and drains on certain os.Signals
// and returns a future to block the server until fully drained.
func (d *DrainHTTP) Drain(ctx context.Context) error {
	return d.Shutdown(ctx)
}
