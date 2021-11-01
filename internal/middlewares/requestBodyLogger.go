package middlewares

import (
	"bytes"
	"go-react-app/internal/util/logger"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func NewRequestBodyLogger(log *logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return requestBodyLoggerInternal(log, next)
	}
}

func requestBodyLoggerInternal(log *logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			err := errors.Wrap(err, "unable to read request body")
			log.Error().Stack().Caller().Err(err).Msg("")
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		log.Debug().RawJSON("request", bodyBytes).Msg("")
	})
}
