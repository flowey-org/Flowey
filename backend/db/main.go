package db

import (
	"database/sql"
	"flag"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func Main() error {
	flagSet := flag.NewFlagSet("flowey db", flag.ExitOnError)
	path := flagSet.String("db", "flowey.db", "path to the database file")
	flagSet.Parse(os.Args[2:])

	if flagSet.NArg() > 0 {
		flagSet.Usage()
		return nil
	}

	db, err := sql.Open("sqlite3", *path)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return err
	}

	query := `
CREATE TABLE users(
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
