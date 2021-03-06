package account

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"go-react-app/internal/dal"
	"go-react-app/internal/util/logger"

	"github.com/gofrs/uuid"
)

type LoginRequest struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

func (l *LoginRequest) FromJSON(r io.Reader) error {
	return json.NewDecoder(r).Decode(l)
}

type LoginResponse struct {
	SessionToken uuid.UUID `json:"sessionToken"`
}

func (l *LoginResponse) ToJSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(l)
}

type UserByUsernameGetter interface {
	GetUserByUsername(string) (*dal.User, error)
}

type SessionCreator interface {
	CreateSession(int64) (uuid.UUID, error)
}

type login struct {
	log      *logger.Logger
	users    UserByUsernameGetter
	sessions SessionCreator
}

func NewLogin(log *logger.Logger, users UserByUsernameGetter, sessions SessionCreator) http.Handler {
	return &login{log, users, sessions}
}

func (h *login) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.log.Debug().Msg("lol")
	var request LoginRequest
	if err := request.FromJSON(r.Body); err != nil {
		h.log.Error().Stack().Caller().Err(err).Msg("Bad request")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.users.GetUserByUsername(request.Username)
	if err != nil {
		h.log.Error().Stack().Caller().Err(err).Interface("request", request).Msg("Can't find user")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.PasswordHash != request.PasswordHash {
		err = errors.New("incorrect password")
		h.log.Error().Stack().Caller().Err(err).Interface("request", request).Msg("")
		http.Error(rw, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := h.sessions.CreateSession(user.Id)
	if err != nil {
		h.log.Error().Stack().Caller().Err(err).Interface("request", request).Msg("")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	response := LoginResponse{token}

	rw.Header().Set("Content-Type", "application/json")
	response.ToJSON(rw)
}
