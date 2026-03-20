package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/zyxevls/internal/config"
)

func NewPostgres(cfg *config.Config) *sqlx.DB {
	dsn := cfg.DatabaseURL

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("DB Connection error", err)
	}

	fmt.Println("Successfully connected database✅")
	return db
}
