package main

import "net/http"

var (
	sem     chan struct{}
	acquire = func() { sem <- struct{}{} }
	release = func() { <-sem }
)

// MaxAllowedMiddleware
func MaxAllowedMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acquire()
		defer release()
		next.ServeHTTP(w, r)
	})
}
