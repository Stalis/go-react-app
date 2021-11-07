package middlewares

import (
	"bytes"
	"go-react-app/internal/util/logger"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type RequestBodyLogger struct {
	log *logger.Logger
}

func NewRequestBodyLogger(log *logger.Logger) *RequestBodyLogger {
	return &RequestBodyLogger{log}
}

func (m *RequestBodyLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			err := errors.Wrap(err, "unable to read request body")
			m.log.Error().Stack().Caller().Err(err).Msg("")
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		m.log.Debug().RawJSON("request", bodyBytes).Msg("")

		next.ServeHTTP(rw, r)
	})
}
