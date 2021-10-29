package dal

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4/stdlib"
	migrate "github.com/rubenv/sql-migrate"
)

func applyMigrations(pool *pgxpool.Pool) (int, error) {
	migrations := &migrate.FileMigrationSource{
		Dir: "./migrations",
	}

	stdDb := stdlib.OpenDB(*pool.Config().ConnConfig)
	defer stdDb.Close()

	return migrate.Exec(stdDb, "postgres", migrations, migrate.Up)
}
