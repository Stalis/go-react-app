package middlewares

import (
	"github.com/gorilla/mux"
)

func Apply(router *mux.Router) {

	recovery := &recovery{}
	router.Use(recovery.Middleware)

	logging := &logger{}
	router.Use(logging.Middleware)
}
