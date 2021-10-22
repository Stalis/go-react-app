package middlewares

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/phuslu/log"
)

type logger struct {
}

func (m *logger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		buf := bytes.NewBuffer(bodyBytes)

		r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		log.Debug().
			RawJSONStr("request", buf.String()).
			Str("uri", r.RequestURI).
			Msg("")
		next.ServeHTTP(rw, r)
	})
}
