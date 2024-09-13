package api

import (
	"context"
	"log"
	"net/http"
	"strconv"

	db "www.github.com/TCaprioli/Apptrack-go/db/sqlc"
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
