package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type APIServer struct {
	ListenAddr string
	Client     *http.Client
	Store      *MySqlStore
}

func NewAPIServer(store *MySqlStore, listenAddr string) *APIServer {
	return &APIServer{
		ListenAddr: listenAddr,
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		Store: store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/", MakeHttpFunc(s.HandleRootRoute)).Methods(http.MethodGet)
	router.HandleFunc("/protected", WithAuth(MakeHttpFunc(s.HandleRootRoute), s.Store))
	router.HandleFunc("/provider/{provider}", MakeHttpFunc(s.HandleProviderRoute)).Methods(http.MethodGet)
	router.HandleFunc("/user/register", MakeHttpFunc(s.HandlerUserRegistration)).Methods(http.MethodPost)
	router.HandleFunc("/user/validate", MakeHttpFunc(s.HandleUserLogin)).Methods(http.MethodPost)
	log.Printf("API running on address: %s", s.ListenAddr)

	err := http.ListenAndServe(s.ListenAddr, router)
	if err != nil {
		log.Panicf("Could not listen on port %s: %v", s.ListenAddr, err)
	}
}

type (
	apiFunc  func(http.ResponseWriter, *http.Request) error
	ApiError struct {
		Error string `json:"error"`
	}
)

func MakeHttpFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			err := WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
			if err != nil {
				log.Printf("Could not write a response: %v", err)
				return
			}
		}
	}
}

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func ExtractRequestVarOrInvalid(key string, vars map[string]string) (string, error) {
	val := vars[key]

	if val == "" {
		return "", fmt.Errorf("Could not find variable with key %s", key)
	}

	return val, nil
}
