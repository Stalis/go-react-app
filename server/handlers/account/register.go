package account

import (
	"encoding/json"
	"net/http"

	"go-react-app/server/dal"
	"go-react-app/server/util/logger"
)

type RegisterRequest struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

type register struct {
	log   *logger.Logger
	users dal.UserRepository
}

func NewRegister(log *logger.Logger, users dal.UserRepository) http.Handler {
	return &register{log, users}
}

func (h *register) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var request RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.log.Error().Stack().Err(err).Msg("")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	entity := dal.User{
		Username:     request.Username,
		PasswordHash: request.PasswordHash,
	}

	_, err := h.users.CreateUser(&entity)
	if err != nil {
		h.log.Error().Stack().Err(err).Msg("")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
