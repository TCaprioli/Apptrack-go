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
	router.Handle("/users/", userHandler(s.store, s.ctx))
	router.Handle("/applications", authMiddleware(applicationHandler(s.store, s.ctx)))

	server := http.Server{Addr: s.addr, Handler: router}
	return server.ListenAndServe()
}
