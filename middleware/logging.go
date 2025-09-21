package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body = append(rw.body, b...)
	return rw.ResponseWriter.Write(b)
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log request
		var requestBody []byte
		if r.Body != nil {
			requestBody, _ = io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Wrap response writer to capture response
		rw := &responseWriter{
			ResponseWriter: w,
			status:        http.StatusOK,
		}

		// Process request
		next.ServeHTTP(rw, r)

		// Log response
		duration := time.Since(start)
		log.Printf(
			"[%s] %s %s | Status: %d | Duration: %v | Request: %s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			rw.status,
			duration,
			string(requestBody),
		)

		// Log response body for errors
		if rw.status >= 400 {
			log.Printf("Response Error: %s", string(rw.body))
		}
	})
}