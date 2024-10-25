package server

import (
	"encoding/json"
	"io"
	"net/http"

	"flowey/db"
)

const (
	sessionKeyCookieName        = "flowey_session_key"
	sessionKeyPresentCookieName = "flowey_session_key_present"
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

	sessionKey, err := db.CreateSessionKey(userID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(writer, &http.Cookie{
		Name:     sessionKeyCookieName,
		Value:    sessionKey,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   34560000,
	})
	http.SetCookie(writer, &http.Cookie{
		Name:     sessionKeyPresentCookieName,
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
		Name:     sessionKeyCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})
	http.SetCookie(writer, &http.Cookie{
		Name:     sessionKeyPresentCookieName,
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
