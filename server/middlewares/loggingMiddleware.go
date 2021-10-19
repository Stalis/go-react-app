package middlewares

import (
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(rw, r)
	})
}
