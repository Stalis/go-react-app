package middlewares

import "github.com/gorilla/mux"

func Apply(router *mux.Router) {
	router.Use(LoggingMiddleware)
}
