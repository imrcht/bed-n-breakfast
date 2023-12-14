package driver

import (
	"database/sql"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB: is a struct that holds the database connection
type DB struct {
	SQL *sql.DB
}

// * One of the popular ORM for Go is upperDb/upper.io

func ConnectSql(dsn string) (*DB, error) {
	d, err := NewDbConnection(dsn)

	if err != nil {
		panic(err)
	}

	err = testDBPing(d)
	if err != nil {
		panic(err)
	}

	db := &DB{SQL: d}

	return db, nil
}

func testDBPing(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func NewDbConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
