package middlewares

import (
	"bytes"
	"go-react-app/server/util/logger"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type requestBodyLogger struct {
	log *logger.Logger
}

func NewRequestBodyLogger(log *logger.Logger) Middleware {
	return &requestBodyLogger{log}
}

func (m *requestBodyLogger) Middleware(next http.Handler) http.Handler {
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
	})
}
