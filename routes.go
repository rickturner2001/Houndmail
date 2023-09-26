package main

import "net/http"

type ServerResponse struct {
	message string
}

func NewServerResponse(message string) *ServerResponse {
	return &ServerResponse{message: message}
}

func RootRoute(w http.ResponseWriter, r *http.Request) {
	WriteJson(w, http.StatusOK, NewServerResponse("Good"))
}
