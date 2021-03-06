package session

import (
	"encoding/json"
	"net/http"
	"time"

	"go-react-app/internal/dal"
	"go-react-app/internal/util/logger"

	"github.com/gofrs/uuid"
)

type CheckSessionRequest struct {
	Token uuid.UUID
}

type CheckSessionResponse struct {
	IsExpired   bool
	ExpiredTime time.Time
}

type SessionGetter interface {
	GetSessionByToken(uuid.UUID) (*dal.Session, error)
}

type checkSession struct {
	log      *logger.Logger
	sessions SessionGetter
}

func NewCheck(log *logger.Logger, sessions SessionGetter) http.Handler {
	return &checkSession{log, sessions}
}

func (h *checkSession) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var request CheckSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.log.Error().Stack().Err(err).Msg("")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	session, err := h.sessions.GetSessionByToken(request.Token)
	if err != nil {
		h.log.Error().Stack().Err(err).Msg("")
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	response := CheckSessionResponse{
		ExpiredTime: session.ExpiredDate,
		IsExpired:   session.ExpiredDate.After(time.Now()),
	}

	payload, _ := json.Marshal(response)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(payload)
}
