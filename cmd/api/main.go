package main

import (
	"log"

	"github.com/Kaungmyatkyaw2/go-social/internal/db"
	"github.com/Kaungmyatkyaw2/go-social/internal/env"
	"github.com/Kaungmyatkyaw2/go-social/internal/store"
	"github.com/joho/godotenv"
)

const version = "0.0.2"

//	@title			GoSocial API
//	@description	API for GOSocial Project
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/v1

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Ahutorization
func main() {

	godotenv.Load()

	cfg := config{
		addr:   env.GetString("ADDR", ":8080"),
		env:    env.GetString("ENV", "development"),
		apiURL: env.GetString("API_URL", "localhost:8080"),
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
