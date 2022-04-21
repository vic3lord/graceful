# Graceful

Graceful shutdown helpers for long running servers.
Databases, HTTP, gRPC, etc.

```go
package main

import (
	"errors"
	"net/http"

	"github.com/vic3lord/graceful"
)

func main() {
	srv := &http.Server{Addr: ":3000"}
	go func() {
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			// Handle error here.
		}
	}()
	graceful.Drain(context.Background(), &graceful.DrainHTTP{srv})
}
```
