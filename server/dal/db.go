package dal

import (
	"context"

	"github.com/jackc/pgtype"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func ConnectDB() (*DB, error) {
	databaseUrl := "postgresql://api:123@localhost:9423/TestDB"
	dbconfig, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		return nil, err
	}

	dbconfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		conn.ConnInfo().RegisterDataType(pgtype.DataType{
			Value: &pgtypeuuid.UUID{},
			Name:  "uuid",
			OID:   pgtype.UUIDOID,
		})
		return nil
	}

	db, err := pgxpool.ConnectConfig(context.Background(), dbconfig)
	if err != nil {
		return nil, err
	}

	return &DB{
		pool: db,
	}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) CreateUser(user *User) (int64, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return -1, err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`INSERT INTO users(username, passwordHash) VALUES ($1, $2) RETURNING id`,
		user.Username, user.PasswordHash)

	var userId int64
	if err = row.Scan(&userId); err != nil {
		return -1, err
	}

	return userId, nil
}

func (db *DB) GetUserById(id int64) (*User, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`SELECT id, username, password_hash FROM users WHERE id = $1`,
		id)
	var res User
	if err = row.Scan(&res.Id, &res.Username, &res.PasswordHash); err != nil {
		return nil, err
	}

	return &res, nil
}

func (db *DB) GetUserByUsername(username string) (*User, error) {
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row := conn.QueryRow(context.Background(),
		`SELECT id, username, password_hash FROM users WHERE username = $1`,
		username)
	var res User
	if err = row.Scan(&res.Id, &res.Username, &res.PasswordHash); err != nil {
		return nil, err
	}

	return &res, nil
}
