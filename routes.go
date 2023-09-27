package main

import (
	"log"
	"net/http"
	"rickturner2001/houndmail/providers"

	"github.com/gorilla/mux"
)

type ServerResponse struct {
	Message string `json:"message"`
}

type ServerProviderResponse struct {
	Provider     string `json:"provider"`
	IsRegistered bool   `json:"isRegistered"`
}

func NewServerProviderResponse(provider string, isRegistered bool) *ServerProviderResponse {
	return &ServerProviderResponse{Provider: provider, IsRegistered: isRegistered}
}

func NewServerResponse(message string) *ServerResponse {
	return &ServerResponse{Message: message}
}

func (s *APIServer) HandleRootRoute(w http.ResponseWriter, r *http.Request) error {
	log.Println("Root route was called")
	return WriteJson(w, http.StatusOK, NewServerResponse("Hey"))
}

func (s *APIServer) HandleProviderRoute(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	providerName, err := ExtractRequestVarOrInvalid("provider", vars)
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}

	provider, err := providers.GetProviderByName(providerName)
	provider.HandleResponse(s.Client)
	if err != nil {
		return WriteJson(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
	}
	return WriteJson(w, http.StatusOK, NewServerProviderResponse(provider.Name, true))
}
