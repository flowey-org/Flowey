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

func (handler *sessionHandler) handlePost(writer http.ResponseWriter, request *http.Request) {
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
	if err == db.Unathorized {
		http.Error(writer, err.Error(), http.StatusUnauthorized)
		return
	} else if err == db.InternalServerError {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	sessionKey, err := db.CreateSessionKey(userID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(writer, &http.Cookie{
		Name:     "flowey_session_key",
		Value:    sessionKey,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   34560000,
	})
	http.SetCookie(writer, &http.Cookie{
		Name:     "flowey_session_key_present",
		Value:    "true",
		Path:     "/",
		HttpOnly: false,
		Secure:   true,
		MaxAge:   34560000,
	})
	writer.WriteHeader(http.StatusOK)
}

func (handler *sessionHandler) handleDelete(writer http.ResponseWriter, request *http.Request) {
	cookie, err := request.Cookie("flowey_session_key")

	if err == nil {
		db.DeleteSessionKey(cookie.Value)
	}

	http.SetCookie(writer, &http.Cookie{
		Name:     "flowey_session_key",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})
	http.SetCookie(writer, &http.Cookie{
		Name:     "flowey_session_key_present",
		Value:    "",
		Path:     "/",
		HttpOnly: false,
		Secure:   true,
		MaxAge:   -1,
	})
	writer.WriteHeader(http.StatusOK)
}

func (handler *sessionHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		handler.handlePost(writer, request)
	case http.MethodDelete:
		handler.handleDelete(writer, request)
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}
