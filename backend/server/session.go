package server

import (
	"encoding/json"
	"io"
	"net/http"

	"flowey/db"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type sessionHandler struct{}

func (handler *sessionHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "failed to read the request body", http.StatusBadRequest)
		return
	}

	var credentials Credentials
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		http.Error(writer, "couldn't parse the body as a JSON object", http.StatusBadRequest)
		return
	}

	userID, err := db.Authenticate(credentials.Username, credentials.Password)
	if err == db.AuthenticateErrCredentials {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return
	} else if err == db.AuthenticateErrQuery {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionKey, err := db.CreateSessionKey(userID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(writer, &http.Cookie{
		Name:     "session_key",
		Value:    sessionKey,
		HttpOnly: true,
		Secure:   true,
	})
	writer.WriteHeader(http.StatusOK)
}
