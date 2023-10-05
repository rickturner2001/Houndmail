package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"rickturner2001/houndmail/providers"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

type ServerResponse struct {
	Message string `json:"message"`
}

type ServerProviderResponse struct {
	Provider     string `json:"provider"`
	IsRegistered bool   `json:"isRegistered"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewServerProviderResponse(provider string, isRegistered bool) *ServerProviderResponse {
	return &ServerProviderResponse{Provider: provider, IsRegistered: isRegistered}
}

func NewServerResponse(message string) *ServerResponse {
	return &ServerResponse{Message: message}
}

func WithAuth(handlerFunc http.HandlerFunc, s *MySqlStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			log.Println("Token string is empty")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("my-secret-key"), nil
		})
		if err != nil || !token.Valid {
			log.Printf("Token is invalid: %s", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check token expiration
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Token claims are invalid")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		exp, ok := claims["exp"].(float64)
		if !ok {
			log.Println("Token expiration claim is missing or invalid")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		currentTime := time.Now().Unix()
		if currentTime >= int64(exp) {
			log.Println("Token has expired")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}

		fmt.Printf("Authenticated user: %s\n", claims["username"])
		handlerFunc(w, r)
	}
}

func createToken(username string) (string, error) {
	expirationTime := time.Now().Add(time.Hour)

	claims := jwt.MapClaims{
		"username": username,
		"exp":      expirationTime.Unix(),
	}

	// Convert the secret key to a byte slice
	key := []byte("my_secret_key")

	// Create the token with the claims and sign it using the byte slice key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *APIServer) HandlerUserRegistration(w http.ResponseWriter, r *http.Request) error {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Failed to parse JSON data", http.StatusBadRequest)
		return err
	}

	err = s.Store.RegisterUser(&user)
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, user)
}

type TokenResponse struct {
	Token string `json:"token"`
}

func (s *APIServer) HandleUserLogin(w http.ResponseWriter, r *http.Request) error {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return err
	}

	isValid := s.Store.AuthenticateUser(user)

	if !isValid {
		return WriteJson(w, http.StatusNotFound, ServerResponse{Message: "Invalid credentials"})
	}

	token, err := createToken(user.Username)
	if err != nil {
		return err
	}

	r.Header.Set("Authorization", "Bearer "+token)
	return WriteJson(w, http.StatusOK, TokenResponse{Token: token})
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
