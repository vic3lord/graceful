package graceful

import (
	"errors"
	"net/http"
)

func ExampleDrain() {
	srv := &http.Server{Addr: ":3000"}
	go func() {
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			// Handle error here.
		}
	}()
	<-Drain(&DrainHTTP{srv})
}
