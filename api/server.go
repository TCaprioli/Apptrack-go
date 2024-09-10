package api

import (
	"encoding/json"
	"net/http"
)

type ApiServer struct {
	addr string
}

func NewServer(addr string) *ApiServer {
	return &ApiServer{addr: addr}
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