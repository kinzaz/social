package main

import (
	"log"

	"github.com/kinzaz/social/internal/db"
	"github.com/kinzaz/social/internal/env"
	"github.com/kinzaz/social/internal/store"
)

func main() {
	addr := env.GetString("DSN", "postgres://admin:pass@localhost:5432/social?sslmode=disable")

	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(store, conn)
}
