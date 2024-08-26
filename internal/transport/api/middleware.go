package api

import (
	"fmt"
	"net/http"

	"github.com/saikrir/keep-notes/internal/logger"
)

func JSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("-> %s %s", r.Method, r.URL.Path))
		next.ServeHTTP(w, r)
	})
}
