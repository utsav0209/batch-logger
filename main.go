package main

import (
	customMiddleware "batch-logger/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

func main() {

	// Create the router
	r := chi.NewRouter()
	r.Use(customMiddleware.RequestLogger())
	r.Use(middleware.Recoverer)

	// Route Endpoints
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			log.Print("Response failed with error: ", err)
		}
	})

	// Start the server
	log.Fatal(http.ListenAndServe(":3000", r))
}
