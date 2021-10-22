package middlewares

import (
	"log"
	"net/http"
)

type logger struct {
	l *log.Logger
}

func (m *logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		m.l.Println(r.RequestURI)
		next.ServeHTTP(rw, r)
	})
}
