package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/vladopadikk/go-chat/internal/config"
)

func Connect(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("faileed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("faileed to ping database: %v", err)
	}

	log.Println("Database connected successfully")
	return db

}
