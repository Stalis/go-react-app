package requestlog

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"go-react-app/server/util/logger"
)

type RequestLogger struct {
	log *logger.Logger
}

func New(l *logger.Logger) *RequestLogger {
	return &RequestLogger{l}
}

func (l *RequestLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		le := &logEntry{
			ReceivedTime:      start,
			RequestMethod:     r.Method,
			RequestURL:        r.URL.String(),
			RequestHeaderSize: headerSize(r.Header),
			UserAgent:         r.UserAgent(),
			Referer:           r.Referer(),
			Proto:             r.Proto,
			RemoteIP:          ipFromHostPort(r.RemoteAddr),
		}

		if addr, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr); ok {
			le.ServerIP = ipFromHostPort(addr.String())
		}
		r2 := new(http.Request)
		*r2 = *r

		// read request body
		bodyBytes, _ := ioutil.ReadAll(r.Body)

		r.Body.Close()
		// return new io.ReadCloser to Body
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		rcc := &readCounterCloser{r: r.Body}
		r2.Body = rcc
		w2 := &responseStats{w: w}

		defer func() {
			le.Latency = time.Since(start)
			if rcc.err == nil && rcc.r != nil {
				// If the handler hasn't encountered an error in the Body (like EOF),
				// then consume the rest of the Body to provide an accurate rcc.n.
				io.Copy(ioutil.Discard, rcc)
			}
			le.RequestBodySize = rcc.n
			le.Status = w2.code
			if le.Status == 0 {
				le.Status = http.StatusOK
			}
			le.ResponseHeaderSize, le.ResponseBodySize = w2.size()
			l.log.Info().
				Time("received_time", le.ReceivedTime).
				Str("method", le.RequestMethod).
				Str("url", le.RequestURL).
				Int64("header_size", le.RequestHeaderSize).
				Int64("body_size", le.RequestBodySize).
				Str("agent", le.UserAgent).
				Str("referer", le.Referer).
				Str("proto", le.Proto).
				Str("remote_ip", le.RemoteIP).
				Str("server_ip", le.ServerIP).
				Int("status", le.Status).
				Int64("resp_header_size", le.ResponseHeaderSize).
				Int64("resp_body_size", le.ResponseBodySize).
				Dur("latency", le.Latency).
				RawJSON("request", bodyBytes).
				Msg("")
		}()

		next.ServeHTTP(w2, r2)
	})
}
