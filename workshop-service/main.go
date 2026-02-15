package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	greeting := map[string]string{
		"message": "Welcome to the ALM Workshop!",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Encode the struct to JSON and write it to the response
	json.NewEncoder(w).Encode(greeting)
}

func main() {
	http.HandleFunc("/", defaultHandler)
	// Default endpoint that we can use for Kubernetes health checks
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct{ Status string }{Status: "OK"})
	})
	// Define the /workshop endpoint
	http.HandleFunc("/workshop", WorkshopHandler)

	// Start the server on port 3000
	slog.Info("Starting server on localhost:3000")
	http.ListenAndServe(":3000", nil)
	slog.Info("Exiting server")
}
