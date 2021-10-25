package dal

import (
	"context"
	"go-react-app/server/config"
	"go-react-app/server/util/logger"

	"github.com/jackc/pgtype"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func ConnectDB(log *logger.Logger, conf *config.DatabaseConfig) (*DB, error) {
	dbconfig, err := pgxpool.ParseConfig(conf.Url)
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
