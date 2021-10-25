package dal

import (
	"context"
	"time"

	"github.com/gofrs/uuid"
)

type Session struct {
	Id          int64
	Token       uuid.UUID
	UserId      int64
	CreatedDate time.Time
	ExpiredDate time.Time
}

type SessionRepository interface {
	CreateSession(int64) (uuid.UUID, error)
	GetSessionByToken(uuid.UUID) (*Session, error)
}

func (db *DB) CreateSession(userId int64) (uuid.UUID, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return uuid.Nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`INSERT INTO sessions(user_id) VALUES ($1) RETURNING token`,
		userId)

	var sessionToken uuid.UUID
	if err = row.Scan(&sessionToken); err != nil {
		return uuid.Nil, err
	}

	return sessionToken, nil
}

func (db *DB) GetSessionByToken(token uuid.UUID) (*Session, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`SELECT id, token, user_id, created_date, expired_date FROM users WHERE token = $1`,
		token)
	var res Session
	if err = row.Scan(&res.Id, &res.Token, &res.UserId, &res.CreatedDate, &res.ExpiredDate); err != nil {
		return nil, err
	}

	return &res, nil
}
