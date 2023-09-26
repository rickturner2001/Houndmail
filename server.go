package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
}

func NewServer() *mux.Router {
	return mux.NewRouter()
}

func (Server) init(s *Server) {
	s.router.HandleFunc("/", RootRoute)
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}
