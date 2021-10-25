package middlewares

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/Stalis/go-react-app/server/util/logger"
)

type logging struct {
	log *logger.Logger
}

func (m *logging) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := ioutil.ReadAll(r.Body)

		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		m.log.Debug().
			RawJSON("request", bodyBytes).
			Str("uri", r.RequestURI).
			Msg("")
		next.ServeHTTP(rw, r)
	})
}
