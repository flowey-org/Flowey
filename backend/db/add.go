package db

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"flag"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func confirmed() bool {
	var input string
	fmt.Println("Are you sure? (y/n)")
	fmt.Scanln(&input)

	if input == "y" || input == "Y" {
		fmt.Print("\033[F\033[F\033[2K\r")
		return true
	}
	return false
}

func generatePassword(length int) (string, error) {
	bytePassword := make([]byte, length)
	_, err := rand.Read(bytePassword)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytePassword), nil
}

func hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func add(path string, username string, passwordLength int) error {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	fmt.Printf("username: %s\n", username)

	password, err := generatePassword(passwordLength)
	if err != nil {
		return err
	}

	fmt.Printf("password: %s\n", password)

	hashedPassword, err := hash(password)
	if err != nil {
		return err
	}

	fmt.Printf("    hash: %s\n", hashedPassword)

	query := `INSERT INTO users (username, password) VALUES (?, ?)`
	_, err = db.Exec(query, username, hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func Add(args []string, path string) error {
	flagSet := flag.NewFlagSet("flowey db add", flag.ExitOnError)
	passwordLength := flagSet.Int("l", 40, "password length")
	skipConfirmation := flagSet.Bool("y", false, "skip confirmation")

	flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: flowey db add [OPTIONS] USERNAME")
		flagSet.PrintDefaults()
	}

	flagSet.Parse(args)

	if flagSet.NArg() != 1 {
		flagSet.Usage()
		return nil
	}

	if !*skipConfirmation && !confirmed() {
		return nil
	}

	username := flagSet.Arg(0)
	return add(path, username, *passwordLength)
}
