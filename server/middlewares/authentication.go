package middlewares

import (
	"go-react-app/server/dal"
	"go-react-app/server/util/logger"
	"net/http"
	"net/url"
	"strings"

	gsessions "github.com/gorilla/sessions"
	"github.com/pkg/errors"
)

type authentication struct {
	log      *logger.Logger
	sessions dal.SessionRepository
	store    gsessions.Store
}

func NewAuthentication(log *logger.Logger, sessions dal.SessionRepository) *authentication {
	return &authentication{log, sessions, gsessions.NewCookieStore([]byte{})}
}

func (auth *authentication) CheckToken(r *http.Request) error {
	uri := r.URL.RequestURI()
	parsedUrl, _ := url.ParseRequestURI(uri)
	parts := strings.Split(parsedUrl.Path, "/")
	auth.log.Debug().Interface("parsed", parsedUrl).Msg("")
	switch parts[len(parts)-1] {
	case "login":
	case "register":
		return nil
	}

	session, err := auth.store.Get(r, "myapp-session")
	if err != nil {
		return errors.Wrap(err, "token not found")
	}

	userId := 0
	if session.Values["userId"] != nil {
		userId = session.Values["userId"].(int)
	}
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		return errors.New("not authenticated")
	}

	auth.log.Debug().Msgf("Authenticated user: %d", userId)
	return nil
}

func (auth *authentication) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		if err := auth.CheckToken(r); err != nil {
			err = errors.Wrap(err, "error while authenticate")

			auth.log.Error().Stack().Caller().Err(err).Msg("")
			http.Error(rw, err.Error(), http.StatusForbidden)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
