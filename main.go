package main

import (
	logHandler "batch-logger/pkg/log"
	customMiddleware "batch-logger/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {

	// Create the router
	r := chi.NewRouter()

	// Use Custom Middleware for request logging
	r.Use(middleware.RequestID)
	r.Use(customMiddleware.RequestLogger)

	r.Use(middleware.Recoverer)

	// Route Endpoints
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			log.Errorf("Response failed with error: %s", err)
			return
		}
	})

	// Create Log Handler
	h := logHandler.CreateHandler()

	r.Route("/log", func(r chi.Router) {
		r.Use(customMiddleware.AddLogToCtx)
		r.Post("/", h.HandleLogRequest)
	})

	go h.SyncAtIntervals()

	// Start the server
	log.Infof("Starting server on port: %d", 3000)
	log.Fatal(http.ListenAndServe(":3000", r))
}
