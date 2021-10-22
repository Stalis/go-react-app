package account

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
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.db.GetUserByUsername(request.Username)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	response := LoginResponse{Success: false, SessionToken: uuid.Nil}

	hash := hashPassword(request.Password)
	if user.PasswordHash == hash {
		token, err := h.db.CreateSession(user.Id)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		response.Success = true
		response.SessionToken = token
	}

	payload, _ := json.Marshal(response)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(payload)
}
