package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kzkei/natureAdvice/api/config"

	_ "github.com/lib/pq"
)

func Connect(cfg *config.Config) (*sql.DB, error) {

	log.Printf("connecting to db at %s:%s/%s...", cfg.DBHost, cfg.DBPort, cfg.DBName)
	// use getDSN for connection string
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db at connect: %w", err)
	}

	return db, nil
}
