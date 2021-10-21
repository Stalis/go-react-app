package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Stalis/go-react-app/dal"
	"github.com/gofrs/uuid"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	IsSuccess    bool      `json:"isSuccess"`
	SessionToken uuid.UUID `json:"sessionToken"`
}

type Login struct {
	db *dal.DB
}

func NewLogin(db *dal.DB) *Login {
	return &Login{db}
}

func (h *Login) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.db.GetUserByUsername(request.Username)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	response := LoginResponse{IsSuccess: false, SessionToken: uuid.Nil}

	hash := hashPassword(request.Password)
	if user.PasswordHash == hash {
		token, err := h.db.CreateSession(user.Id)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		response.IsSuccess = true
		response.SessionToken = token
	}

	payload, _ := json.Marshal(response)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(payload)
}
