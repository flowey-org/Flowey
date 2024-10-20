package db

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
)

func Occupied(path string) (bool, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			return true, nil
		} else {
			return false, err
		}
	}
	file.Close()
	return false, nil
}

func Create(path string) error {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	query := `CREATE TABLE users(
  id INTEGER NOT NULL PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL
);`
	if _, err := db.Exec(query); err != nil {
		log.Println(err)
		return nil
	}

	return nil
}

func Init(args []string, path string) error {
	flagSet := flag.NewFlagSet("flowey db init", flag.ExitOnError)

	flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: flowey db init")
	}

	flagSet.Parse(args)

	if flagSet.NArg() > 0 {
		flagSet.Usage()
		return nil
	}

	if occupied, err := Occupied(path); err != nil {
		return err
	} else if occupied {
		fmt.Printf("the file %s already exists, doing nothing\n", path)
		return nil
	}

	return Create(path)
}
