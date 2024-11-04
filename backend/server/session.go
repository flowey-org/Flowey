package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"flowey/db"
)

type sessionHandler struct{}

func (handler *sessionHandler) handlePost(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "failed to read the request body", http.StatusBadRequest)
		return
	}

	var credentials db.Credentials
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		http.Error(writer, "couldn't parse the body as a JSON object", http.StatusBadRequest)
		return
	}

	userID, err := db.AuthenticateByCredentials(credentials)
	if err == db.Unathorized {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return
	} else if err == db.InternalServerError {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionToken, err := db.CreateSessionToken(userID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	type loginResponse struct {
		SessionToken string `json:"sessionToken"`
	}
	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(loginResponse{sessionToken})
}

func (handler *sessionHandler) handleDelete(writer http.ResponseWriter, request *http.Request) {
	authHeader := request.Header.Get("Authorization")
	if authHeader == "" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionToken := parts[1]
	db.DeleteSessionToken(sessionToken)

	writer.WriteHeader(http.StatusOK)
}

func (handler *sessionHandler) handleOptions(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	writer.Header().Set("Access-Control-Allow-Methods", "POST, DELETE, OPTIONS")
	writer.WriteHeader(http.StatusOK)
}

func (handler *sessionHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if origin := request.Header.Get("Origin"); origin != "" {
		writer.Header().Set("Access-Control-Allow-Origin", origin)
	}
	writer.Header().Set("Vary", "Origin")

	switch request.Method {
	case http.MethodPost:
		handler.handlePost(writer, request)
	case http.MethodDelete:
		handler.handleDelete(writer, request)
	case http.MethodOptions:
		handler.handleOptions(writer, request)
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}
