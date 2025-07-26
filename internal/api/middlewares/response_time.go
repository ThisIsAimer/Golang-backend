package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

func ResponseTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// create a custom response writer to capture the status code
		wrappedWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(start)
		//wrappedWriter.Header().Set("X-Response-Time",duration.String())
		fmt.Printf("method: %s, url: %s, status: %d, duration: %v\n", r.Method, r.URL, wrappedWriter.status, duration.String())
	})
}

type responseWriter struct {
	// adds all the properties of http.ResponseWriter to our struct
	http.ResponseWriter
	status int
}

// we are overwriting WriteHeader method
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code // now we can log the status codes
	rw.ResponseWriter.WriteHeader(code)
}
