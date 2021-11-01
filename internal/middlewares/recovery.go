package middlewares

import (
	"encoding/json"
	"net/http"

	"go-react-app/internal/util/logger"

	"github.com/gorilla/mux"
)

func NewRecovery(log *logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return recoveryInternal(log, next)
	}
}

func recoveryInternal(log *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Error().Interface("panic", err).Msgf("Panic recover")

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
