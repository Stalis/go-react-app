package middlewares

import (
	"log"

	"github.com/gorilla/mux"
)

func Apply(router *mux.Router, l *log.Logger) {

	logging := &loggingMiddleware{l}
	router.Use(logging.Middleware)
}
