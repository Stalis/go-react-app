package account

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Stalis/go-react-app/server/dal"
	"github.com/gofrs/uuid"
)

type CheckSessionRequest struct {
	Token uuid.UUID
}

type CheckSessionResponse struct {
	IsExpired   bool
	ExpiredTime time.Time
}

type checkSession struct {
	db *dal.DB
}

func NewCheckSession(db *dal.DB) http.Handler {
	return &checkSession{db}
}

func (h *checkSession) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var request CheckSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	session, err := h.db.GetSessionByToken(request.Token)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	response := CheckSessionResponse{
		ExpiredTime: session.ExpiredDate,
		IsExpired:   session.ExpiredDate.After(time.Now()),
	}

	payload, _ := json.Marshal(response)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(payload)
}
