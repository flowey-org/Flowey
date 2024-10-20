package utils

import (
	"os"
	"path/filepath"
)

func GetDataDir() (string, error) {
	dataHome, present := os.LookupEnv("XDG_DATA_HOME")
	if !present {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		dataHome = filepath.Join(homeDir, ".local", "share")
	}

	dataDir := filepath.Join(dataHome, "flowey")
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return "", err
	}

	return dataDir, nil
}
