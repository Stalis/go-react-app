package middlewares

import (
	"go-react-app/internal/util/logger"
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog/hlog"
)

type requestlogger struct {
	log *logger.Logger
}

func NewRequestLogger(log *logger.Logger) Middleware {
	return &requestlogger{log}
}

func (r *requestlogger) Middleware(next http.Handler) http.Handler {
	c := alice.New()

	c = c.Append(hlog.NewHandler(r.log.Output(nil)))
	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))
	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

	return c.Then(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		hlog.FromRequest(r).
			Info().
			Msg("")

		next.ServeHTTP(rw, r)
	}))

}
