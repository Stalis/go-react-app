package dal

import (
	"context"
	"time"

	"github.com/pkg/errors"
)

const (
	errFailedAcquireConnection = "failed acquire pgxpool connection"
	errUserInserting           = "error while inserting user"
	errUserReading             = "error while reading user"
)

type User struct {
	Id           int64
	Username     string
	PasswordHash string
	CreatedDate  time.Time
}

type UserRepository interface {
	CreateUser(*User) (int64, error)
	GetUserById(int64) (*User, error)
	GetUserByUsername(string) (*User, error)
}

func (db *DB) CreateUser(user *User) (int64, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return -1, errors.Wrap(err, errFailedAcquireConnection)
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`INSERT INTO users(username, password_hash) VALUES ($1, $2) RETURNING id`,
		user.Username, user.PasswordHash)

	var userId int64
	if err = row.Scan(&userId); err != nil {
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

	row := conn.QueryRow(context.Background(),
		`SELECT id, username, password_hash, created_date FROM users WHERE id = $1`,
		id)

	var res User
	if err = row.Scan(&res.Id, &res.Username, &res.PasswordHash, &res.CreatedDate); err != nil {
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

	row := conn.QueryRow(context.Background(),
		`SELECT id, username, password_hash, created_date FROM users WHERE username = $1`,
		username)

	var res User
	if err = row.Scan(&res.Id, &res.Username, &res.PasswordHash, &res.CreatedDate); err != nil {
		return nil, errors.Wrap(err, errUserReading)
	}

	return &res, nil
}
