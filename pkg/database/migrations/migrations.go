package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

//go:embed sql/*.sql
var fs embed.FS

var (
	CantDoMigrationError = errors.New("can't do migration")
)

func MustMigrate(db *sql.DB) error {
	d, err := iofs.New(fs, "sql")
	if err != nil {
		log.Fatal().Err(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		log.Fatal().Err(err)
	}
	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%w, %s", CantDoMigrationError, err)
	}
	return nil
}
