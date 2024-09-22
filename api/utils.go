package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	db "www.github.com/TCaprioli/Apptrack-go/db/sqlc"
	"www.github.com/TCaprioli/Apptrack-go/utils"
)

func parseId(idStr string, w http.ResponseWriter) int32 {
	var id int32
	if idStr == "" {
		http.Error(w, "Missing Id", http.StatusBadRequest)
		log.Panicf("Missing Id")
	} else {
		parsedId, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			log.Println("Invalid id")
		}
		id = int32(parsedId)
	}
	return id
}

func handleFuncWithCtx(cb func(w http.ResponseWriter, r *http.Request, store *db.Store, ctx context.Context), store *db.Store, ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { cb(w, r, store, ctx) })
}

func verifyToken(tokenString string) (*jwt.Token, error) {
	utils.LoadEnv()
	symmetricKey := os.Getenv("SECRET_KEY")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(symmetricKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid Token")
	}

	return token, nil
}

func createToken(id int32, email string) (string, error) {
	utils.LoadEnv()
	symmetricKey := os.Getenv("SECRET_KEY")

	if symmetricKey == "" {
		return "", fmt.Errorf("key is unset")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    id,
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(symmetricKey))
}
