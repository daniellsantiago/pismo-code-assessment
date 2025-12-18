package database

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/nubank/pismo-code-assessment/internal/infrastructure/config"
)

func NewPostgresDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.URL)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
