package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"www.github.com/TCaprioli/Apptrack-go/utils"
)
func Connect()(*sql.DB, error) {
	utils.LoadEnv()

	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		return sql.Open("postgres", databaseURL)
	} else {
		user := os.Getenv("POSTGRES_USER")
		password := os.Getenv("POSTGRES_PASSWORD")
		name := os.Getenv("POSTGRES_NAME")
		portStr := os.Getenv("POSTGRES_PORT")
		port, err := strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid port number: %v", err)
		}
		psqlInfo := fmt.Sprintf("port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			port, user, password, name)
		return sql.Open("postgres", psqlInfo)
	}

}