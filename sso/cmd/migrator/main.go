package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storageDSN, migrationsPath, migrationsTable string

	flag.StringVar(&storageDSN, "storage-dsn", "", "Postgres DSN (e.g. postgres://user:pass@host:5432/dbname?sslmode=disable)")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations (filesystem) e.g. ./migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()

	if storageDSN == "" {
		log.Fatalf("storage-dsn is required")
	}
	if migrationsPath == "" {
		log.Fatalf("migrations-path is required")
	}

	dbURL, err := url.Parse(storageDSN)
	if err != nil {
		log.Fatalf("invalid storage-dsn: %v", err)
	}

	q := dbURL.Query()
	q.Set("x-migrations-table", migrationsTable)
	dbURL.RawQuery = q.Encode()

	m, err := migrate.New(
		"file://"+migrationsPath,
		dbURL.String(),
	)
	if err != nil {
		log.Fatalf("migrate.New: %v", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil || dbErr != nil {
			if srcErr != nil {
				fmt.Fprintf(os.Stderr, "m.Close source error: %v\n", srcErr)
			}
			if dbErr != nil {
				fmt.Fprintf(os.Stderr, "m.Close db error: %v\n", dbErr)
			}
		}
	}()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		log.Fatalf("m.Up: %v", err)
	}

	fmt.Println("migrations applied")
}

type Log struct {
	verbose bool
}

func (l *Log) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *Log) Verbose() bool {
	return false
}
