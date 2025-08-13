package app

import (
	"aitring/handlers"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// CORS middleware to set the headers
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from specific frontend origins
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace with your production client URL

		// Allow specific methods (GET, POST, OPTIONS, etc.)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Allow specific headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func runServer(envPort string, h handlers.Store) {

	// Ensure envPort has a valid default value
	if envPort == "" {
		envPort = "8080"
	}

	// Create a new router
	r := mux.NewRouter()

	// Public routes

	r.HandleFunc("/public/chunks/{id}", h.AudioHandler.GetChunkByID).Methods(http.MethodGet)

	r.HandleFunc("/public/sessions/{user_id}", h.AudioHandler.GetChunksByUser).Methods(http.MethodGet)

	r.HandleFunc("/public/upload", h.AudioHandler.Upload).Methods(http.MethodPost)

	// WebSocket route
	r.HandleFunc("/ws", h.AudioHandler.WSHandler).Methods(http.MethodGet)

	//    r.HandleFunc("/events",h.FieldsHandler.SSEHandler)

	// Wrap the router with CORS middleware
	http.Handle("/", enableCORS(r))

	// Start the server
	fmt.Printf("Server listening on port %s...\n", envPort)
	log.Fatal(http.ListenAndServe(":"+envPort, nil))
}
