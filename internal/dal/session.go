package dal

import (
	"context"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

const (
	errSessionInserting = "error while session inserting"
	errSessionReading   = "error while session reading"
)

type Session struct {
	Id          int64     `db:"id"`
	Token       uuid.UUID `db:"token"`
	UserId      int64     `db:"user_id"`
	CreatedDate time.Time `db:"created_date"`
	ExpiredDate time.Time `db:"expired_date"`
}

type SessionRepository interface {
	CreateSession(int64) (uuid.UUID, error)
	GetSessionByToken(uuid.UUID) (*Session, error)
}

func (db *DB) CreateSession(userId int64) (uuid.UUID, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return uuid.Nil, errors.Wrap(err, errFailedAcquireConnection)
	}
	defer conn.Release()

	var sessionToken uuid.UUID
	err = pgxscan.Get(context.Background(), conn, &sessionToken,
		`INSERT INTO sessions(user_id) VALUES ($1) RETURNING token`,
		userId)
	if err != nil {
		return uuid.Nil, errors.Wrap(err, errSessionInserting)
	}

	return sessionToken, nil
}

func (db *DB) GetSessionByToken(token uuid.UUID) (*Session, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var res Session
	err = pgxscan.Get(context.Background(), conn, &res,
		`SELECT id, token, user_id, created_date, expired_date FROM users WHERE token = $1`,
		token)
	if err != nil {
		errors.Wrap(err, errSessionReading)
	}

	return &res, nil
}
