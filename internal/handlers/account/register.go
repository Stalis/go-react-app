package account

import (
	"encoding/json"
	"net/http"

	"go-react-app/internal/dal"
	"go-react-app/internal/util/logger"
)

type RegisterRequest struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

type UserCreator interface {
	CreateUser(*dal.User) (int64, error)
}

type register struct {
	log   *logger.Logger
	users UserCreator
}

func NewRegister(log *logger.Logger, users UserCreator) http.Handler {
	return &register{log, users}
}

func (h *register) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var request RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.log.Error().Stack().Caller().Err(err).Msg("")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	entity := dal.User{
		Username:     request.Username,
		PasswordHash: request.PasswordHash,
	}

	_, err := h.users.CreateUser(&entity)
	if err != nil {
		h.log.Error().Stack().Caller().Err(err).Msg("")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
