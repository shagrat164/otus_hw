package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(logg Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		latency := time.Since(start)
		logg.Info(fmt.Sprintf(
			"%s [%s] %s %s %s %d %d %s",
			r.Host,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.URL.Path,
			r.Proto,
			http.StatusOK,
			latency.Milliseconds(),
			r.UserAgent(),
		))
	})
}
