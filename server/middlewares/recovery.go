package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/phuslu/log"
)

type recovery struct {
}

func (m *recovery) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Error().Interface("panic", err).Msgf("Panic recover")

				jsonBody, _ := json.Marshal(map[string]string{
					"error": "There was an internal server error",
				})

				rw.Header().Set("Content-Type", "application/json")
				rw.WriteHeader(http.StatusInternalServerError)
				rw.Write(jsonBody)
			}
		}()

		next.ServeHTTP(rw, r)
	})
}
