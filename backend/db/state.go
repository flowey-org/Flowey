package db

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"maps"
)

type State struct {
	Version int `json:"version"`
}

func GetState(userID UserID) (stateString string, stateVersion int, err error) {
	query := `SELECT state FROM states WHERE user_id = ?`
	err = db.QueryRow(query, userID).Scan(&stateString)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", 0, nil
		}
		log.Println(err)
		return "", 0, InternalServerError
	}

	var state State
	err = json.Unmarshal([]byte(stateString), &state)
	if err != nil {
		log.Println(err)
		return "", 0, InternalServerError
	}

	stateVersion = state.Version
	return stateString, stateVersion, nil
}

func incrementStateVersion(stateString string) (newStateString string, err error) {
	decoder := json.NewDecoder(bytes.NewReader([]byte(stateString)))
	decoder.UseNumber()

	state := make(map[string]interface{})
	if err := decoder.Decode(&state); err != nil {
		log.Printf("failed to parse old state: %v", err)
		return "", InternalServerError
	}

	version, err := state["version"].(json.Number).Int64()
	if err != nil {
		return "", err
	}

	state["version"] = (version + 1) % 1000000

	newStateBytes, err := json.Marshal(state)
	if err != nil {
		log.Printf("failed to marshal new state: %v", err)
		return "", nil
	}

	return string(newStateBytes), nil
}

func equalStates(clientStateString string, serverStateString string) bool {
	clientState := make(map[string]interface{})
	serverState := make(map[string]interface{})

	if err := json.Unmarshal([]byte(clientStateString), &clientState); err != nil {
		return false
	}

	if err := json.Unmarshal([]byte(serverStateString), &serverState); err != nil {
		return false
	}

	return maps.Equal(clientState, serverState)
}

func ChooseState(userID UserID, clientStateString string) (push bool, stateString string, err error) {
	var clientState State
	if err := json.Unmarshal([]byte(clientStateString), &clientState); err != nil {
		return false, "", err
	}
	clientStateVersion := clientState.Version

	serverStateString, serverStateVersion, err := GetState(userID)
	if err != nil {
		return false, "", err
	}

	if equalStates(clientStateString, serverStateString) {
		return false, "", nil
	}

	if clientStateVersion == serverStateVersion {
		newClientStateString, err := incrementStateVersion(clientStateString)
		if err != nil {
			return false, "", err
		}

		SetState(userID, newClientStateString)
		return true, newClientStateString, nil
	} else {
		return true, serverStateString, nil
	}
}

func SetState(userID UserID, stateString string) error {
	query := `INSERT INTO states (user_id, state) VALUES (?, ?) ON CONFLICT (user_id) DO UPDATE SET state = ?`
	_, err := db.Exec(query, userID, stateString, stateString)
	if err != nil {
		log.Println(err)
		return InternalServerError
	}

	return nil
}
