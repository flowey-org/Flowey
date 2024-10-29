package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

var (
	Unathorized         = errors.New("unauthorized")
	InternalServerError = errors.New("internal server error")
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AuthenticateByCredentials(credentials Credentials) (int, error) {
	var userID UserID
	var hashedPassword string

	query := `SELECT id, password FROM users WHERE username = ?`
	err := db.QueryRow(query, credentials.Username).Scan(&userID, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, Unathorized
		}
		log.Println(err)
		return -1, InternalServerError
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password))
	if err != nil {
		return -1, Unathorized
	}

	return userID, nil
}

func AuthenticateBySessionKey(sessionKey string) (int, error) {
	var userID UserID

	query := `SELECT user_id FROM sessions WHERE session_key = ?`
	err := db.QueryRow(query, sessionKey).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, Unathorized
		}
		log.Println(err)
		return -1, InternalServerError
	}

	return userID, nil
}

func CreateSessionKey(userID UserID) (string, error) {
	byteSessionKey := make([]byte, 40)
	if _, err := rand.Read(byteSessionKey); err != nil {
		log.Println(err)
		return "", fmt.Errorf("failed to create a session key")
	}
	sessionKey := base64.URLEncoding.EncodeToString(byteSessionKey)

	query := `INSERT INTO sessions (session_key, user_id) VALUES (?, ?)`
	_, err := db.Exec(query, sessionKey, userID)
	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("failed to store a session key")
	}

	return sessionKey, nil
}

func DeleteSessionKey(sessionKey string) error {
	query := `DELETE FROM sessions WHERE session_key = ?`
	_, err := db.Exec(query, sessionKey)
	if err != nil {
		log.Println(err)
		return nil
	}

	return nil
}
