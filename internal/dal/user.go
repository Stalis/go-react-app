package dal

import (
	"context"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
)

const (
	errFailedAcquireConnection = "failed acquire pgxpool connection"
	errUserInserting           = "error while inserting user"
	errUserReading             = "error while reading user"
)

type User struct {
	Id           int64     `db:"id"`
	Username     string    `db:"username"`
	PasswordHash string    `db:"password_hash"`
	CreatedDate  time.Time `db:"created_date"`
}

func (db *DB) CreateUser(user *User) (int64, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return -1, errors.Wrap(err, errFailedAcquireConnection)
	}
	defer conn.Release()

	var userId int64
	err = pgxscan.Get(context.Background(), conn, &userId,
		`INSERT INTO users(username, password_hash) VALUES ($1, $2) RETURNING id`,
		user.Username, user.PasswordHash)
	if err != nil {
		return -1, errors.Wrap(err, errUserInserting)
	}

	return userId, nil
}

func (db *DB) GetUserById(id int64) (*User, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, errFailedAcquireConnection)
	}
	defer conn.Release()

	var res User
	err = pgxscan.Get(context.Background(), conn, &res,
		`SELECT id, username, password_hash, created_date FROM users WHERE id = $1`,
		id)
	if err != nil {
		return nil, errors.Wrap(err, errUserReading)
	}

	return &res, nil
}

func (db *DB) GetUserByUsername(username string) (*User, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, errFailedAcquireConnection)
	}
	defer conn.Release()

	var res User
	err = pgxscan.Get(context.Background(), conn, &res,
		`SELECT id, username, password_hash, created_date FROM users WHERE username = $1`,
		username)
	if err != nil {
		return nil, errors.Wrap(err, errUserReading)
	}

	return &res, nil
}
