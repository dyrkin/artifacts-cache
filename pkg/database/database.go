package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"gitlab-cache/pkg/database/migrations"
)

var (
	CantConnectToDatabaseError = errors.New("can't connect to database")
)

type Database interface {
	Connect() error
	Select(statement *sql.Stmt, args ...any) (*sql.Rows, error)
	Update(statement *sql.Stmt, args ...any) (sql.Result, error)
	Statement(query string) (*sql.Stmt, error)
	Migrate() error
}

type database struct {
	connectionString string
	db               *sql.DB
}

func NewDatabase(connectionString string) *database {
	return &database{connectionString: connectionString}
}

func (d *database) Connect() error {
	db, err := sql.Open("postgres", d.connectionString)
	if err != nil {
		return fmt.Errorf("%w. %s", CantConnectToDatabaseError, err)
	}
	d.db = db
	return nil
}

func (d *database) Select(statement *sql.Stmt, args ...any) (*sql.Rows, error) {
	return statement.Query(args...)
}

func (d *database) Update(statement *sql.Stmt, args ...any) (sql.Result, error) {
	return statement.Exec(args...)
}

func (d *database) Statement(query string) (*sql.Stmt, error) {
	return d.db.Prepare(query)
}

func (d *database) Migrate() error {
	return migrations.MustMigrate(d.db)
}
