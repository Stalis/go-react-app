package middlewares

import (
	"log"

	"github.com/gorilla/mux"
)

func Apply(router *mux.Router, l *log.Logger) {
	recovery := &recovery{l}
	router.Use(recovery.Middleware)

	logging := &logger{l}
	router.Use(logging.Middleware)
}
