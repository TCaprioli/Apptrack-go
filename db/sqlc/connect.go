package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)
func Connect()(*sql.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load environment")
	}
	user := os.Getenv("POSTGRES_USER")
    host := os.Getenv("POSTGRES_HOST")
    password := os.Getenv("POSTGRES_PASSWORD")
    name := os.Getenv("POSTGRES_NAME")
    portStr := os.Getenv("POSTGRES_PORT")
    port, err := strconv.Atoi(portStr)
    if err != nil {
        return nil, fmt.Errorf("invalid port number: %v", err)
    }
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
        "password=%s dbname=%s sslmode=disable",
        host, port, user, password, name)


    return sql.Open("postgres", psqlInfo)

}