package middlewares

import (
	"go-react-app/server/dal"
	"go-react-app/server/util/logger"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

type authentication struct {
	log      *logger.Logger
	sessions dal.SessionRepository
}

func NewAuthentication(log *logger.Logger, sessions dal.SessionRepository) *authentication {
	return &authentication{log, sessions}
}

func (auth *authentication) CheckToken(token string) error {
	if token == "initial" {
		auth.log.Debug().Msgf("Initial user")
		return nil
	}

	session, err := auth.sessions.GetSessionByToken(uuid.FromStringOrNil(token))
	if err != nil {
		return err
	}

	auth.log.Debug().Msgf("Authenticated user: %d", session.UserId)
	return nil
}

func (auth *authentication) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")
		if err := auth.CheckToken(token); err != nil {
			err = errors.Wrap(err, "error while authenticate")

			auth.log.Error().Stack().Caller().Err(err).Msg("")
			http.Error(rw, err.Error(), http.StatusForbidden)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
