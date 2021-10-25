package middlewares

import (
	"encoding/json"
	"net/http"

	"go-react-app/server/util/logger"
)

type recovery struct {
	log *logger.Logger
}

func NewRecovery(log *logger.Logger) Middleware {
	return &recovery{log}
}

func (m *recovery) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				m.log.Error().Interface("panic", err).Msgf("Panic recover")

				jsonBody, _ := json.Marshal(map[string]string{
					"error": "There was an internal server error",
				})

				rw.WriteHeader(http.StatusInternalServerError)
				rw.Header().Set("Content-Type", "application/json")
				rw.Write(jsonBody)
			}
		}()

		next.ServeHTTP(rw, r)
	})
}
