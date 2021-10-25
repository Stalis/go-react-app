package dal

import (
	"context"
	"go-react-app/server/config"
	"go-react-app/server/util/logger"

	"github.com/jackc/pgtype"
	pgtypeuuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

const (
	errParsingConnectionUrl = "error while parsing connection URL"
	errDbConnection         = "error while connection to database"
	errApplyMigrations      = "error while applying migrations"
)

type DB struct {
	pool *pgxpool.Pool
}

func ConnectDB(log *logger.Logger, conf *config.DatabaseConfig) (*DB, error) {
	dbconfig, err := pgxpool.ParseConfig(conf.Url)
	if err != nil {
		return nil, errors.Wrap(err, errParsingConnectionUrl)
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
		return nil, errors.Wrap(err, errDbConnection)
	}

	migrationsNum, err := applyMigrations(pool)
	if err != nil {
		return nil, errors.Wrap(err, errApplyMigrations)
	}
	log.Debug().Msgf("%d migration applied", migrationsNum)

	return &DB{pool}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) Ping(ctx context.Context) error {
	return db.pool.Ping(ctx)
}
