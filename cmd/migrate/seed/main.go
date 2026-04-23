package main

import (
	"log"

	"github.com/Kaungmyatkyaw2/go-social/internal/db"
	"github.com/Kaungmyatkyaw2/go-social/internal/env"
	"github.com/Kaungmyatkyaw2/go-social/internal/store"
)

func main() {

	addr := env.GetString("DB_DSN","postgres://admin:adminpassword@localhost:5435/social?sslmode=disable")
	
	conn, err := db.New(addr,3,3,"15m")


	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()


	store := store.NewStorage(conn)

	db.Seed(store)
}