package dal

import (
	"context"

	"github.com/jackc/pgtype"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/phuslu/log"
)

type DB struct {
	pool *pgxpool.Pool
}

func ConnectDB(connectionUri string) (*DB, error) {
	dbconfig, err := pgxpool.ParseConfig(connectionUri)
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

	pool, err := pgxpool.ConnectConfig(context.Background(), dbconfig)
	if err != nil {
		return nil, err
	}

	migrationsNum, err := applyMigrations(pool)
	if err != nil {
		return nil, err
	}
	log.Debug().Msgf("%d migration applied", migrationsNum)

	return &DB{pool}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}
