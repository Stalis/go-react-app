package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/Stalis/go-react-app/api/handlers/account"
	"github.com/Stalis/go-react-app/config"
	"github.com/Stalis/go-react-app/dal"
	"github.com/gorilla/mux"
)

func Route(r *mux.Router, cfg *config.Config) {
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	db, err := dal.ConnectDB(cfg.Database.Url)
	if err != nil {
		log.Println(err.Error())
		os.Exit(-1)
	}

	accountRouter := r.PathPrefix("/account").Subrouter()
	accountRouter.Handle("/login", account.NewLogin(db))
	accountRouter.Handle("/register", account.NewRegister(db))

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
