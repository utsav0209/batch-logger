package middleware

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			log.Infof("Received request with method: %s, route: %s and served in %s", r.Method, r.RequestURI, time.Since(start))
		}()

		next.ServeHTTP(w, r)
	})
}
