package dal

import "github.com/google/uuid"

type User struct {
	Id           int64
	Username     string
	PasswordHash string
}

type Session struct {
	Id     int64
	Token  uuid.UUID
	UserId int64
}
