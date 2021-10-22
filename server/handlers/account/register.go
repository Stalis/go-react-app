package account

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Stalis/go-react-app/server/dal"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Success bool `json:"success"`
}

type register struct {
	db *dal.DB
}

func NewRegister(db *dal.DB) http.Handler {
	return &register{db}
}

func (h *register) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var request RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	entity := dal.User{
		Username:     request.Username,
		PasswordHash: hashPassword(request.Password),
	}

	_, err := h.db.CreateUser(&entity)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	payload, _ := json.Marshal(RegisterResponse{Success: true})

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(payload)
}

func hashPassword(input string) string {
	data := []byte(input)
	return fmt.Sprintf("%x", md5.Sum(data))
}
