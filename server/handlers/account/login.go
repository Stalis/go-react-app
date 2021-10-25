package account

import (
	"encoding/json"
	"net/http"

	"go-react-app/server/dal"

	"github.com/gofrs/uuid"
	"github.com/phuslu/log"
)

type LoginRequest struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

type LoginResponse struct {
	Success      bool      `json:"success"`
	SessionToken uuid.UUID `json:"sessionToken"`
}

type login struct {
	db *dal.DB
}

func NewLogin(db *dal.DB) http.Handler {
	return &login{db}
}

func (h *login) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Error().Err(err).Msg("Bad request")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.db.GetUserByUsername(request.Username)
	if err != nil {
		log.Error().Err(err).Interface("request", request).Msg("Can't find user")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	response := LoginResponse{Success: false, SessionToken: uuid.Nil}

	hash := request.PasswordHash
	if user.PasswordHash == hash {
		token, err := h.db.CreateSession(user.Id)
		if err != nil {
			log.Error().Err(err).Interface("request", request).Msg("Incorrect password")
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		response.Success = true
		response.SessionToken = token
	}

	payload, _ := json.Marshal(response)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(payload)
}
