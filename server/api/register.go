package api

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Stalis/go-react-app/dal"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Success bool `json:"success"`
}

func hashPassword(input string) string {
	data := []byte(input)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func RegisterHandler(rw http.ResponseWriter, r *http.Request) {
	var request RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	entity := dal.User{
		Username:     request.Username,
		PasswordHash: hashPassword(request.Password),
	}

	db, err := dal.ConnectDB()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.CreateUser(&entity)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	payload, _ := json.Marshal(RegisterResponse{Success: true})

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(payload)
}
