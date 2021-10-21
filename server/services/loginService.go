package services

import "github.com/gofrs/uuid"

type LoginService interface {
	Register(username string, password string) error
	Login(username string, password string) (uuid.UUID, error)
}
