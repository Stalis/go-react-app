package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func Route(r *mux.Router) {
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	r.HandleFunc("/hello", HelloApiHandler)
}

func HelloApiHandler(w http.ResponseWriter, r *http.Request) {
	type HelloResponse struct {
		Message string `json:"message"`
	}

	response := &HelloResponse{
		Message: "Hello, World!",
	}

	formatted, _ := json.Marshal(response)
	w.Header().Add("Content-Type", "application/json")
	w.Write(formatted)
}
