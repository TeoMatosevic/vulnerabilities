package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Close() error

	TableName() string

	CreateIfNotExists() error

	Query(query string) (*sql.Rows, error)

	QueryResponable(query string, args ...interface{}) (*sql.Rows, error)
}

type service struct {
	db         *sql.DB
	table_name string
}

var (
	connStr    = os.Getenv("DB_CONN_STR")
	dbInstance *service
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	dbInstance = &service{
		db:         db,
		table_name: "users",
	}

	err = dbInstance.CreateIfNotExists()

	if err != nil {
		log.Fatal(err)
	}

	return dbInstance
}

func (s *service) Close() error {
	return s.db.Close()
}

func (s *service) CreateIfNotExists() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            id TEXT PRIMARY KEY,
            first_name TEXT NOT NULL,
            last_name TEXT NOT NULL,
            email TEXT NOT NULL,
            salt TEXT NOT NULL,
            pwd_hash TEXT NOT NULL
        )`, s.table_name))

	if err != nil {
		return err
	}

	return nil
}

func (s *service) TableName() string {
	return s.table_name
}

func (s *service) Query(query string) (*sql.Rows, error) {
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (s *service) QueryResponable(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
