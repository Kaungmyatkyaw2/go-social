package main

import (
	"log"

	"github.com/Kaungmyatkyaw2/go-social/internal/db"
	"github.com/Kaungmyatkyaw2/go-social/internal/env"
	"github.com/Kaungmyatkyaw2/go-social/internal/store"
	"github.com/joho/godotenv"
)

const version = "0.0.1"

func main() {

	godotenv.Load()

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		env:  env.GetString("ENV", "development"),
		db: dbConfig{
			dsn:          env.GetString("DB_DSN", "postgres://admin:adminpassword@localhost:5435/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(cfg.db.dsn, cfg.db.maxOpenConns, cfg.db.maxOpenConns, cfg.db.maxIdleTime)

	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Println("database connection pool established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
