package api

import (
	"context"
	"net/http"

	db "www.github.com/TCaprioli/Apptrack-go/db/sqlc"
)

type ApiServer struct {
	addr  string
	store *db.Store
	ctx   context.Context
}

func NewServer(addr string, store *db.Store, ctx context.Context) *ApiServer {
	return &ApiServer{addr: addr, store: store, ctx: ctx}
}

func (s ApiServer) Run() error {
	router := http.NewServeMux()
	// Handles login and register
	router.Handle("/users/", userHandler(s.store, s.ctx))
	// Handles create and get one
	router.Handle("/applications", authMiddleware(applicationHandler(s.store, s.ctx)))
	// Handles get one, update, and delete
	router.Handle("/applications/", authMiddleware(applicationIdHandler(s.store, s.ctx)))

	server := http.Server{Addr: s.addr, Handler: router}
	return server.ListenAndServe()
}
