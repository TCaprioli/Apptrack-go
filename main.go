package main

import (
	"log"

	"www.github.com/TCaprioli/Apptrack-go/api"
)
var addr = ":8080"
func main() {
	server := api.NewServer(addr)
	log.Printf("Server starting at address %v", addr)
	if serverErr := server.Run();serverErr != nil {
		log.Fatal(serverErr)
	} 
	
}