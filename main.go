package main

import (
	"context"
	"log"

	_ "github.com/lib/pq"
	"www.github.com/TCaprioli/Apptrack-go/api"
	db "www.github.com/TCaprioli/Apptrack-go/db/sqlc"
)

var addr = ":8080"

func main() {
	dbCtx := context.Background()
	conn, connErr := db.Connect()
	if connErr != nil {
		log.Fatalf("Failed to connect to the database. %v", connErr)
	}
	log.Print("Connected to the database...")
	store := db.NewStore(conn)
	server := api.NewServer(addr, store, dbCtx)
	log.Printf("Server starting at address %v", addr)
	if serverErr := server.Run(); serverErr != nil {
		log.Fatal(serverErr)
	}

}
