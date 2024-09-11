package api

import (
	"encoding/json"
	"net/http"

	db "www.github.com/TCaprioli/Apptrack-go/db/sqlc"
)

type ApiServer struct {
	addr string
	store *db.Store
}

func NewServer(addr string, store *db.Store) *ApiServer {
	return &ApiServer{addr: addr, store: store}
}

func (s ApiServer) Run() error {
	router := http.NewServeMux()
	// initialize root route
	// TODO: Remove after adding handlers
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		resp:= struct{ Message string `json:"message"`}{Message: "Hello, world"}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w,"failed to encode resp", http.StatusInternalServerError)
		}
	})
	server := http.Server{Addr: s.addr, Handler: router}
	return server.ListenAndServe()
}