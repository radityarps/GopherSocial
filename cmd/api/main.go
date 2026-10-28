package main

import (
	"log"

	"github.com/radityarps/gophersocial/internal/db"
	"github.com/radityarps/gophersocial/internal/env"
	"github.com/radityarps/gophersocial/internal/store"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		dbConfig: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://postgres:postgres@localhost:5432/gophersocial?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 10),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.NewDB(cfg.dbConfig.addr, cfg.dbConfig.maxOpenConns, cfg.dbConfig.maxIdleConns, cfg.dbConfig.maxIdleTime)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("database connection pool established")

	store := store.NewPostgresStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	if err := app.run(mux); err != nil {
		log.Fatal(err)
	}
}
