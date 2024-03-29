package pkg

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

type PostgresStorage struct {
	db            *sql.DB
	migrationPath string
	sourceUrl     string
}

func NewPostgresStorage(connection string, migrationPath string) *PostgresStorage {
	db, err := sql.Open("postgres", connection)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &PostgresStorage{db: db, migrationPath: migrationPath, sourceUrl: connection}
}

func (s *PostgresStorage) GetDb() *sql.DB {
	return s.db
}

func (s *PostgresStorage) Migrate() error {
	m, err := migrate.New(
		s.migrationPath,
		s.sourceUrl)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		return err
	}
	return nil
}
